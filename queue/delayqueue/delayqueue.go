package delayqueue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type DelayQueue struct {
	name          string
	redisCli      *redis.Client
	cb            func(string) bool
	pendingKey    string
	readyKey      string
	unAckKey      string
	retryKey      string
	retryCountKey string
	garbageKey    string
	ticker        *time.Ticker
	logger        *log.Logger
	close         chan struct{}

	maxConsumerDuration time.Duration
	msgTTL              time.Duration
	defaultRetryCount   uint
	fetchInterval       time.Duration
	fetchLimit          uint

	concurrent uint
}

func NewQueue(name string, cli *redis.Client, callback func(string) bool) *DelayQueue {
	if name == "" {
		panic("name is required")
	}
	if cli == nil {
		panic("cli is required")
	}
	if callback == nil {
		panic("callback is required")
	}

	return &DelayQueue{
		name:                name,
		redisCli:            cli,
		cb:                  callback,
		pendingKey:          "dp:" + name + ":pending",
		readyKey:            "dp:" + name + ":ready",
		unAckKey:            "dp:" + name + ":unack",
		retryKey:            "dp:" + name + ":retry",
		retryCountKey:       "dp:" + name + ":retry:cnt",
		garbageKey:          "dp:" + name + ":garbage",
		logger:              log.Default(),
		close:               make(chan struct{}),
		maxConsumerDuration: 5 * time.Second,
		msgTTL:              1 * time.Hour,
		defaultRetryCount:   3,
		fetchInterval:       1 * time.Second,
		fetchLimit:          0,
		concurrent:          1,
	}
}

func (q *DelayQueue) WithLogger(logger *log.Logger) *DelayQueue {
	q.logger = logger
	return q
}

func (q *DelayQueue) WithFetchInterval(d time.Duration) *DelayQueue {
	q.fetchInterval = d
	return q
}

func (q *DelayQueue) WithMaxConsumeDuration(d time.Duration) *DelayQueue {
	q.maxConsumerDuration = d
	return q
}

func (q *DelayQueue) WithFetchLimit(limit uint) *DelayQueue {
	q.fetchLimit = limit
	return q
}

func (q *DelayQueue) WithConcurrent(c uint) *DelayQueue {
	if c == 0 {
		return q
	}
	q.concurrent = c
	return q
}

func (q *DelayQueue) WithDefaultRetryCount(count uint) *DelayQueue {
	q.defaultRetryCount = count
	return q
}

func (q *DelayQueue) genMsgKey(idStr string) string {
	return "dp:" + q.name + ":msg:" + idStr
}

type retryCountOpt int

func WithRetryCount(count int) interface{} {
	return retryCountOpt(count)
}

func (q *DelayQueue) SendScheduleMsg(payload string, t time.Time, opts ...interface{}) error {
	retryCount := q.defaultRetryCount
	for _, opt := range opts {
		switch o := opt.(type) {
		case retryCountOpt:
			retryCount = uint(o)
		}
	}

	idStr := uuid.Must(uuid.NewRandom()).String()
	ctx := context.Background()
	now := time.Now()
	msgTTL := t.Sub(now) + q.msgTTL
	err := q.redisCli.Set(ctx, q.genMsgKey(idStr), payload, msgTTL).Err()
	if err != nil {
		return fmt.Errorf("store msg failed: %v", err)
	}
	err = q.redisCli.HSet(ctx, q.retryCountKey, idStr, retryCount).Err()
	if err != nil {
		return fmt.Errorf("store retry count failed: %v", err)
	}
	err = q.redisCli.ZAdd(ctx, q.pendingKey, &redis.Z{
		Score:  float64(t.Unix()),
		Member: idStr,
	}).Err()
	if err != nil {
		return fmt.Errorf("push to pending failed: %v", err)
	}
	return nil
}

func (q *DelayQueue) SendDelayMsg(payload string, duration time.Duration, opts ...interface{}) error {
	t := time.Now().Add(duration)
	return q.SendScheduleMsg(payload, t, opts...)
}

const pending2ReadyScript = `
local msgs = redis.call('ZRangeByScore', KEYS[1], '0', ARGV[1])  -- get ready msg
if (#msgs == 0) then return end
local args2 = {'LPush', KEYS[2]} -- push into ready
for _,v in ipairs(msgs) do
	table.insert(args2, v) 
end
redis.call(unpack(args2))
redis.call('ZRemRangeByScore', KEYS[1], '0', ARGV[1])  -- remove msgs from pending
`

func (q *DelayQueue) pending2Ready() error {
	now := time.Now().Unix()
	ctx := context.Background()
	keys := []string{q.pendingKey, q.retryKey}
	err := q.redisCli.Eval(ctx, pending2ReadyScript, keys, now).Err()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("pending2ReadyScript failed: %v", err)
	}
	return nil
}

const ready2UnackScript = `
local msg = redis.call('RPop', KEYS[1])
if (not msg) then return end
redis.call('ZAdd', KEYS[2], ARGV[1], msg)
return msg
`

func (q *DelayQueue) ready2Unack() (string, error) {
	retryTime := time.Now().Add(q.maxConsumerDuration).Unix()
	ctx := context.Background()
	keys := []string{q.readyKey, q.unAckKey}
	ret, err := q.redisCli.Eval(ctx, ready2UnackScript, keys, retryTime).Result()
	if err == redis.Nil {
		return "", err
	}
	if err != nil {
		return "", fmt.Errorf("ready2UnackScript failed: %v", err)
	}
	str, ok := ret.(string)
	if !ok {
		return "", fmt.Errorf("illagal result: %#v", ret)
	}
	return str, nil
}

func (q *DelayQueue) retry2Unack() (string, error) {
	retryTime := time.Now().Add(q.maxConsumerDuration).Unix()
	ctx := context.Background()
	keys := []string{q.retryKey, q.unAckKey}
	ret, err := q.redisCli.Eval(ctx, ready2UnackScript, keys, retryTime).Result()
	if err == redis.Nil {
		return "", redis.Nil
	}
	if err != nil {
		return "", fmt.Errorf("ready2UnackScript failed: %v", err)
	}
	str, ok := ret.(string)
	if !ok {
		return "", fmt.Errorf("illagal result: %#v", ret)
	}
	return str, nil
}

func (q *DelayQueue) callback(idStr string) error {
	ctx := context.Background()
	payload, err := q.redisCli.Get(ctx, q.genMsgKey(idStr)).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return fmt.Errorf("get message payload failed: %v", err)
	}
	ack := q.cb(payload)
	if ack {
		err = q.ack(idStr)
	} else {
		err = q.nack(idStr)
	}
	return err
}

func (q *DelayQueue) batchCallback(ids []string) {
	if len(ids) == 1 || q.concurrent == 1 {
		for _, id := range ids {
			err := q.callback(id)
			if err != nil {
				q.logger.Printf("consume msg %s failed: %v", id, err)
			}
		}
		return
	}

	ch := make(chan string, len(ids))
	for _, id := range ids {
		ch <- id
	}
	close(ch)
	wg := sync.WaitGroup{}
	concurrent := int(q.concurrent)
	if concurrent > len(ids) {
		concurrent = len(ids)
	}
	wg.Add(concurrent)

	for i := 0; i < concurrent; i++ {
		go func() {
			defer wg.Done()
			for id := range ch {
				err := q.callback(id)
				if err != nil {
					q.logger.Printf("consume msg %s failed: %v", id, err)
				}
			}
		}()
	}
	wg.Wait()
}

func (q *DelayQueue) ack(idStr string) error {
	ctx := context.Background()
	err := q.redisCli.ZRem(ctx, q.unAckKey, idStr).Err()
	if err != nil {
		return fmt.Errorf("remove from unack failed: %v", err)
	}

	_ = q.redisCli.Del(ctx, q.genMsgKey(idStr)).Err()
	q.redisCli.HDel(ctx, q.retryCountKey, idStr)
	return nil
}

func (q *DelayQueue) nack(idStr string) error {
	ctx := context.Background()
	err := q.redisCli.ZAdd(ctx, q.unAckKey, &redis.Z{
		Member: idStr,
		Score:  float64(time.Now().Unix()),
	}).Err()
	if err != nil {
		return fmt.Errorf("negative ack failed: %v", err)
	}
	return nil
}

const unack2RetryScript = `
local msgs = redis.call('ZRangeByScore', KEYS[1], '0', ARGV[1])  -- get retry msg
if (#msgs == 0) then return end
local retryCounts = redis.call('HMGet', KEYS[2], unpack(msgs)) -- get retry count
for i,v in ipairs(retryCounts) do
	local k = msgs[i]
	if tonumber(v) > 0 then
		redis.call("HIncrBy", KEYS[2], k, -1) -- reduce retry count
		redis.call("LPush", KEYS[3], k) -- add to retry
	else
		redis.call("HDel", KEYS[2], k) -- del retry count
		redis.call("SAdd", KEYS[4], k) -- add to garbage
	end
end
redis.call('ZRemRangeByScore', KEYS[1], '0', ARGV[1])  -- remove msgs from unack
`

func (q *DelayQueue) unack2Retry() error {
	ctx := context.Background()
	keys := []string{q.unAckKey, q.retryCountKey, q.retryKey, q.garbageKey}
	now := time.Now()
	err := q.redisCli.Eval(ctx, unack2RetryScript, keys, now.Unix()).Err()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("unack to retry script failed: %v", err)
	}
	return nil
}

func (q *DelayQueue) garbageCollect() error {
	ctx := context.Background()
	msgIds, err := q.redisCli.SMembers(ctx, q.garbageKey).Result()
	if err != nil {
		return fmt.Errorf("smembers failed: %v", err)
	}
	if len(msgIds) == 0 {
		return nil
	}
	msgKeys := make([]string, 0, len(msgIds))
	for _, idStr := range msgIds {
		msgKeys = append(msgKeys, q.genMsgKey(idStr))
	}
	err = q.redisCli.Del(ctx, msgKeys...).Err()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("del msgs failed: %v", err)
	}
	err = q.redisCli.SRem(ctx, q.garbageKey, msgIds).Err()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("remove from garbage key failed: %v", err)
	}
	return nil
}

func (q *DelayQueue) consume() error {
	err := q.pending2Ready()
	if err != nil {
		return err
	}
	ids := make([]string, 0, q.fetchLimit)
	for {
		idStr, err := q.ready2Unack()
		if err == redis.Nil {
			break
		}
		if err != nil {
			return err
		}
		ids = append(ids, idStr)
		if q.fetchLimit > 0 && len(ids) >= int(q.fetchLimit) {
			break
		}
	}
	if len(ids) > 0 {
		q.batchCallback(ids)
	}
	err = q.unack2Retry()
	if err != nil {
		return err
	}
	err = q.garbageCollect()
	if err != nil {
		return err
	}
	ids = make([]string, 0, q.fetchLimit)
	for {
		idStr, err := q.retry2Unack()
		if err == redis.Nil {
			break
		}
		if err != nil {
			return err
		}
		ids = append(ids, idStr)
		if q.fetchLimit > 0 && len(ids) >= int(q.fetchLimit) {
			break
		}
	}
	if len(ids) > 0 {
		q.batchCallback(ids)
	}
	return nil
}

func (q *DelayQueue) StartConsume() (done <-chan struct{}) {
	q.ticker = time.NewTicker(q.fetchInterval)
	done0 := make(chan struct{})
	go func() {
	tickerLoop:
		for {
			select {
			case <-q.ticker.C:
				err := q.consume()
				if err != nil {
					log.Printf("consume error: %v", err)
				}
			case <-q.close:
				break tickerLoop
			}
		}
		close(done0)
	}()
	return done0
}

func (q *DelayQueue) StopConsume() {
	close(q.close)
	if q.ticker != nil {
		q.ticker.Stop()
	}
}

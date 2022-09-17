package main

import (
	"strconv"
	"time"

	"github.com/Muskchen/toolkits/queue/delayqueue"
	"github.com/go-redis/redis/v8"
)

func main() {
	redisCli := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	queue := delayqueue.NewQueue("dome", redisCli, func(payload string) bool {
		println(payload)
		return true
	}).WithConcurrent(4)

	for i := 0; i < 10; i++ {
		err := queue.SendDelayMsg(strconv.Itoa(i), time.Second, delayqueue.WithRetryCount(3))
		if err != nil {
			panic(err)
		}
	}

	for i := 0; i < 10; i++ {
		err := queue.SendScheduleMsg(strconv.Itoa(i), time.Now().Add(time.Second))
		if err != nil {
			panic(err)
		}
	}
	done := queue.StartConsume()
	<-done
}

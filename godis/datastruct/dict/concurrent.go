package dict

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
)

type ConcurrentDict struct {
	table []*shard
	count int32
}

type shard struct {
	m     map[string]interface{}
	mutex sync.RWMutex
}

func computeCapacity(param int) (size int) {
	if param <= 16 {
		return 16
	}
	n := param - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n < 0 {
		return math.MaxInt32
	}
	return n + 1
}

func MakeConcurrent(shardCount int) *ConcurrentDict {
	shardCount = computeCapacity(shardCount)
	table := make([]*shard, shardCount)
	for i := 0; i < shardCount; i++ {
		table[i] = &shard{
			m: make(map[string]interface{}),
		}
	}
	return &ConcurrentDict{
		table: table,
		count: 0,
	}
}

const prime32 = uint32(16777619)

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func (dict *ConcurrentDict) spread(hashCode uint32) uint32 {
	if dict == nil {
		panic("dict is nil")
	}
	tableSize := uint32(len(dict.table))
	return (tableSize - 1) & uint32(hashCode)
}

func (dict *ConcurrentDict) getShard(index uint32) *shard {
	if dict == nil {
		panic("dict is nil")
	}
	return dict.table[index]
}

func (dict *ConcurrentDict) makeShard(key string) *shard {
	if dict == nil {
		panic("dict is nil")
	}
	hashCode := fnv32(key)
	index := dict.spread(hashCode)
	shard := dict.getShard(index)
	return shard
}

func (dict *ConcurrentDict) Get(key string) (val interface{}, exists bool) {
	shard := dict.makeShard(key)
	shard.mutex.RLock()
	defer shard.mutex.RUnlock()
	val, exists = shard.m[key]
	return val, exists
}

func (dict *ConcurrentDict) Len() int {
	if dict == nil {
		panic("dict is nil")
	}
	return int(atomic.LoadInt32(&dict.count))
}

func (dict *ConcurrentDict) Put(key string, val interface{}) (result int) {
	shard := dict.makeShard(key)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()
	if _, ok := shard.m[key]; ok {
		shard.m[key] = val
		return 0
	}
	shard.m[key] = val
	dict.addCount()
	return 1
}

func (dict *ConcurrentDict) PutIfAbsent(key string, val interface{}) (result int) {
	shard := dict.makeShard(key)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()

	if _, ok := shard.m[key]; ok {
		return 0
	}
	shard.m[key] = val
	dict.addCount()
	return 1
}

func (dict *ConcurrentDict) PutIfExists(key string, val interface{}) (result int) {
	shard := dict.makeShard(key)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()

	if _, ok := shard.m[key]; ok {
		shard.m[key] = val
		return 1
	}
	return 0
}

func (dict *ConcurrentDict) Remove(key string) (result int) {
	shard := dict.makeShard(key)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()

	if _, ok := shard.m[key]; ok {
		delete(shard.m, key)
		return 1
	}
	return 0
}

func (dict *ConcurrentDict) addCount() int32 {
	return atomic.AddInt32(&dict.count, 1)
}

func (dict *ConcurrentDict) ForEach(consumer Consumer) {
	if dict == nil {
		panic("dict is nil")
	}
	for _, shard := range dict.table {
		shard.mutex.Lock()
		func() {
			defer shard.mutex.Unlock()
			for key, value := range shard.m {
				if !consumer(key, value) {
					return
				}
			}
		}()
	}
}

func (dict *ConcurrentDict) Keys() []string {
	keys := make([]string, dict.Len())
	i := 0
	dict.ForEach(func(key string, val interface{}) bool {
		if i < len(keys) {
			keys[i] = key
			i++
		} else {
			keys = append(keys, key)
		}
		return true
	})
	return keys
}

func (shard *shard) RandomKey() string {
	if shard == nil {
		panic("shard is nil")
	}
	shard.mutex.Lock()
	defer shard.mutex.Unlock()
	for key := range shard.m {
		return key
	}
	return ""
}

func (dict *ConcurrentDict) RandomKeys(limit int) []string {
	size := dict.Len()
	if limit >= size {
		return dict.Keys()
	}
	shardCount := len(dict.table)
	result := make([]string, limit)
	for i := 0; i < limit; {
		shard := dict.getShard(uint32(rand.Intn(shardCount)))
		if shard == nil {
			continue
		}
		key := shard.RandomKey()
		if key != "" {
			result[i] = key
			i++
		}
	}
	return result
}

func (dict *ConcurrentDict) RandomDistinctKeys(limit int) []string {
	size := dict.Len()
	if limit >= size {
		return dict.Keys()
	}
	shardCount := len(dict.table)
	result := make(map[string]bool)
	for len(result) < limit {
		shard := dict.getShard(uint32(rand.Intn(shardCount)))
		if shard == nil {
			continue
		}
		key := shard.RandomKey()
		if key != "" {
			result[key] = true
		}
	}
	arr := make([]string, limit)
	i := 0
	for k := range result {
		arr[i] = k
		i++
	}
	return arr
}

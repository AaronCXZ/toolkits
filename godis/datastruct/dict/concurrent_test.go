package dict

import (
	"errors"
	"strconv"
	"sync"
	"testing"

	"github.com/Muskchen/toolkits/godis/lib/utils"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentDictPut(t *testing.T) {
	d := MakeConcurrent(0)
	count := 100
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			key := "k" + strconv.Itoa(i)
			ret := d.Put(key, i)
			assert.Equal(t, 1, ret)
			val, ok := d.Get(key)
			if ok {
				intVal, _ := val.(int)
				assert.Equal(t, i, intVal)
			} else {
				_, ok := d.Get(key)
				assert.Equal(t, true, ok)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestConcurrentDictPutIfAbsent(t *testing.T) {
	d := MakeConcurrent(0)
	count := 100
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			key := "k" + strconv.Itoa(i)
			ret := d.PutIfAbsent(key, i)
			assert.Equal(t, 1, ret)

			val, ok := d.Get(key)
			if ok {
				intVal, _ := val.(int)
				assert.Equal(t, i, intVal)
			} else {
				_, ok := d.Get(key)
				assert.Equal(t, true, ok)
			}

			ret = d.PutIfAbsent(key, i*10)
			assert.Equal(t, 0, ret)
			val, ok = d.Get(key)
			if ok {
				intVal, _ := val.(int)
				assert.Equal(t, i, intVal)
			} else {
				_, ok := d.Get(key)
				assert.Equal(t, true, ok)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestConcurrentDictRemove(t *testing.T) {
	t.Run("remove head node", func(t *testing.T) {
		d := MakeConcurrent(0)
		for i := 0; i < 100; i++ {
			k := "k" + strconv.Itoa(i)
			d.Put(k, i)
		}

		for i := 0; i < 100; i++ {
			key := "k" + strconv.Itoa(i)
			val, ok := d.Get(key)
			if ok {
				intVal, _ := val.(int)
				assert.Equal(t, i, intVal)
			} else {
				assert.Error(t, errors.New("put test failed: expected true, actual: false"))
			}

			ret := d.Remove(key)
			assert.Equal(t, 1, ret)
			_, ok = d.Get(key)
			assert.Equal(t, false, ok)
			ret = d.Remove(key)
			assert.Equal(t, 0, ret)
		}
	})

	t.Run("remove tail node", func(t *testing.T) {
		d := MakeConcurrent(0)
		for i := 0; i < 100; i++ {
			key := "k" + strconv.Itoa(i)
			d.Put(key, i)
		}
		for i := 9; i >= 0; i-- {
			key := "k" + strconv.Itoa(i)
			val, ok := d.Get(key)
			if ok {
				intVal, _ := val.(int)
				assert.Equal(t, i, intVal)
			} else {
				assert.Error(t, errors.New("put test failed: expected true, actual: false"))
			}
			ret := d.Remove(key)
			assert.Equal(t, 1, ret)
			_, ok = d.Get(key)
			assert.Equal(t, false, ok)
			ret = d.Remove(key)
			assert.Equal(t, 0, ret)
		}
	})

	t.Run("remove middle node", func(t *testing.T) {
		d := MakeConcurrent(0)
		d.Put("head", 0)
		for i := 0; i < 10; i++ {
			key := "k" + strconv.Itoa(i)
			d.Put(key, i)
		}
		d.Put("tail", 0)
		for i := 9; i >= 0; i-- {
			key := "k" + strconv.Itoa(i)
			val, ok := d.Get(key)
			if ok {
				assert.Equal(t, i, val)
			} else {
				assert.Error(t, errors.New("put test failed: expected true, actual: false"))
			}
			ret := d.Remove(key)
			assert.Equal(t, 1, ret)
			_, ok = d.Get(key)
			assert.Equal(t, false, ok)
			ret = d.Remove(key)
			assert.Equal(t, 0, ret)
		}
	})
}

func TestConcurrentDictForEach(t *testing.T) {
	d := MakeConcurrent(0)
	size := 100
	for i := 0; i < size; i++ {
		key := "k" + strconv.Itoa(i)
		d.Put(key, i)
	}
	i := 0
	d.ForEach(func(key string, val interface{}) bool {
		intVal, _ := val.(int)
		expectedKey := "k" + strconv.Itoa(intVal)
		assert.Equal(t, key, expectedKey)
		i++
		return true
	})
	assert.Equal(t, size, i)
}

func TestConcurrentDictRandomKeys(t *testing.T) {
	d := MakeConcurrent(0)
	count := 100
	for i := 0; i < count; i++ {
		key := "k" + strconv.Itoa(i)
		d.Put(key, i)
	}
	fetchSize := 10
	result := d.RandomKeys(fetchSize)
	assert.Equal(t, fetchSize, len(result))
	result = d.RandomDistinctKeys(fetchSize)
	distinct := make(map[string]struct{})
	for _, key := range result {
		distinct[key] = struct{}{}
	}
	assert.Equal(t, fetchSize, len(result))
	assert.Equal(t, len(distinct), len(result))
}

func TestConcurrentDictKeys(t *testing.T) {
	d := MakeConcurrent(0)
	size := 10
	for i := 0; i < size; i++ {
		d.Put(utils.RandString(5), utils.RandString(5))
	}
	assert.Equal(t, size, len(d.Keys()))
}

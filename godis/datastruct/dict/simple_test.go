package dict

import (
	"testing"

	"github.com/Muskchen/toolkits/godis/lib/utils"
	"github.com/stretchr/testify/assert"
)

func TestSimpleDictKeys(t *testing.T) {
	d := MakeSimple()
	size := 10
	for i := 0; i < size; i++ {
		d.Put(utils.RandString(5), utils.RandString(5))
	}
	assert.Equal(t, size, len(d.Keys()))
}

func TestSimpleDictPutIfExists(t *testing.T) {
	d := MakeSimple()
	key := utils.RandString(5)
	val := key + "1"
	ret := d.PutIfExists(key, val)
	assert.Equal(t, 0, ret)
	d.Put(key, val)
	val = key + "2"
	ret = d.PutIfExists(key, val)
	assert.Equal(t, 1, ret)
	v, _ := d.Get(key)
	assert.Equal(t, val, v)
}

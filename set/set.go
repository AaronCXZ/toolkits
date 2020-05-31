package set

import (
	"bytes"
	"fmt"
)

type Set interface {
	Add(e interface{}) bool
	Adds(items ...interface{})
	Remove(e interface{})
	Removes(items ...interface{})
	Clear()
	Contains(e interface{}) bool
	Len() int
	Same(other Set) bool
	Elements() []interface{}
	String() string
}

type HashSet struct {
	m map[interface{}]bool
}

func NewSet(items ...interface{}) Set {
	s := &HashSet{}
	s.m = make(map[interface{}]bool)
	for _, item := range items {
		s.Add(item)
	}
	return s
}

// 添加
func (set *HashSet) Add(e interface{}) bool {
	if !set.m[e] {
		set.m[e] = true
		return true
	}
	return false
}

func (set *HashSet) Adds(items ...interface{}) {
	for _, item := range items {
		set.Add(item)
	}
}

// 删除
func (set *HashSet) Remove(e interface{}) {
	if set.Contains(e) {
		delete(set.m, e)
	}
}

func (set *HashSet) Removes(items ...interface{}) {
	for _, item := range items {
		set.Remove(item)
	}
}

// 清空
func (set *HashSet) Clear() {
	set.m = make(map[interface{}]bool)
}

// 元素是否存在
func (set *HashSet) Contains(e interface{}) bool {
	return set.m[e]
}

// 长度
func (set *HashSet) Len() int {
	return len(set.m)
}

// 判断是否相同
func (set *HashSet) Same(other Set) bool {
	if other == nil {
		return false
	}
	if set.Len() != other.Len() {
		return false
	}
	for key := range set.m {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

// 获取所有元素
func (set *HashSet) Elements() []interface{} {
	initialLen := len(set.m)
	snapshot := make([]interface{}, initialLen)
	actualLen := 0
	for key := range set.m {
		if actualLen < initialLen {
			snapshot[actualLen] = key
		} else {
			snapshot = append(snapshot, key)
		}
		actualLen++
	}
	if actualLen < initialLen {
		snapshot = snapshot[:actualLen]
	}
	return snapshot
}

func (set *HashSet) String() string {
	var buf bytes.Buffer
	buf.WriteString("Set{")
	first := true
	for key := range set.m {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%v", key))
	}
	buf.WriteString("}")
	return buf.String()
}

// 超集
func IsSuperset(set, other Set) bool {
	if other == nil {
		return false
	}
	oneLen := set.Len()
	otherLen := other.Len()
	if oneLen == 0 || otherLen == 0 {
		return false
	}
	if oneLen > 0 && otherLen == 0 {
		return true
	}
	for _, v := range other.Elements() {
		if !set.Contains(v) {
			return false
		}
	}
	return true
}

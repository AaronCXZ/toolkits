package set

import (
	"bytes"
	"fmt"
	"sync"
)

type Set struct {
	m map[interface{}]struct{}
	sync.RWMutex
}

// 新建
func NewSet(items ...interface{}) *Set {
	s := &Set{}
	s.m = make(map[interface{}]struct{})
	s.Add(items...)
	return s
}

// 输出
func (s *Set) String() string {
	s.RLock()
	defer s.RUnlock()
	var buf bytes.Buffer

	buf.WriteString("set(")

	for key := range s.m {
		buf.WriteString(fmt.Sprintf("%v,", key))
	}
	buf.WriteString(")")
	return buf.String()
}

// 添加
func (s *Set) Add(items ...interface{}) {
	s.Lock()
	defer s.Unlock()
	for _, item := range items {
		s.m[item] = struct{}{}
	}
}

// 删除
func (s *Set) Del(items ...interface{}) {
	s.Lock()
	defer s.Unlock()
	for _, item := range items {
		delete(s.m, item)
	}
}

// 包含
func (s *Set) Contains(item interface{}) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

// 长度
func (s *Set) Len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.m)
}

// 清空
func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = make(map[interface{}]struct{})
}

// 是否为空
func (s *Set) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

// 转换为list
func (s *Set) List() []interface{} {
	s.RLock()
	defer s.RUnlock()
	var list []interface{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

// 相等
func (s *Set) Equal(other *Set) bool {
	s.RLock()
	defer s.RUnlock()
	if s.Len() != other.Len() {
		return false
	}
	for key := range s.m {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

// 子集
func (s *Set) IsSubset(other *Set) bool {
	s.RLock()
	defer s.RUnlock()
	if s.Len() > other.Len() {
		return false
	}
	for key := range s.m {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

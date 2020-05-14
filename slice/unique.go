package slice

import (
	"fmt"
	"reflect"
)

func Unique(data interface{}) interface{} {
	dataV := reflect.ValueOf(data)
	if dataV.Kind() != reflect.Slice && dataV.Kind() != reflect.Array {
		return data
	}
	m := make(map[interface{}]bool)
	new := reflect.MakeSlice(dataV.Type(), 0, dataV.Len())
	for i := 0; i < dataV.Len(); i++ {
		iVal := dataV.Index(i)
		if _, ok := m[iVal.Interface()]; !ok {
			new = reflect.Append(new, dataV.Index(i))
			m[iVal.Interface()] = true
		}
	}
	return new.Interface()
}

func UniqueOrigial(data interface{}) {
	dataV := reflect.ValueOf(data)
	if dataV.Kind() != reflect.Ptr {
		fmt.Println("输入的数据不是指针类型")
		return
	}
	tmpData := Unique(dataV.Elem().Interface())
	tmpDataV := reflect.ValueOf(tmpData)

	dataV.Elem().Set(tmpDataV)
}

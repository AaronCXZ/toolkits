// 单例模式
/*
	使用懒惰模式的单例模式，使用双重检查加锁保证线程安全
*/

package singleton

import "sync"

// 单例模式类
type Singleton struct {
}

var (
	singleton *Singleton
	once      sync.Once
)

// 用于获取单例模对象
func GetInstance() *Singleton {
	once.Do(func() {
		singleton = &Singleton{}
	})
	return singleton
}

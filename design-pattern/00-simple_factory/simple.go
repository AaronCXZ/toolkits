// 简单工厂模式
package implefacfory

import "fmt"

type API interface {
	Say(name string) string
}

// NewXXX函数返回接口时就是简单工厂模式
func NewAPI(t int) API {
	if t == 1 {
		return &hiAPI{}
	} else if t == 2 {
		return &HelloAPI{}
	}
	return nil
}

type hiAPI struct {
}

func (h *hiAPI) Say(name string) string {
	return fmt.Sprintf("Hi, %s", name)
}

type HelloAPI struct {
}

func (h *HelloAPI) Say(name string) string {
	return fmt.Sprintf("Hello, %s", name)
}

// 代理模式
/*
代理模式用于延迟处理操作或者在进行实际操作前后进行其它处理。
代理模式的常见用法有
	虚代理
	COW代理
	远程代理
	保护代理
	Cache 代理
	防火墙代理
	同步代理
	智能指引
	等。。。
*/

package proxy

type Subject interface {
	Do() string
}

type RealSubject struct {
}

func (RealSubject) Do() string {
	return "real"
}

type Proxy struct {
	real RealSubject
}

func (p Proxy) Do() string {
	var res string
	res += "pre:"
	res += p.real.Do()
	res += ":after"
	return res
}

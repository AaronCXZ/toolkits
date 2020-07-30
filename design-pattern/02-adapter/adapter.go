// 适配器模式
/*
	适配器模式用于转换一种接口适配另一种接口。
	实际使用中Adaptee一般为接口，并且使用工厂函数生成实例。
	在Adapter中匿名组合Adaptee接口，所以Adapter类也拥有SpecificRequest实例方法，又因为Go语言中非入侵式接口特征，其实Adapter也适配Adaptee接口。
*/
package adapter

// Target是适配的目标接口
type Target interface {
	Request() string
}

// Adaptee是被适配的目标接口
type Adaptee interface {
	SpecificRequest() string
}

// 被适配接口的工厂函数
func NewAdaptee() Adaptee {
	return &adapteeTmpl{}
}

// adapterTmpl是被适配的目标类
type adapteeTmpl struct {
}

// 目标类方法的实现
func (a *adapteeTmpl) SpecificRequest() string {
	return "adaptee method"
}

// Adaptee的工厂函数
func Newadapter(adaptee Adaptee) Target {
	return &adapter{Adaptee: adaptee}
}

// adapter是转换Adaptee为Target接口的适配器
type adapter struct {
	Adaptee
}

// Request实现Target接口
func (a *adapter) Request() string {
	return a.SpecificRequest()
}

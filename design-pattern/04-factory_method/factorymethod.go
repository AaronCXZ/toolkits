//  工厂方法模式
/*
	工厂方法模式使用子类的方式延迟生成对象到子类中实现。
	Go中不存在继承 所以使用匿名组合来实现
*/
package factorymethod

// 被封装的实际接口
type Operator interface {
	SetA(int)
	SetB(int)
	Result() int
}

// 工厂接口
type OperatorFactory interface {
	Create() Operator
}

// Operator接口实现的基类，封装共用方法
type OperatorBase struct {
	a, b int
}

func (o *OperatorBase) SetA(i int) {
	o.a = i
}

func (o *OperatorBase) SetB(i int) {
	o.b = i
}

// PlusOperator的工厂类
type PlusOperatorFactory struct {
}

func (p PlusOperatorFactory) Create() Operator {
	return &PlusOperator{
		OperatorBase: &OperatorBase{},
	}
}

//Operator的实际加法实现
type PlusOperator struct {
	*OperatorBase
}

func (o PlusOperator) Result() int {
	return o.a + o.b
}

// MinusOperator的工厂类
type MinusOperatorFactory struct {
}

func (m MinusOperatorFactory) Create() Operator {
	return &MinusOperator{
		OperatorBase: &OperatorBase{},
	}
}

// Operator的实际减法实现
type MinusOperator struct {
	*OperatorBase
}

func (o MinusOperator) Result() int {
	return o.a - o.b
}

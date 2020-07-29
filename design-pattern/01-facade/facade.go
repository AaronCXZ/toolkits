// 外观模式
/*
	同时暴露a和b两个module的NewXXX和interface，其它代码如果需要使用细节功能时可以直接调用。
*/

package facade

import "fmt"

type API interface {
	Test() string
}

func NewAPI() API {
	return &apiTmpl{
		a: NewAModuleAPI(),
		b: NewBModuleAPI(),
	}
}

type apiTmpl struct {
	a AModuleAPI
	b BModuleAPI
}

func (a *apiTmpl) Test() string {
	aRet := a.a.TestA()
	bRet := a.b.TestB()
	return fmt.Sprintf("%s\n%s", aRet, bRet)
}

type AModuleAPI interface {
	TestA() string
}

func NewAModuleAPI() AModuleAPI {
	return &aModuleAPI{}
}

type aModuleAPI struct {
}

func (b *aModuleAPI) TestA() string {
	return "A module running"
}

type BModuleAPI interface {
	TestB() string
}

func NewBModuleAPI() BModuleAPI {
	return &bModuleTmpl{}
}

type bModuleTmpl struct {
}

func (b *bModuleTmpl) TestB() string {
	return "B module running"
}

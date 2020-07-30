// 抽象工厂模式
/*
	抽象工厂模式用于生成产品族的工厂，所生成的对象是有关联的。
	如果抽象工厂退化成生成的对象无关联则成为工厂函数模式。
	比如本例子中使用RDB和XML存储订单信息，抽象工厂分别能生成相关的主订单信息和订单详情信息。 如果业务逻辑中需要替换使用的时候只需要改动工厂函数相关的类就能替换使用不同的存储方式了。
*/

package abstractfactory

import "fmt"

// 订单主记录
type OrderMainDAO interface {
	SaveOrderMain()
}

// 订单详情记录
type OrderDetailDAO interface {
	SaveOrderDetail()
}

// 抽象模式工厂接口
type DAOFactory interface {
	CreateOrderMainDAO() OrderMainDAO
	CreateOrderDetailDAO() OrderDetailDAO
}

// 关系型数据库的OrderMainDAO实现
type RDBMainDAO struct {
}

func (R *RDBMainDAO) SaveOrderMain() {
	fmt.Print("rdb main save\n")
}

// 关系型数据库的OrderDetailDAO实现
type RDBDetailDAO struct {
}

func (R *RDBDetailDAO) SaveOrderDetail() {
	fmt.Print("rdb detail save\n")
}

// RDB抽象工厂实现
type RDBDAOFactory struct {
}

func (R *RDBDAOFactory) CreateOrderMainDAO() OrderMainDAO {
	return &RDBMainDAO{}
}

func (R *RDBDAOFactory) CreateOrderDetailDAO() OrderDetailDAO {
	return &RDBDetailDAO{}
}

// XML存储
type XMLMainDAO struct {
}

func (X *XMLMainDAO) SaveOrderMain() {
	fmt.Print("xml main save\n")
}

// XMl存储
type XMLDetailDAO struct {
}

func (X *XMLDetailDAO) SaveOrderDetail() {
	fmt.Print("xml detail dave\n")
}

// RDB抽象工厂实现
type XMLDAOFactory struct {
}

func (X *XMLDAOFactory) CreateOrderMainDAO() OrderMainDAO {
	return &XMLMainDAO{}
}

func (X *XMLDAOFactory) CreateOrderDetailDAO() OrderDetailDAO {
	return &XMLDetailDAO{}
}

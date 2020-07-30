package abstractfactory

import "testing"

var (
	rdbwant = "rdb main save\nrdb detail save\n"
	xmlwant = "xml main save\nxml detail dave\n"
)

func getMainAndDetail(factory DAOFactory) {
	factory.CreateOrderMainDAO().SaveOrderMain()
	factory.CreateOrderDetailDAO().SaveOrderDetail()
}

func ExampleRdbFactory() {
	var factory DAOFactory
	factory = &RDBDAOFactory{}
	getMainAndDetail(factory)
}

func ExampleXmlFactory() {
	var factory DAOFactory
	factory = &XMLDAOFactory{}
	getMainAndDetail(factory)
}

func TestFactory(t *testing.T) {
	t.Run("rdb", func(t *testing.T) {
		ExampleRdbFactory()
	})
	t.Run("xml", func(t *testing.T) {
		ExampleXmlFactory()
	})
}

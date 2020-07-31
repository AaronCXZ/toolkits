package strategy

import "testing"

func ExampleCash_Pay() {
	payment := NewPayment("ada", "", 123, &Cash{})
	payment.Pay()
}

func ExampleBank_Pay() {
	payment := NewPayment("Bob", "00002", 888, &Bank{})
	payment.Pay()
}

func TestPayment(t *testing.T) {
	ExampleCash_Pay()
	ExampleBank_Pay()
}

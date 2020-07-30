package iterator

import "testing"

func ExampleIterator() {
	var aggrsgate Aggregate
	aggrsgate = NewNumbers(1, 10)
	IteratorPrint(aggrsgate.Iterator())
}

func TestIterator(t *testing.T) {
	ExampleIterator()
}

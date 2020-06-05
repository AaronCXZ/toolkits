package snowflake

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {
	node, err := NewWorker(1)
	if err != nil {
		assert.Error(t, err)
	}
	for {
		fmt.Println(node.GetId())
	}
}

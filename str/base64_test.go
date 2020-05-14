package str

import (
	"fmt"
	"testing"
)

func TestBase64(t *testing.T) {
	s := "abc.1234%"
	t.Run("encode", func(t *testing.T) {
		encode := Base64Encode(s)
		fmt.Println(encode)
	})
	t.Run("decode", func(t *testing.T) {
		str := "YWJjLjEyMzQl"
		decode, err := Base64Decode(str)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(decode)
	})
}

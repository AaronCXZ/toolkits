package str

import (
	"fmt"
	"testing"
)

func TestMD5(t *testing.T) {
	s := "abc.1234%"
	md5, err := MD5(s)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(md5)

}

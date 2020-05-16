package str

import (
	"fmt"
	"testing"
)

func TestMD5(t *testing.T) {
	s := "abc.1234%"
	t.Run("md5", func(t *testing.T) {
		md5, err := MD5(s)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(md5)
	})
	t.Run("sha512", func(t *testing.T) {
		sha, err := Sha512(s)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(sha)
	})

	t.Run("sha256", func(t *testing.T) {
		sha, err := Sha256(s)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(sha)

	})
}

package str

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func MD5(str string) (string, error) {
	h := md5.New()
	if _, err := io.WriteString(h, str); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

package str

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
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

func Sha512(str string) (string, error) {
	h := sha512.New()
	if _, err := io.WriteString(h, str); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func Sha256(str string) (string, error) {
	h := sha256.New()

	if _, err := io.WriteString(h, str); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

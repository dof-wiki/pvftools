package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

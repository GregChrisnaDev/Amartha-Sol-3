package common

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Hasher(param string) string {
	hasher := md5.New()

	hasher.Write([]byte(param))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
)

func Hash(b []byte) string {
	h := sha256.New()
	h.Write(b)
	h.Write(h.Sum(nil))
	return hex.EncodeToString(h.Sum(nil))
}

func GetFileHash(filename string) (s string, err error) {
	var (
		f os.FileInfo
		b []byte
	)
	if f, err = os.Stat(filename); err != nil {
		return
	}
	if f.IsDir() {
		return "", errors.New("this cmd is dir")
	}
	// hash 文件
	if b, err = os.ReadFile(filename); err != nil {
		return
	}
	s = Hash(b)
	return
}

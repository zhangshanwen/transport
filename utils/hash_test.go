package utils

import "testing"

func TestHash(t *testing.T) {
	a := "111111"
	t.Log(Hash([]byte(a)))
}

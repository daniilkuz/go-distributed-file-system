package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCopyEncryptDecrypt(t *testing.T) {
	src := bytes.NewReader([]byte("Good moring!"))
	dst := new(bytes.Buffer)
	key := newEncryptionKey()
	_, err := copyEncrypt(key, src, dst)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(dst.Bytes())
}

package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCopyEnrypt(t *testing.T) {
	src := bytes.NewReader([]byte("Good moring!"))
	dst := new(bytes.Buffer)
	key := newEncryptionKey()
	_, err := coppyEncrypt(key, src, dst)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(dst.Bytes())
}

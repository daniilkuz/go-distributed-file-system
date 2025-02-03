package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCopyEncryptDecrypt(t *testing.T) {
	payload := "Good moring!"
	src := bytes.NewReader([]byte(payload))
	dst := new(bytes.Buffer)
	key := newEncryptionKey()
	_, err := copyEncrypt(key, src, dst)
	if err != nil {
		t.Error(err)
	}

	// fmt.Println(dst.String())

	fmt.Println(len(payload))
	fmt.Println(len(dst.String()))

	out := new(bytes.Buffer)
	if _, err := copyDecrypt(key, dst, out); err != nil {
		t.Error(err)
	}

	if out.String() != payload {
		t.Error("Decryption failed!")
	}

	fmt.Println(out.String())
}

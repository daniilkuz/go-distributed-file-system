package main

import (
	"bytes"
	"testing"
)

func TestCopyEnrypt(t *testing.T) {
	src := bytes.NewReader([]byte("Good moring!"))
	dst := new(bytes.Buffer)
}

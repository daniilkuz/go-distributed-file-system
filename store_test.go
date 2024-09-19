package main

import (
	"bytes"
	"testing"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: DefaultPathTransofrmFunc,
	}
	s := NewStore(opts)

	data := bytes.NewBuffer([]byte("random jpg image"))
	if err := s.writeStream("key", data); err != nil {
		t.Error(err)
	}
}

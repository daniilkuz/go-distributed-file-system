package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "something random to test PathTRansformFunc"
	pathKey := CASPathTransformFunc(key)
	fmt.Println(pathKey)
	expectedOriginalKey := "667d6978205dda2be5c1d0562e8e546e5c793c89"
	expectedPathname := "667d6/97820/5dda2/be5c1/d0562/e8e54/6e5c7/93c89"
	if pathKey.Pathname != expectedPathname {
		t.Errorf("got %s, but expected %s", pathKey.Pathname, expectedPathname)
	}

	if pathKey.Original != expectedOriginalKey {
		t.Errorf("got %s, but expected %s", pathKey.Original, expectedOriginalKey)
	}

}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewBuffer([]byte("random jpg image"))
	if err := s.writeStream("key", data); err != nil {
		t.Error(err)
	}
}

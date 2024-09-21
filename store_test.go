package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "something random to test PathTRansformFunc"
	pathname := CASPathTransformFunc(key)
	fmt.Println(pathname)
	expectedPathname := "667d6/97820/5dda2/be5c1/d0562/e8e54/6e5c7/93c89"
	if pathname != expectedPathname {
		t.Errorf("got %s, but expected %s", pathname, expectedPathname)
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

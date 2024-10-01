package main

import (
	"bytes"
	"fmt"
	"io"
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

	if pathKey.Filename != expectedOriginalKey {
		t.Errorf("got %s, but expected %s", pathKey.Filename, expectedOriginalKey)
	}

}

func TestStoreDeleteKey(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "key for jpg"
	data := []byte("random jpg image")
	if err := s.writeStream(key, bytes.NewBuffer(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}

}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "key for jpg"

	data := []byte("random jpg image")
	if err := s.writeStream(key, bytes.NewBuffer(data)); err != nil {
		t.Error(err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)
	if string(b) != string(data) {
		t.Errorf("got %s, but expected %s", b, data)
	}

	s.Delete(key)
}

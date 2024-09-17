package main

import (
	"fmt"
	"testing"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: DefaultPathTransofrmFunc,
	}
	s := NewStore(opts)
	fmt.Println("%+v\n", s)
}

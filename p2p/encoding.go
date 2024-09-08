package p2p

import (
	"encoding/gob"
	"fmt"
	"io"
)

type Decoder interface {
	Decode(io.Reader, any) error
}

type GOBDecoder struct{}

func (dec GOBDecoder) Decode(r io.Reader, v any) error {
	return gob.NewDecoder(r).Decode(v)
}

type NOPDecoder struct{}

func (dec NOPDecoder) Decode(r io.Reader, msg Message) error {
	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	fmt.Println(string(buf[:n]))

	return nil
}

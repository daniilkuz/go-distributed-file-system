package main

import (
	"fmt"
	"log"

	"github.com/daniilkuz/go-distributed-file-system/p2p"
)

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		HandshakeFunc: p2p.NOPHandshakeFunc,
		ListenAddr:    ":3000",
		Decoder:       p2p.DefaultDecoder{},
	}

	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}

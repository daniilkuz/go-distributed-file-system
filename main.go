package main

import (
	"log"

	"github.com/daniilkuz/go-distributed-file-system/p2p"
)

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		HandshakeFunc: p2p.NOPHandshakeFunc,
		ListenAddr:    ":3000",
		Decoder:       p2p.GOBDecoder{},
	}

	tr := p2p.NewTCPTransport(tcpOpts)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}

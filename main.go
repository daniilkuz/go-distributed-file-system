package main

import (
	"log"

	"github.com/daniilkuz/go-distributed-file-system/p2p"
)

func main() {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    "3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)
	fileServerOpts := FileServerOpts{
		ListenAddr:        "3000",
		StoreageRoot:      "3000_netweork",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         *tcpTransport,
	}
	s := NewFileServer(fileServerOpts)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

	select {}
}

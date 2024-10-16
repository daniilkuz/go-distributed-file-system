package main

import (
	"log"
	"time"

	"github.com/daniilkuz/go-distributed-file-system/p2p"
)

func main() {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)
	fileServerOpts := FileServerOpts{
		StoreageRoot:      "3000_netweork",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}
	s := NewFileServer(fileServerOpts)

	go func() {
		time.Sleep(time.Second)
		s.Stop()
	}()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

	select {}
}

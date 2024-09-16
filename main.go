package main

import (
	"fmt"
	"log"

	"github.com/daniilkuz/go-distributed-file-system/p2p"
)

func OnPeer(peer p2p.Peer) error {
	// fmt.Println("doing some logic with the peer outside of TCPTransport")
	peer.Close()
	return nil
}

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		HandshakeFunc: p2p.NOPHandshakeFunc,
		ListenAddr:    ":3000",
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
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

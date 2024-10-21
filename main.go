package main

import (
	"log"

	"github.com/daniilkuz/go-distributed-file-system/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)
	fileServerOpts := FileServerOpts{
		StoreageRoot:      listenAddr + "_netweork",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}
	return NewFileServer(fileServerOpts)

}

func main() {
	// tcpTransportOpts := p2p.TCPTransportOpts{
	// 	ListenAddr:    ":3000",
	// 	HandshakeFunc: p2p.NOPHandshakeFunc,
	// 	Decoder:       p2p.DefaultDecoder{},
	// }
	// tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)
	// fileServerOpts := FileServerOpts{
	// 	StoreageRoot:      "3000_netweork",
	// 	PathTransformFunc: CASPathTransformFunc,
	// 	Transport:         tcpTransport,
	// 	BootstrapNodes:    []string{":4000"},
	// }
	// s := NewFileServer(fileServerOpts)

	// go func() {
	// 	time.Sleep(time.Second)
	// 	s.Stop()
	// }()

	// if err := s.Start(); err != nil {
	// 	log.Fatal(err)
	// }

	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	s2.Start()

	// select {}
}

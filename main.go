package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

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
		EncKey:            newEncryptionKey(),
		StoreageRoot:      listenAddr + "_netweork",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}
	s := NewFileServer(fileServerOpts)
	tcpTransport.OnPeer = s.OnPeer
	return s

}

func main() {

	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(4 * time.Second)

	go s2.Start()
	time.Sleep(4 * time.Second)

	// data := bytes.NewReader([]byte("Something to say"))
	// s2.Store("picture.jpg", data)
	// time.Sleep(5 * time.Millisecond)

	r, err := s2.Get("picture.jpg")
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	// select {}
}

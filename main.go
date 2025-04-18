package main

import (
	"bytes"
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
	s3 := makeServer(":5000", ":3000", ":4000")

	// go func() {
	go func() { log.Fatal(s1.Start()) }()
	// time.Sleep(time.Second * 1)
	time.Sleep(500 * time.Millisecond)
	go func() { log.Fatal(s2.Start()) }()
	// }()

	time.Sleep(4 * time.Second)

	go s3.Start()
	time.Sleep(4 * time.Second)

	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("pricture_%d.png", i)
		data := bytes.NewReader([]byte("Something to say"))
		s3.Store(key, data)
		time.Sleep(5 * time.Millisecond)

		if err := s3.store.Delete(key); err != nil {
			log.Fatal(err)
		}

		r, err := s3.Get(key)
		if err != nil {
			log.Fatal(err)
		}

		b, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(b))
	}

	// select {}
}

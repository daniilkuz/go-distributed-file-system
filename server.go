package main

import (
	"fmt"

	"github.com/daniilkuz/go-distributed-file-system/p2p"
)

type FileServerOpts struct {
	ListenAddr        string
	StoreageRoot      string
	PathTransformFunc PathTransformFunc
	Transport         p2p.TCPTransport
}

type FileServer struct {
	FileServerOpts

	store  *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StoreageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}

	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
	}
}

func (s *FileServer) loop() {
	for {
		select {
		case msg := <-s.Transport.Consume():
			fmt.Println(msg)
		}
	}
}

func (s *FileServerOpts) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	return nil
}

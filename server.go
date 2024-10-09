package main

import "github.com/daniilkuz/go-distributed-file-system/p2p"

type FileServerOpts struct {
	ListenAddr        string
	StoreageRoot      string
	PathTransformFunc PathTransformFunc
	Transport         p2p.TCPTransport
}

type FileServer struct {
	FileServerOpts

	store *Store
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StoreageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}

	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
	}
}

func (s *FileServerOpts) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	return nil
}

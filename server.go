package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/daniilkuz/go-distributed-file-system/p2p"
)

type FileServerOpts struct {
	// ListenAddr        string
	StoreageRoot      string
	PathTransformFunc PathTransformFunc
	Transport         *p2p.TCPTransport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts
	peerLock sync.Mutex
	peers    map[string]p2p.Peer

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
		peers:          make(map[string]p2p.Peer),
	}
}

type Message struct {
	// From    string
	Payload any
}

type MessageStoreFile struct {
	Key  string
	Size int64
}

// type DataMessage struct {
// 	Key  string
// 	Data []byte
// }

func (s *FileServer) broadcast(msg *Message) error {

	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(msg)
}

func (s *FileServer) StoreData(key string, r io.Reader) error {

	buf := new(bytes.Buffer)
	tee := io.TeeReader(r, buf)

	size, err := s.store.Write(key, tee)
	if err != nil {
		return err
	}

	msg := Message{
		Payload: MessageStoreFile{
			Key:  key,
			Size: size,
		},
	}

	msgBuf := new(bytes.Buffer)
	if err := gob.NewEncoder(msgBuf).Encode(msg); err != nil {
		return err
	}

	for _, peer := range s.peers {
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}

	time.Sleep(time.Second * 3)

	// payload := []byte("super large file")

	for _, peer := range s.peers {
		// n, err := io.Copy(peer, bytes.NewReader(payload))
		n, err := io.Copy(peer, r)
		if err != nil {
			return err
		}
		fmt.Println("received and written bytes to disk: ", n)
		// if err := peer.Send(payload); err != nil {
		// 	return err
		// }
	}

	return nil

	// buf := new(bytes.Buffer)
	// tee := io.TeeReader(r, buf)

	// if err := s.store.Write(key, tee); err != nil {
	// 	return err
	// }

	// p := &DataMessage{
	// 	Key:  key,
	// 	Data: buf.Bytes(),
	// }
	// // fmt.Printf("written %d bytes\n", n)
	// fmt.Println(buf.Bytes())
	// return s.broadcast(&Message{
	// 	From:    "todo",
	// 	Payload: p,
	// })
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[p.RemoteAddr().String()] = p
	log.Printf("connected with remote %s", p.RemoteAddr())
	return nil
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped due to user quit")
		s.Transport.Close()
	}()
	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Println(err)
			}

			if err := s.handleMessage(rpc.From, &msg); err != nil {
				log.Println(err)
				return
			}

			// fmt.Printf("%+v\n", msg.Payload)
			// // fmt.Printf("recv: %s\n", string(msg.Payload.([]byte)))

			// peer, ok := s.peers[rpc.From]
			// if !ok {
			// 	panic("peer not found in peer map")
			// }

			// // fmt.Println(peer)
			// buf := make([]byte, 1000)
			// if _, err := peer.Read(buf); err != nil {
			// 	panic(err)
			// }

			// fmt.Printf("%s\n", string(buf))

			// peer.(*p2p.TCPPeer).Wg.Done()

			// if err := s.handleMessage(&m); err != nil {
			// 	log.Println(err)
			// }
			// fmt.Printf("%+v\n", p)
		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		// fmt.Printf("received data %+v\n", v)
		return s.handleMessageStoreFile(from, v)
	}
	return nil
}

func (s *FileServer) handleMessageStoreFile(from string, msg MessageStoreFile) error {
	peer, ok := s.peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) could not be found in the peer list", from)
	}

	if _, err := s.store.Write(msg.Key, io.LimitReader(peer, msg.Size)); err != nil {
		return err
	}

	peer.(*p2p.TCPPeer).Wg.Done()

	// fmt.Printf("recv store file msg: %+v\n", msg)
	return nil
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			fmt.Println("attempting to connect with remote: ", addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error: ", err)
				// panic(err)
			}
		}(addr)
	}
	return nil
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	s.bootstrapNetwork()
	s.loop()

	return nil
}

func init() {
	gob.Register(MessageStoreFile{})
}

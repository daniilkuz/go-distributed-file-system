package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

// TCPPeer represents the remote over a TCP established connections
type TCPPeer struct {
	// conn is the underlying connection of the peer
	net.Conn

	// if we dial a connection a conn => outbound=true
	// if we accept and retrieve a conn=>outbuond = false
	outbound bool

	wg *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		wg:       &sync.WaitGroup{},
	}
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}

// func (p *TCPPeer) RemoteAddr() net.Addr {
// 	return p.conn.RemoteAddr()
// }

// func (p *TCPPeer) Close() error {
// 	return p.conn.Close()
// }

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC

	// mu    sync.RWMutex
	// peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
	}
}

func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	go t.handleConn(conn, true)
	return nil
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()
	log.Printf("TCP transport listening on port: %s\n", t.ListenAddr)

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			fmt.Printf("TCP accpet error: %s\n", err)
		}

		fmt.Printf("new incomming connection %+v\n", conn)
		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {

	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)

	if err := t.HandshakeFunc(peer); err != nil {
		// conn.Close()
		// fmt.Printf("TCP handshake error: %s\n", err)
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	rpc := RPC{}
	for {
		err = t.Decoder.Decode(conn, &rpc)
		// fmt.Println(reflect.TypeOf(err))
		// panic(err)
		// if err == net.ErrClosed {
		// 	return
		// }

		if err != nil {
			fmt.Printf("TCP error %s\n", err)
			return
		}

		rpc.From = conn.RemoteAddr().String()
		peer.wg.Add(1)
		fmt.Println("waiting until stream is done")
		t.rpcch <- rpc
		peer.wg.Wait()
		fmt.Println("stream is done, continuing normal loop")
		// fmt.Printf("message: %+v\n", rpc)
	}

}

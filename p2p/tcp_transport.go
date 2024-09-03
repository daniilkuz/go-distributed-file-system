package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote over a TCP established connections
type TCPPeer struct {
	// conn is the underlying connection of the peer
	conn net.Conn

	// if we dial a connection a conn => outbound=true
	// if we accept and retrieve a conn=>outbuond = false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	listenAddr string
	listener   net.Listener
	shakeHands HandshakeFunc
	decoder    Decoder

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddr: listenAddr,
		shakeHands: NOPHandshakeFunc,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accpet error: %s\n", err)
		}

		fmt.Printf("new incomming connection %+v\n", conn)
		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {

	peer := NewTCPPeer(conn, true)

	if err := t.shakeHands(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", err)
		return
	}

	// buf := new(bytes.Buffer)
	msg := &Temp{}

	for {
		if err := t.decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error %s\n", err)
			continue
		}
	}

}

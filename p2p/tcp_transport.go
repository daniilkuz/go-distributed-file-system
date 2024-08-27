package p2p

import (
	"net"
	"sync"
)

type TCPTransport struct {
	listenAddr string
	listener   net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddr: listenAddr,
	}
}

func (t *TCPTransport) listenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}
}

func (t *TCPTransport) acceptLoop() {

}

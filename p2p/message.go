package p2p

const (
	IncommingMessage = 0x1
	IncommingStream  = 0x2
)

type RPC struct {
	// From    net.Addr
	From    string
	Payload []byte
	stream  bool
}

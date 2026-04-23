package p2p

import "net"

// Peer is remote node
type Peer interface {
	RemoteAddr() net.Addr
	Close() error
}

// Handles communication between nodes in the network.
// Can be TCP, UDP, websockets etc...

type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
	Dial(string) error
}

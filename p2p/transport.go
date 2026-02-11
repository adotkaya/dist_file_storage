package p2p

// Peer is remote node
type Peer interface {
	Close() error
}

// Handles communication between nodes in the network.
// Can be TCP, UDP, websockets etc...

type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}

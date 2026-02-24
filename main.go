package main

import (
	"dist_file_storage/p2p"
	"log"
	"time"
)

func main() {
	tcptransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		//TODO OnPeer:        OnPeer,
	}
	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)
	fileServerOpts := FileServerOpts{
		StorageRoot:       "/tmp_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}
	fs := NewFileServer(fileServerOpts)
	go func() {
		time.Sleep(time.Second * 3)
		fs.Quit()
	}()
	if err := fs.Start(); err != nil {
		log.Fatal(err)
	}

}

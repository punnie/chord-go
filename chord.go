package main

import (
	"flag"
)

func main() {
	listenPtr := flag.String("listen", "127.0.0.1:5556", "address to listen on")
	joinPtr := flag.String("join", "address:port", "hostname to join to")
	flag.Parse()

	// TODO: autogenerate node id
	self := NewNode(*listenPtr)
	dht := NewDHT(self)

	if *joinPtr != "address:port" {
		node := NewNode(*joinPtr)
		dht.Join(node)
	}

	dht.Listen()
}

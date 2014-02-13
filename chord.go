package main

import (
	"flag"
  "fmt"
)

func main() {
	listenPtr := flag.String("listen", "127.0.0.1:5556", "address to listen on")
	joinPtr := flag.String("join", "address:port", "hostname to join to")
	flag.Parse()

	selfKey, err := NewKeyID().GenerateKeyID(*listenPtr)

  fmt.Printf("\x1b[32mkey: %s \x1b[0m\n", selfKey.String())

	if err != nil {
		println("LOLWUT")
		panic(err)
	}

	self := NewNode(selfKey, *listenPtr)
	dht := NewDHT(self)

	if *joinPtr != "address:port" {
		node := NewNode(nil, *joinPtr)
		dht.Join(node)
	}

	dht.Listen()
}

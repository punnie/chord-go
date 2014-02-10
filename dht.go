package main

import (
	"bufio"
	//"fmt"
	"net"
	"os"
	//"time"
)

const BITS = 160 // sha1

type DHT struct {
	self        *Node
	finger      []*Node
	predecessor *Node
	successor   *Node
	messages    chan *Message
}

func NewDHT(self *Node) *DHT {
	dht := DHT{
		self:        self,
		finger:      make([]*Node, BITS),
		predecessor: nil,
		successor:   self,
		messages:    make(chan *Message, 100),
	}

	for i, _ := range dht.finger {
		dht.finger[i] = self
	}

	for i := 0; i < 1; i++ { // TODO: implement worker number
		go dht.Worker()
	}

	return &dht
}

func (d *DHT) Store(object []byte) error {
	return nil
}

func (d *DHT) Retrieve(id int64) ([]byte, error) {
	return nil, nil
}

func (d *DHT) Join(n *Node) {
	err := n.Connect()

	if err != nil {
		panic(err)
	}

	d.findSuccessor(d.self) // blocks
}

func (d *DHT) Listen() {
	sock, err := net.Listen("tcp", d.self.Address())

	if err != nil {
		println("Error listening:", err)
		os.Exit(1)
	}

	println("Listening on", d.self.Address())

	for {
		conn, err := sock.Accept()

		if err != nil {
			println("Error accepting!")
		}

		go d.handleConnection(conn)
	}
}

func (d *DHT) handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		mesBuffer, err := readMessage(reader)

		if err != nil {
			conn.Close()
			return
		}

		mes, err := MessageDecode(mesBuffer)

		if err != nil {
			println("Error decoding message!")
			continue
		}

		println("message from:", conn.RemoteAddr().String())
		d.messages <- mes
	}
}

//
//
//
//
//

func (d *DHT) Worker() {
	println("dht worker started")

	for {
		m := <-d.messages
		println("message     :", m.String())
	}
}

//
//
//
//
//

func (d *DHT) findSuccessor(node *Node) *Node {
	if node.Id().elementOf(d.self.Id(), d.successor.Id()) {
		return d.successor
	} else {
		queryNode := d.closestPrecedingNode(node)
		println(queryNode)
		// send "findSuccessor" message to queryNode
	}

	return &Node{}
}

func (d *DHT) closestPrecedingNode(node *Node) *Node {
	for i := BITS; i > 0; i-- {
		if d.finger[i].Id().elementOf(d.self.Id(), node.Id()) {
			return d.finger[i]
		}
	}

	return d.self
}

func (d *DHT) stabilize() {
}

func (d DHT) notify(node Node) {
}

func (d *DHT) fixFingers() {
}

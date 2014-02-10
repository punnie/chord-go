package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type Node struct {
	id           *KeyID
	address      string
	conn         net.Conn
	writer       *bufio.Writer
	localInbound chan *Message
}

func NewNode(address string) *Node {
	println("new connection:", address)

	i := new(KeyID)
	err := i.GenerateNodeKeyID(address)

	if err != nil {
		panic(err)
	}

	n := &Node{
		id:           i,
		address:      address,
		localInbound: make(chan *Message),
	}

	fmt.Println("new node id hex:", i.String())
	return n
}

func FakeNode(hex string) *Node {
	i := NewKeyID(hex)

	return &Node{
		id: i,
	}
}

func (n *Node) Accept(conn net.Conn, globalInbound chan<- *Message) {
	n.conn = conn

	go n.handleConnection(globalInbound)
}

func (n Node) Id() *KeyID {
	return n.id
}

func (n Node) Address() string {
	return n.address
}

func (n *Node) Connect(globalInbound chan<- *Message) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", n.address)

	if err != nil {
		return err
	}

	n.conn, err = net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		return err
	}

	go n.handleConnection(globalInbound)

	return nil
}

func (n *Node) handleConnection(globalInbound chan<- *Message) {
	reader := bufio.NewReader(n.conn)

	for {
		mesBuffer, err := readMessage(reader)

		if err != nil {
			n.conn.Close()
			return
		}

		mes, err := MessageDecode(mesBuffer)

		if err != nil {
			println("Error decoding message!")
			continue
		}

		switch mes.Intent {
		case REPLY_SUCCESSOR:
			fallthrough
		case REPLY_PING:
			n.localInbound <- mes
		case REQUEST_SUCCESSOR:
			fallthrough
		case REQUEST_PING:
			mes.Sender = n
			globalInbound <- mes
		}
	}
}

func (n *Node) sendMessage(m *Message) (int, error) {
	// TODO: verify we have an active connection
	if n.writer == nil {
		n.writer = bufio.NewWriter(n.conn) // this is not threadsafe!
	}

	println("sending        :", m.String())

	buf := new(bytes.Buffer)
	payload, err := m.MessageEncode()

	if err != nil {
		return 0, err
	}

	size := uint32(len(payload))

	err = binary.Write(buf, binary.LittleEndian, size)

	if err != nil {
		return 0, err
	}

	w, err := buf.Write(payload)

	if err != nil {
		return 0, err
	}

	_, err = buf.WriteTo(n.writer)

	if err != nil {
		return 0, err
	}

	n.writer.Flush()

	return w, nil
}

func (n *Node) RequestSuccessor(node *Node) (*Node, error) {
	queryKey := node.Id()
	n.sendMessage(NewFindSuccessorMessage([]string{queryKey.String()}))

	reply := <-n.localInbound
	println("received reply :", reply.String())

	return &Node{}, nil
}

func (n *Node) ReplySuccessor(node *Node) error {
	reply := &Message{
		Intent:     REPLY_SUCCESSOR,
		Parameters: []string{node.Id().String(), node.Address()},
		Timestamp:  time.Now().UTC(),
	}

	w, err := n.sendMessage(reply)

	println("sending size   :", w)

	if err != nil {
		return err
	}

	return nil
}

func (n *Node) RequestPredecessor() (*Node, error) {
	return nil, nil
}

func (n *Node) SendPing() error {
	_, err := n.sendMessage(&Message{Intent: REQUEST_PING, Timestamp: time.Now().UTC()})

	if err != nil {
		return err
	}

	return nil
}

func (n *Node) ReplyPing() error {
	_, err := n.sendMessage(&Message{Intent: REPLY_PING, Timestamp: time.Now().UTC()})

	if err != nil {
		return err
	}

	return nil
}

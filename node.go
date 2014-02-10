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
		id:      i,
		address: address,
	}

	fmt.Println("new node id hex:", i.String())
	return n
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
      println(mes.String())
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
	writer := bufio.NewWriter(n.conn)
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

	_, err = buf.WriteTo(writer)

	if err != nil {
		return 0, err
	}

	writer.Flush()
	return w, nil
}

func (n *Node) GetSuccessor(node *Node) (*Node, error) {
	queryKey := node.Id()
	n.sendMessage(NewFindSuccessorMessage([]string{queryKey.String()}))

  reply := <-n.localInbound

  println("received reply:", reply)

	return &Node{}, nil
}

func (n *Node) GetPredecessor() (*Node, error) {
	return nil, nil
}

func (n *Node) SendPing() error {
  n.sendMessage(&Message{Intent: REQUEST_PING, Timestamp: time.Now().UTC()})

  return nil
}

func (n *Node) ReplyPing() error {
  n.sendMessage(&Message{Intent: REPLY_PING, Timestamp: time.Now().UTC()})

  return nil
}

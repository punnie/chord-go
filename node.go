package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type Node struct {
	id      *KeyID
	address string
	conn    net.Conn
}

func NewNode(address string) *Node {
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

func (n Node) Id() KeyID {
	return *n.id
}

func (n Node) Address() string {
	return n.address
}

func (n *Node) Connect() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", n.address)

	if err != nil {
		return err
	}

	n.conn, err = net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		return err
	}

	return nil
}

func (n *Node) SendMessage(m *Message) (int, error) {
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

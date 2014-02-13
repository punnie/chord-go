package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"net"
  "time"
  "errors"
)

type Node struct {
	id            *KeyID
	dht           *DHT
	globalAddress string
	localAddress  string
	conn          net.Conn
	reader        *bufio.Reader
	writer        *bufio.Writer
	localInbound  chan *Message
	localOutbound chan *Message
	quit          chan bool
}

func NewNode(id *KeyID, globalAddress string) *Node {
	return &Node{
		id:            id,
		globalAddress: globalAddress,
		localInbound:  make(chan *Message),
		localOutbound: make(chan *Message),
		quit:          make(chan bool),
	}
}

func (n *Node) Accept(d *DHT, conn net.Conn) {
	n.dht = d
	n.conn = conn

	go n.handleConnection()
}

func (n Node) Id() *KeyID {
	return n.id
}

func (n Node) GlobalAddress() string {
	return n.globalAddress
}

func (n *Node) Connect(d *DHT) error {
	n.dht = d

	tcpAddr, err := net.ResolveTCPAddr("tcp", n.globalAddress)

	if err != nil {
		return err
	}

	n.conn, err = net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		return err
	}

	go n.handleConnection()

	return nil
}

func (n *Node) Close() {
	n.quit <- true
}

func (n *Node) handleConnection() {
	reader := bufio.NewReader(n.conn)

	for {
		select {
		case <-n.quit:
			defer n.conn.Close()
			println("connection to node closed!")
			break

		default:
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

			if n.id == nil {
				senderHash := mes.SenderHash
				senderAddr := mes.SenderAddr

				nodeId := NewKeyID().SetHash(senderHash)
				n.id = nodeId
				n.globalAddress = senderAddr

				n.dht.nodePool.PutNode(n)
			}

			switch mes.Intent {
			case _M_REPLY_SUCCESSOR:
				fallthrough
			case _M_REPLY_PREDECESSOR:
				fallthrough
			case _M_REPLY_PING:
				n.localInbound <- mes

			case _M_REQUEST_SUCCESSOR:
				fallthrough
			case _M_REQUEST_PREDECESSOR:
				fallthrough
			case _M_NOTIFY:
				fallthrough
			case _M_REQUEST_PING:
				envelope := mes.NewEvelope()
				envelope.sender = n

				n.dht.globalInbound <- envelope
			}
		}

	}
}

func (n *Node) sendMessage(m Message) (int, error) {
	// TODO: verify we have an active connection
	if n.writer == nil {
		n.writer = bufio.NewWriter(n.conn) // this is not threadsafe!
	}

  //println("sending        :", m.String())

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

func (n *Node) RequestSuccessor(id *KeyID) (*Node, error) {
	requestMessage := n.dht.self.NewMessage(_M_REQUEST_SUCCESSOR, []string{id.String()})
	n.sendMessage(requestMessage)

	replyMessage := <-n.localInbound
	successorId := NewKeyID().SetHash(replyMessage.Parameters[0])
	successorAddr := replyMessage.Parameters[1]

	return n.dht.nodePool.GetNode(successorId, successorAddr), nil
}

func (n *Node) ReplySuccessor(node *Node) error {
	reply := n.dht.self.NewMessage(_M_REPLY_SUCCESSOR, []string{
		node.Id().String(),
		node.GlobalAddress(),
	})

	_, err := n.sendMessage(reply)

	if err != nil {
		return err
	}

	return nil
}

func (n *Node) RequestPredecessor() (*Node, error) {
	requestMessage := n.dht.self.NewMessage(_M_REQUEST_PREDECESSOR, []string{})
	n.sendMessage(requestMessage)

	replyMessage := <-n.localInbound

	if len(replyMessage.Parameters) > 0 {
		predecessorId := NewKeyID().SetHash(replyMessage.Parameters[0])
		predecessorAddr := replyMessage.Parameters[1]

		return n.dht.nodePool.GetNode(predecessorId, predecessorAddr), nil
	} else {
		return nil, nil
	}
}

func (n *Node) ReplyPredecessor(node *Node) error {
	var reply Message

	if node == nil {
		reply = n.dht.self.NewMessage(_M_REPLY_PREDECESSOR, []string{})
	} else {
		reply = n.dht.self.NewMessage(_M_REPLY_PREDECESSOR, []string{
			node.Id().String(),
			node.GlobalAddress(),
		})
	}

	_, err := n.sendMessage(reply)

	if err != nil {
		return err
	}

	return nil
}

func (n *Node) Notify(node *Node) {
	requestMessage := n.dht.self.NewMessage(_M_NOTIFY, []string{})
	n.sendMessage(requestMessage)
}

func (n *Node) SendPing() error {
	queryMessage := n.dht.self.NewMessage(_M_REQUEST_PING, nil)
	_, err := n.sendMessage(queryMessage)

	if err != nil {
		return err
	}

  select {
  case _ = <-n.localInbound:
    return nil
  case <-time.After(time.Millisecond * 100):
    return errors.New("LOL")
  }
}

func (n *Node) ReplyPing() error {
	replyMessage := n.dht.self.NewMessage(_M_REPLY_PING, nil)
	_, err := n.sendMessage(replyMessage)

	if err != nil {
		return err
	}

	return nil
}

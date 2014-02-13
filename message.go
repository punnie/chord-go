package main

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

const (
	_M_REQUEST_SUCCESSOR = byte(iota)
	_M_REQUEST_PREDECESSOR
	_M_REQUEST_PING
	_M_NOTIFY
	_M_REPLY_SUCCESSOR
	_M_REPLY_PREDECESSOR
	_M_REPLY_PING
)

type Message struct {
	Intent     byte      `json:"intent"`
	Parameters []string  `json:"parameters"`
	Timestamp  time.Time `json:"timestamp"`
	SenderAddr string    `json:"senderaddr"`
	SenderHash string    `json:"senderhash"`
}

type Envelope struct {
	message *Message
	sender  *Node
}

func (m *Message) NewEvelope() *Envelope {
	return &Envelope{
		message: m,
	}
}

func (n *Node) NewMessage(intent byte, params []string) Message {
	return Message{
		Intent:     intent,
		Parameters: params,
		Timestamp:  time.Now().UTC(),
		SenderHash: n.Id().String(),
		SenderAddr: n.GlobalAddress(),
	}
}

func MessageDecode(buf []byte) (*Message, error) {
	decBuf, err := base64.URLEncoding.DecodeString(string(buf))

	if err != nil {
		return nil, err
	}

	m := new(Message)
	err = json.Unmarshal(decBuf, m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Message) MessageEncode() ([]byte, error) {
	mJSON, err := json.Marshal(m)

	if err != nil {
		return nil, err
	}

	encBuf := base64.StdEncoding.EncodeToString([]byte(mJSON))
	return []byte(encBuf), nil
}

func (m *Message) String() string {
	res, err := json.Marshal(m)

	if err != nil {
		panic(err)
	}

	return string(res)
}

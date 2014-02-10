package main

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

const (
	REQUEST_SUCCESSOR = "REQ SUCCESSOR"
	REPLY_SUCCESSOR   = "ACK SUCCESSOR"
	REQUEST_PING      = "PING"
	REPLY_PING        = "PONG"
)

type Message struct {
	Intent     string    `bson:"intent"`
	Parameters []string  `bson:"parameters"`
	Timestamp  time.Time `bson:"timestamp"`
  Sender     *Node
}

func NewFindSuccessorMessage(params []string) *Message {
	return &Message{
		Intent:     REQUEST_SUCCESSOR,
		Parameters: params,
		Timestamp:  time.Now().UTC(),
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

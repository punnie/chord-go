package main

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type Message struct {
	Intent     string    `bson:"intent"`
	Parameters []string  `bson:"parameters"`
	Timestamp  time.Time `bson:"timestamp"`
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

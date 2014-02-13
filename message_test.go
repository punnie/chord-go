package main

import (
	"encoding/base64"
	"testing"
	"time"
)

func TestMarshalling(t *testing.T) {
	ts := time.Now().UTC()
	m1 := &Message{Intent: _M_REPLY_PING, Timestamp: ts, Parameters: []string{"65", "áéíõú"}}

	buf, err := m1.MessageEncode()

	if err != nil {
		t.Fatal(err.Error())
	}

	m2, err := MessageDecode(buf)

	if err != nil {
		decBuf, err := base64.URLEncoding.DecodeString(string(buf))
		t.Log(string(decBuf))
		t.Fatal(err.Error())
	}

	if m1.Intent != m2.Intent {
		t.Fatal("Intents don't match!")
	}

	if m1.Timestamp.String() != m2.Timestamp.String() {
		t.Log(m1.Timestamp.String())
		t.Log(m2.Timestamp.String())
		t.Fatal("Timestamps don't match!")
	}

	for i, param := range m1.Parameters {
		if param != m2.Parameters[i] {
			t.Fatal("Parameters don't match!")
		}
	}
}

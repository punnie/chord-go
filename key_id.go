package main

import (
	"crypto/sha1"
	"encoding/hex"
  "bytes"
)

const (
  sha1Min = "0000000000000000000000000000000000000000"
  sha1Max = "ffffffffffffffffffffffffffffffffffffffff"
)

type KeyID struct {
	hex     string
}

func NewKeyID(hex string) *KeyID {
  return &KeyID{
    hex: hex,
  }
}

func (k *KeyID) GenerateNodeKeyID(seed string) error {
	s := sha1.New()
	s.Write([]byte(seed))

	h := hex.EncodeToString(s.Sum(nil))
	k.hex = h

	return nil
}

func (k *KeyID) String() string {
	return k.hex
}

//func (k *KeyID) elementOfLeftOpen(left *KeyID, right *KeyID) bool {
//}

//func (k *KeyID) elementOfBothOpen(left *KeyID, right *KeyID) bool {
//}

func (k *KeyID) elementOf(left *KeyID, right *KeyID) bool {
  keyHex := []byte(k.String())
  leftHex := []byte(left.String())
  rightHex := []byte(right.String())

  if(bytes.Compare(leftHex, rightHex) < 0) {
    return bytes.Compare(keyHex, leftHex) > 0 && bytes.Compare(keyHex, rightHex) <= 0
  } else {
    overFlow := bytes.Compare(keyHex, leftHex) > 0 && bytes.Compare(keyHex, []byte(sha1Max)) <= 0
    underFlow := bytes.Compare(keyHex, []byte(sha1Min)) > 0 && bytes.Compare(keyHex, rightHex) <= 0

    return overFlow || underFlow
  }
}

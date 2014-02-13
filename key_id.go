package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

const (
	_SHA256MIN = "0000000000000000000000000000000000000000000000000000000000000000"
	_SHA256MAX = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
)

const (
	// Taken from http://golang.org/src/pkg/math/big/arith.go
	// Compute the size _S of a Word in bytes.
	_m    = ^big.Word(0)
	_logS = _m>>8&1 + _m>>16&1 + _m>>32&1
	_S    = 1 << _logS
	_W    = _S << 3
)

type KeyID struct {
	hex     string
	integer *big.Int
}

func NewKeyID() *KeyID {
	return &KeyID{}
}

func (k *KeyID) SetHash(hex string) *KeyID {
	i, success := big.NewInt(0).SetString(hex, 16)

	if success != true {
		panic("lol errors")
	}

	return &KeyID{
		hex:     hex,
		integer: i,
	}
}

func (k *KeyID) GenerateKeyID(seed string) (*KeyID, error) {
	s := sha256.New()
	s.Write([]byte(seed))

	h := hex.EncodeToString(s.Sum(nil))
	k.hex = h

	i, success := big.NewInt(0).SetString(h, 16)

	if success != true {
		panic("lol errors")
	}

	k.integer = i

	return k, nil
}

func (k *KeyID) String() string {
	return k.hex
}

func (k *KeyID) elementOf(left *KeyID, right *KeyID) bool {
	keyHex := []byte(k.String())
	leftHex := []byte(left.String())
	rightHex := []byte(right.String())

	if bytes.Compare(leftHex, rightHex) < 0 {
		return bytes.Compare(keyHex, leftHex) > 0 && bytes.Compare(keyHex, rightHex) <= 0
	} else {
		overFlow := bytes.Compare(keyHex, leftHex) > 0 && bytes.Compare(keyHex, []byte(_SHA256MAX)) <= 0
		underFlow := bytes.Compare(keyHex, []byte(_SHA256MIN)) >= 0 && bytes.Compare(keyHex, rightHex) <= 0

		return overFlow || underFlow
	}
}

// Credit for this func: Hugo Peixoto (https://github.com/hugopeixoto)
func (key *KeyID) add2nModK(n int, k int) *KeyID {
	words := key.integer.Bits()
	extended := make([]big.Word, (k+_W-1)/_W)
	copy(extended, words)

	for n/_W < k {
		i := n / _W
		w := extended[i]

		new_w := w + (big.Word(1) << uint(n-i*_W))
		extended[i] = new_w

		if new_w < w {
			n = (i + 1) * _W
		} else {
			break
		}
	}

	result := big.NewInt(0).SetBits(extended)
	hex := fmt.Sprintf("%040x", result)

	return &KeyID{
		hex:     hex,
		integer: result,
	}
}

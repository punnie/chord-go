package main

import (
	"testing"
)

func TestKeyElementOf(t *testing.T) {
	k1 := NewKeyID().SetHash("6666666666666666666666666666666666666666")
	k2 := NewKeyID().SetHash("9999999999999999999999999999999999999999")
	k3 := NewKeyID().SetHash("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	test1 := k1.elementOf(k2, k3)

	if test1 {
		t.Fatal("Comparison failed")
	}

	test2 := k2.elementOf(k1, k3)

	if !test2 {
		t.Fatal("Comparison failed")
	}

	test3 := k1.elementOf(k3, k2)

	if !test3 {
		t.Fatal("Comparison failed")
	}

	test4 := k3.elementOf(k2, k1)

	if !test4 {
		t.Fatal("Comparison failed")
	}

	test5 := k1.elementOf(k1, k3)

	if test5 {
		t.Fatal("Comparison failed")
	}
}

func TestFingerCalculations(t *testing.T) {
	k1 := NewKeyID().SetHash("0fffffffffffffffffffffffffffffffffffffff")

	for i := 0; i < BITS; i++ {
		k1.add2nModK(i, BITS)
	}
}

package main

import (
	"bytes"
	"math/big"
	"testing"
)

// Verify hash conversion
func TestHash(t *testing.T) {
	address := "127.0.0.1:4232"
	expected := hash(address)

	haStr := string(expected)
	actual := []byte(haStr)
	if bytes.Compare(actual, expected) != 0 {
		t.Errorf("Want: %x; Got: %x\n", expected, actual)
	}

	actual = []byte("sdsdsjasdjkgasjhdgjhasvdjhasvdjhvsajhdvjhasvdjhsvdjhvshjdvhjsvdjhs")
	if bytes.Compare(actual, expected) == 0 {
		t.Errorf("Want: %x; Got: %x should be different.\n", expected, actual)
	}
}

func TestMsgDataSerialization(t *testing.T) {
	data := make(map[string]interface{})
	data["Key"] = "Name"
	data["Value"] = "chord"

	m := &msg{make(map[string]interface{})}
	m.content["data"] = data
	b, _ := m.Marshal()

	m = &msg{}
	d, _ := m.Unmarshal(b)
	dm := d["data"].(map[string]interface{})

	if dm["Value"] != "chord" {
		t.Errorf("Want %s; Got: %s", data["Value"], dm["Value"])
	}
}

// test if there is any difference between
// between bytes comparison and big int comparison
func TestIdMath(t *testing.T) {
	x := hash("132")
	a := hash("123")
	b := hash("chord")

	xBigInt := new(big.Int).SetBytes(x)
	aBigInt := new(big.Int).SetBytes(a)
	bBigInt := new(big.Int).SetBytes(b)

	actual := bytes.Compare(x, a) | bytes.Compare(x, b) | bytes.Compare(a, b)        // 1 | 1 | -1
	expected := (xBigInt.Cmp(aBigInt) | xBigInt.Cmp(bBigInt) | aBigInt.Cmp(bBigInt)) // -1

	// Test equals
	if actual != expected {
		t.Errorf("Got: %d; Want: %d", actual, expected)
	}

	// Test not equal
	actual &= 1
	if actual == expected {
		t.Errorf("Got: %d; Want: %d", actual, expected)
	}
}

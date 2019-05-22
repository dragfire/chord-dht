package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
)

func hash(id string) []byte {
	h := sha1.New()

	_, err := h.Write([]byte(id))
	checkErrPanic(err)
	return h.Sum(nil)
}

func checkErrPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func checkErrReturn(err error) error {
	if err != nil {
		return err
	}

	return nil
}

func checkErrPrint(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// a < x < b
func between(x, a, b []byte) bool {
	return bytes.Compare(x, a) == 1 && bytes.Compare(x, b) == -1
}

package main

import (
	"fmt"
	"math/big"
)

type fingerEntry struct {
	startId []byte
	node    *remoteNode
}

func createFingerTable(node *Node) []*fingerEntry {
	m := node.options.keySize
	fingerTable := make([]*fingerEntry, m)
	// fmt.Println("createFingerTable:", node.remoteNode)

	for i := range fingerTable {
		// n + 2^i % 2^m
		idBigInt := new(big.Int).SetBytes(node.Id)
		tmp := new(big.Int).Add(idBigInt, new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), nil))
		newNodeID := new(big.Int).Rem(tmp, new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(m)), nil))
		fingerTable[i] = &fingerEntry{newNodeID.Bytes(), node.remoteNode}
	}

	return fingerTable
}

func (fe *fingerEntry) String() string {
	return fmt.Sprintf("FingerEntry startId: %x remoteNode: %s", fe.startId, fe.node)
}

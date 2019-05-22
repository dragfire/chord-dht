package main

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"time"
)

// Options stores all node configuations
type Options struct {
	keySize          int
	intervalDuration time.Duration
}

func defaultOptions() *Options {
	return &Options{20, 10}
}

type Node struct {
	// since each node doesn't have access to the actual object of other node's
	// in the ring, we need remoteNode{id, address} to communicate with each other
	*remoteNode

	store       *kvStore
	predecessor *remoteNode
	successor   *remoteNode
	fingerTable []*fingerEntry
	options     *Options
	next        int
	active      bool
}

// NewNode creates new node
func NewNode(address string) *Node {
	nId := hash(address)
	fmt.Println("NodeID:", address, nId)
	rNode := &remoteNode{nId, address}
	node := &Node{remoteNode: rNode, store: newKvStore(), options: defaultOptions()}
	server := newServer(node)
	node.fingerTable = createFingerTable(node)
	node.successor = node.remoteNode

	ctx, cancel := context.WithCancel(context.Background())

	go server.start(cancel)

	// periodically verify n's immediate successor,
	// and tell the successor about n.
	// ticker := time.NewTicker(time.Millisecond * 500)
	go stabilize(node, ctx)

	// periodically refreshes finger table entries
	go fixFingerTables(node, ctx)

	return node
}

// Join allows (n) Node to join (other) remote node
func (n *Node) Join(other *remoteNode) error {
	node, err := n.FindSuccessor(other.Id)
	checkErrPanic(err)

	fmt.Printf("Node with Address: %s joining Node: %s\n", n.Address, other.Address)
	err = n.buildFingers(other)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// checkErrPanic(err)
	n.active = true
	n.successor = node
	return nil
}

// ask other(remoteNode) to build n's finger table
func (n *Node) buildFingers(other *remoteNode) error {
	id := new(big.Int).SetBytes(n.Id)

	j := 0
	m := n.options.keySize - 1

	// i => 0-indexed
	for i := m; i >= j; i-- {
		// n + 2^(i-1)

		tmp := new(big.Int).Add(id, new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), nil))

		node, err := other.FindSuccessorRemote(tmp.Bytes())
		checkErrPanic(err)

		n.fingerTable[i].node = node
	}

	return nil
}

func (n *Node) printFingerTable() {
	fmt.Printf("\n\n===========Finger Table for node: %v ============\n", n.remoteNode)
	for i := 0; i < len(n.fingerTable); i++ {
		fmt.Println(n.fingerTable[i])
	}
	fmt.Print("\n\n")
}

// FindSuccessor finds a successor
func (n *Node) FindSuccessor(id []byte) (*remoteNode, error) {
	fmt.Printf("Finding successor for node with ID: %x\n", id)

	if bytes.Compare(id, n.Id) == 1 && (bytes.Compare(id, n.successor.Id) == 0 || bytes.Compare(id, n.successor.Id) == -1) {
		return n.successor, nil
	}

	cpNode := n.closestPrecedingNode(id)

	if bytes.Compare(cpNode.Id, n.Id) == 0 {
		fmt.Println("ClosestPrecedingNode:: is node itself", cpNode)
		return cpNode, nil
	}

	fmt.Printf("ClosestPrecedingNode:: %s", cpNode)
	resp, err := cpNode.FindSuccessorRemote(id)
	// fmt.Println(resp, err)
	return resp, err
}

// find the closest preceding node given an Id
func (n *Node) closestPrecedingNode(rid []byte) *remoteNode {
	m := n.options.keySize
	// fmt.Println(m, len(n.fingerTable))
	for i := m - 1; i >= 0; i-- {
		fe := n.fingerTable[i]
		id := fe.startId

		// check if finger entry id is in between n.Id and rid
		if bytes.Compare(id, n.Id) == 1 && bytes.Compare(id, rid) == -1 {
			return fe.node
		}
	}

	return n.remoteNode
}

// this is called periodically when node is created
func stabilize(n *Node, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down stabilize goroutine for node:", n.remoteNode)
			return
		default:
			if n.active == false {
				continue
			}
			fmt.Println("Stabilize:", n.remoteNode)
			succ := n.successor

			pred, err := succ.GetPredecessorRemote()
			checkErrPanic(err)
			// n < pred < successor
			if pred != nil && between(pred.Id, n.Id, n.successor.Id) {
				n.successor = pred
			}

			succ.NotifyRemote(n.Address)

			time.Sleep(n.options.intervalDuration * time.Second)
		}
	}

}

func (n *Node) notify(other *remoteNode) {
	pred := n.predecessor
	if pred == nil || (&remoteNode{} == pred) || between(other.Id, pred.Id, n.Id) {
		n.predecessor = other
	}
}

func fixFingerTables(n *Node, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down fixfingerTables goroutine for node:", n.remoteNode)
			return
		default:
			if n.active == false {
				continue
			}
			n.next++
			m := n.options.keySize
			if n.next == m {
				n.next = 0
			}
			// TODO: n + 2^next or n + 2^next % 2^m ???
			nextId := new(big.Int).Add(new(big.Int).SetBytes(n.Id), new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(n.next)), nil))
			node, err := n.FindSuccessor(nextId.Bytes())
			checkErrPanic(err)
			n.fingerTable[n.next].node = node

			n.printFingerTable()

			time.Sleep(n.options.intervalDuration * time.Second)
		}
	}

}

// TODO: add check-predecessor
// this needs to be called periodically
func (n *Node) checkPredecessor() {
	panic("Implement me")
}

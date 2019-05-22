package main

import (
	"fmt"
	"time"
)

func spawnNode(address string) {
	node := NewNode(address)

	fmt.Printf("NodeID: %s spawned!\n", node.Address)
}

func join(recipientNode, sponsoringNode string) {

	// node2 wants to join node1
	// node1 is the sponsoring-node
	// node2 is the recipient-node

	// contact node2
	client := newClient(&remoteNode{Address: recipientNode})

	msg, err := joinRingMsg(sponsoringNode).Marshal()
	checkErrPanic(err)
	// send join-ring msg, sponsoring-node: node1
	_, err = client.requester.Send(string(msg), 0)
	checkErrPanic(err)
	response, err := client.requester.Recv(0)
	fmt.Println("Response:", response)
}

func leave(address, mode string) {
	client := newClient(&remoteNode{Address: address})
	m := leaveRingMsg(mode)
	client.sendReq(m)
}

func listItems(address string) {
	client := newClient(&remoteNode{Address: address})
	m := listItemsMsg(address)
	client.sendReq(m)

	resp, err := client.requester.Recv(0)
	checkErrPanic(err)
	fmt.Printf("List Items for address: %s : %s", address, resp)
}

func put(address, key, val string) {
	client := newClient(&remoteNode{Address: address})
	m := newMsg()
	client.sendReq(m)

	resp, err := client.requester.Recv(0)
	checkErrPanic(err)
	fmt.Printf("Put for address: %s : %s", address, resp)
}

func get(address, key string) {
	client := newClient(&remoteNode{Address: address})
	m := newMsg()
	client.sendReq(m)

	resp, err := client.requester.Recv(0)
	checkErrPanic(err)
	fmt.Printf("Get for address and key: %s, %s :: %s", address, key, resp)
}

func remove(address, key string) {
	client := newClient(&remoteNode{Address: address})
	m := newMsg()
	client.sendReq(m)

	resp, err := client.requester.Recv(0)
	checkErrPanic(err)
	fmt.Printf("Remove for address and key: %s, %s :: %s", address, key, resp)
}

func main() {
	// Run chord ring
	// TODO: Implement test cases
	node1Addr := "127.0.0.1:5001"
	node2Addr := "127.0.0.1:5002"
	node3Addr := "127.0.0.1:5003"
	node4Addr := "127.0.0.1:5004"
	node5Addr := "127.0.0.1:5005"
	node6Addr := "127.0.0.1:5006"
	node7Addr := "127.0.0.1:5007"
	node8Addr := "127.0.0.1:5008"
	node9Addr := "127.0.0.1:5009"
	node10Addr := "127.0.0.1:5010"
	node11Addr := "127.0.0.1:5011"

	// start node goroutines
	go spawnNode(node1Addr)
	go spawnNode(node2Addr)
	go spawnNode(node3Addr)
	go spawnNode(node4Addr)
	go spawnNode(node5Addr)
	go spawnNode(node6Addr)
	go spawnNode(node7Addr)
	go spawnNode(node8Addr)
	go spawnNode(node9Addr)
	go spawnNode(node10Addr)
	go spawnNode(node11Addr)

	// send join msg to sponsoring node
	join(node1Addr, node2Addr)
	join(node3Addr, node1Addr)
	join(node4Addr, node3Addr)
	join(node5Addr, node3Addr)
	join(node6Addr, node3Addr)
	join(node7Addr, node3Addr)
	join(node8Addr, node5Addr)
	join(node9Addr, node3Addr)
	join(node10Addr, node8Addr)
	join(node11Addr, node9Addr)

	// Issue some random msg
	time.AfterFunc(2*time.Second, func() {
		put(node2Addr, "name", "abc")
		put(node1Addr, "city", "def")
		put(node3Addr, "tel", "ghi")

		time.AfterFunc(2*time.Second, func() {
			listItems(node1Addr)
			get(node3Addr, "tel")
			get(node1Addr, "city")
		})
		leave(node2Addr, OrderlyMode)
		leave(node2Addr, ImmediateMode)
		leave(node3Addr, ImmediateMode)
		leave(node4Addr, ImmediateMode)
	})
	// TODO: Shutdown properly using Notify signal
	c := make(chan int)
	<-c
}

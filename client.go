package main

import (
	"fmt"
	"time"

	zmq "github.com/pebbe/zmq4"
)

type client struct {
	requester *zmq.Socket
	*remoteNode
}

func newClient(rNode *remoteNode) *client {
	fmt.Println("Client with remoteNode:", rNode)
	address := "tcp://" + rNode.Address
	requester, err := zmq.NewSocket(zmq.REQ)
	checkErrPanic(err)
	// Set a 500ms client timeout
	requester.SetLinger(500 * time.Millisecond)

	err = requester.Connect(address)
	checkErrPanic(err)

	return &client{requester, rNode}
}

func unmarshalResponse(response string) (*remoteNode, error) {
	if response == "null" {
		return nil, nil
	}
	fmt.Println("unmarshalResponse:", response, response == "null")
	rn := &remoteNode{}
	rn.Unmarshal(response)
	return rn, nil
}

func (c *client) sendReq(m *msg) error {
	data, err := m.Marshal()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Send request", m)

	// NOTE: Do not change this to SendMessage(data []byte),
	// zmq has a bug or some functionality that I am not aware of
	// if you use SendMessage, it sometimes sends an additional request
	// with request content "0", which is undesired behaviour.
	_, err = c.requester.Send(string(data), 0)

	return err
}

func (c *client) GetPredecessor(replyTo string) (*remoteNode, error) {
	msg := findRingPredecessorMsg(replyTo)
	// fmt.Println("GetPredecessor request:", msg)
	err := c.sendReq(msg)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	response, err := c.requester.Recv(0)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// msg := newMsg()

	fmt.Println("GetPredecessor response:", response)

	return unmarshalResponse(response)
}

// TODO: discuss, we don't actually need reply-to
func (c *client) FindSuccessor(replyTo string, id []byte) (*remoteNode, error) {
	err := c.sendReq(findRingSuccessorMsg(replyTo, string(id)))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	response, err := c.requester.Recv(0)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// msg := newMsg()

	fmt.Println("FindSuccessor response:", response)

	return unmarshalResponse(response)
}

func (c *client) Notify(address string) {
	msg := ringNotifyMsg(address)
	fmt.Println("Notify request:", msg)

	err := c.sendReq(msg)
	checkErrPanic(err)

	response, err := c.requester.Recv(0)
	checkErrPanic(err)
	// msg := newMsg()

	fmt.Println("Notify response:", response)
}

func (c *client) DeleteKeyRPC(replyTo string, key string) error {
	err := c.sendReq(newMsg())
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = c.requester.Recv(0)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("RPC Delete Key Request Successful")

	return nil
}

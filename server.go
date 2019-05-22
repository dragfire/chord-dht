package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	zmq "github.com/pebbe/zmq4"
)

type chordServer struct {
	handlers  map[string]handler
	responder *zmq.Socket
	node      *Node
}

func newServer(node *Node) *chordServer {
	address := node.Address
	address = "tcp://" + address
	// Socket to talk to clients
	responder, err := zmq.NewSocket(zmq.REP)
	checkErrPanic(err)

	// defer responder.Close()
	err = responder.Bind(address)
	checkErrPanic(err)

	fmt.Println("Server started:", address)
	return &chordServer{make(map[string]handler), responder, node}
}

func (c *chordServer) start(cancel context.CancelFunc) {
	c.registerHandlers()
	// start msg receiver worker
	for {
		m, err := c.responder.Recv(0)
		checkErrPanic(err)
		fmt.Println("Request received:", m, "Address:", c.node.Address)

		// TODO: debug
		if len(m) == 1 {
			fmt.Println("WARN: msg body shouldn't be empty")
			c.responder.Send("WARN: unexpected state", 0)
			continue
		}

		msg := newMsg()
		content, err := msg.Unmarshal([]byte(m))
		msg.content = content
		checkErrPanic(err)

		if c.checkLeaveRing(cancel, msg) {
			return
		}

		c.processMsg(msg)
	}
}

func (c *chordServer) checkLeaveRing(cancel context.CancelFunc, m *msg) bool {
	if m.Get(DoKey) == LeaveRing {
		cancel()
		c.responder.Close()

		if m.Get(ModeKey) == OrderlyMode {
			pred := c.node.predecessor
			succ := c.node.successor

			fmt.Println("Leave Orderly, pred:", pred, " succ:", succ, " node:", c.node.remoteNode)

			if pred != nil {
				pred.NotifyRemote(c.node.Address)
			}
			succ.NotifyRemote(c.node.Address)

			// transfer keys to successor
			store := c.node.store
			from := c.node.remoteNode
			to := succ
			for k, v := range store.content {
				hashedKey := hash(k)

				if between(hashedKey, from.Id, to.Id) || bytes.Compare(hashedKey, to.Id) == 1 {
					to.PutRemote(k, v)
				} else if pred != nil {
					pred.PutRemote(k, v)
				}
			}
			c.node.store.DeleteAll()
		}

		return true
	}

	return false
}

type handler interface {
	process(m *msg) error
}

type handlerFunc func(m *msg) error

func (f handlerFunc) process(m *msg) error {
	return f(m)
}

func (c *chordServer) processMsg(m *msg) {
	key := m.content[DoKey].(string)
	// fmt.Println("Handler Key:", key)
	fn := c.handlers[key]
	fn.process(m)
}

func (c *chordServer) registerHandlers() {
	// join-ring
	c.addHandler(JoinRing, handlerFunc(c.joinRingHandler))

	// find-ring-successor
	c.addHandler(FindRingSuccessor, handlerFunc(c.findRingSuccessorHandler))

	// find-ring-predecessor
	c.addHandler(FindRingPredecessor, handlerFunc(c.findRingPredecessorHandler))

	// list-fingers
	c.addHandler(GetRingFingers, handlerFunc(c.getRingFingersHandler))

	// ring-notify
	c.addHandler(RingNotify, handlerFunc(c.ringNotifyHandler))

	// list-items
	c.addHandler(ListItems, handlerFunc(c.listItemsHandler))
}

func (c *chordServer) addHandler(key string, h handler) {
	c.handlers[key] = h
}

func (c *chordServer) joinRingHandler(m *msg) error {
	fmt.Println(JoinRing, m)
	address := m.Get(SponsoringNodeKey).(string)
	rNode := &remoteNode{Id: hash(address), Address: address}
	c.node.Join(rNode)
	_, err := c.responder.Send("Ok", 0)
	return err
}

func (c *chordServer) findRingSuccessorHandler(m *msg) error {
	fmt.Println(FindRingSuccessor, m)
	node, err := c.node.FindSuccessor([]byte(m.Get("id").(string)))
	checkErrPanic(err)
	fmt.Println("FindRingSuccessorHandler", node)
	_, err = c.responder.Send(node.Marshal(), 0)
	return err
}

func (c *chordServer) findRingPredecessorHandler(m *msg) error {
	fmt.Println(FindRingPredecessor, m)
	pred := c.node.predecessor
	_, err := c.responder.Send(pred.Marshal(), 0)
	return err
}

func (c *chordServer) getRingFingersHandler(m *msg) error {
	ft := c.node.fingerTable
	data, err := json.Marshal(ft)
	checkErrPanic(err)
	_, err = c.responder.Send(string(data), 0)
	return err
}

func (c *chordServer) ringNotifyHandler(m *msg) error {
	fmt.Println(RingNotify, m)
	address := m.Get(ReplyToKey).(string)
	c.node.notify(&remoteNode{Id: hash(address), Address: address})
	_, err := c.responder.Send("Ok", 0)
	return err
}

func (c *chordServer) listItemsHandler(m *msg) error {
	fmt.Println(ListItems, m)

	data := c.node.store.Marshal()
	_, err := c.responder.Send(data, 0)

	return err
}

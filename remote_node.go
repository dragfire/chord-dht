package main

import (
	"encoding/json"
	"fmt"
)

// remoteNode contains the actual node Id(hash value)
// and remote address(ip:port)
type remoteNode struct {
	Id      []byte
	Address string // ip:port
}

func (r *remoteNode) String() string {
	return fmt.Sprintf("Id: %x, Address: %s", r.Id, r.Address)
}

// TODO: cache client connections, no need to create every time

// GetPredecessor
func (rn *remoteNode) GetPredecessorRemote() (*remoteNode, error) {
	fmt.Println("GetPredecessorRemote", rn)
	client := newClient(rn)
	return client.GetPredecessor(rn.Address)
}

// FindSuccessorRemote contacts the appropriate remoteNode
// to find the successor
func (rn *remoteNode) FindSuccessorRemote(id []byte) (*remoteNode, error) {
	client := newClient(rn)
	return client.FindSuccessor(rn.Address, id)
}

func (rn *remoteNode) NotifyRemote(address string) {
	client := newClient(rn)
	client.Notify(address)
}

func (rn *remoteNode) Marshal() string {
	data, err := json.Marshal(&rn)
	checkErrPanic(err)
	return string(data)
}

func (rn *remoteNode) Unmarshal(data string) error {
	err := json.Unmarshal([]byte(data), &rn)
	checkErrPanic(err)
	return err
}

func (rn *remoteNode) GetRemote(key string) (string, error) {
	// client := newClient(rn)
	// TODO: Implementz
	// return client.Get(rn.Address, key)
	return "", nil
}

func (rn *remoteNode) GetValRemote(key string) (string, error) {
	// client := newClient(rn)
	// TODO: Implement
	// return client.GetVal(rn.Address, key)
	return "", nil
}

func (rn *remoteNode) PutRemote(key string, value string) {
	// client := newClient(rn)
	// TODO: Implement
	// client.Put(rn.Address, key, value)
}

func (rn *remoteNode) PutKeyValRemote(key string, value string) {
	// client := newClient(rn)
	// TODO: Implement
	// client.PutKeyVal(rn.Address, key, value)
}

func (rn *remoteNode) DeleteRemote(key string) {
	// client := newClient(rn)
	// TODO: Implement
	// client.Delete(rn.Address, key)
}

func (rn *remoteNode) DeleteKeyRemote(key string) {
	// TODO: Implement
	// client := newClient(rn)
	// client.DeleteKey(rn.Address, key)
}

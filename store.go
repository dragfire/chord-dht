package main

import (
	"bytes"
	"encoding/json"
	"errors"
)

type kvStore struct {
	content map[string]string
}

func newKvStore() *kvStore {
	return &kvStore{make(map[string]string)}
}

func (s *kvStore) Put(key, val string) error {
	hashedKey := hash(key)
	s.content[string(hashedKey)] = val
	return nil
}

func (s *kvStore) Get(key string) (string, error) {
	hashedKey := hash(key)
	v, ok := s.content[string(hashedKey)]
	if !ok {
		return "", errors.New("Key not found")
	}
	return v, nil
}

func (s *kvStore) Delete(key string) error {
	hashedKey := hash(key)
	_, ok := s.content[string(hashedKey)]
	if !ok {
		return errors.New("Key not found")
	}
	delete(s.content, string(hashedKey))
	return nil
}

func (s *kvStore) TransferKeys(from, to *remoteNode) {
	for k, v := range s.content {
		hashedKey := hash(k)

		if between(hashedKey, from.Id, to.Id) || bytes.Compare(hashedKey, to.Id) == 1 {
			to.PutRemote(k, v)
		}
	}
}

func (s *kvStore) DeleteAll() {
	for k := range s.content {
		delete(s.content, k)
	}
}

func (s *kvStore) Marshal() string {
	data, err := json.Marshal(s.content)
	checkErrPanic(err)
	return string(data)
}

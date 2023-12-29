package app

import (
	"crypto/sha256"
	"encoding/hex"
)

type Storage interface {
	Save(url string) (string, error)
	Load(key string) (string, error)
}

type InMemoryStorage struct {
	hashTable map[string]string
}

func NewInMemoryStorage() Storage {
	return &InMemoryStorage{hashTable: make(map[string]string)}
}

func (s *InMemoryStorage) Save(url string) (string, error) {
	key := generateKey(url)
	s.hashTable[key] = url
	return key, nil
}

func (s *InMemoryStorage) Load(key string) (string, error) {
	url, exists := s.hashTable[key]
	if !exists {
		return "", ErrNotFound
	}
	return url, nil
}

func generateKey(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	return hex.EncodeToString(hasher.Sum(nil)[:8])
}

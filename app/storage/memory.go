package storage

import "go-fake-smtp/app/message"

// In-memory storage.
//
// Messages are stored in memory and are lost upon application restart.
type MemoryBackend struct {
	messages []*message.Message
}

func (storage *MemoryBackend) Add(msg *message.Message) {
	storage.messages = append(storage.messages, msg)
}

func (storage *MemoryBackend) Count() int {
	return len(storage.messages)
}

func (storage *MemoryBackend) GetAll() []*message.Message {
	return storage.messages
}

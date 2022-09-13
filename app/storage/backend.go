package storage

import "go-fake-smtp/app/message"

// Interface for message storage backend.
//
// Libraries may extends this to implement their own storage.
type StorageBackend interface {
	// Add new message to storage.
	Add(msg *message.Message)

	// Returns the number of stored messages.
	Count() int

	// Returns a slice of all stored messages.
	GetAll() []*message.Message
}

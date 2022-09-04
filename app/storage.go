package app

// Interface for arbitrary message storage.
//
// Libraries may extends this to implement their own storage.
type MessageStorage interface {
	// Add new message to storage.
	Add(msg Message)

	// Returns the number of stored messages.
	Count() int

	// Returns a slice of all stored messages.
	GetAll() []Message
}

// In-memory storage.
//
// Stored messages are lost upon application restart.
type MemoryStorage struct {
	messages []Message
}

func (storage *MemoryStorage) Add(msg Message) {
	storage.messages = append(storage.messages, msg)
}

func (storage *MemoryStorage) Count() int {
	return len(storage.messages)
}

func (storage *MemoryStorage) GetAll() []Message {
	return storage.messages
}

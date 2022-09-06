package app

// Interface for arbitrary message storage backend.
//
// Libraries may extends this to implement their own storage.
type StorageBackend interface {
	// Add new message to storage.
	Add(msg *Message)

	// Returns the number of stored messages.
	Count() int

	// Returns a slice of all stored messages.
	GetAll() []*Message
}

// Central storage for everything mail.
type Storage struct {
	// Storage backend that stores messages.
	Backend StorageBackend

	// Index of all available mailboxes.
	MailboxIndex MailboxIndex
}

// Creates new central storage.
func NewStorage(backend StorageBackend) *Storage {
	return &Storage{
		Backend: backend,
		MailboxIndex: MailboxIndex{
			mailboxes: make(map[string][]*Message),
		},
	}
}

// In-memory storage.
//
// Stored messages are lost upon application restart.
type MemoryStorageBackend struct {
	messages []*Message
}

func (storage *MemoryStorageBackend) Add(msg *Message) {
	storage.messages = append(storage.messages, msg)
}

func (storage *MemoryStorageBackend) Count() int {
	return len(storage.messages)
}

func (storage *MemoryStorageBackend) GetAll() []*Message {
	return storage.messages
}

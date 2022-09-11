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

	// Index of all known mailboxes.
	//
	// The key is unique mailbox ID.
	mailboxes map[string]*Mailbox

	// Index of all known mailboxes.
	//
	// The key is mailbox name.
	mailboxesByName map[string]*Mailbox

	// Index of all messages bound to specific mailbox.
	//
	// The key is unique mailbox ID. The key of nested map is unique message ID.
	mailboxMessages map[string]map[string]*Message

	// Index of all known messages.
	//
	// The key is unique message ID.
	messages map[string]*Message

	// Index of mailboxes to which messages are bound.
	//
	// The key is unique message ID.
	messageMailboxes map[string]*Mailbox
}

// Adds new message to the storage.
//
// Requires name of the mailbox (not unique mailbox ID) to append message to.
func (storage *Storage) Add(message *Message, mailboxName string) {
	mailbox := storage.addMailbox(mailboxName)

	if _, ok := storage.mailboxMessages[mailbox.Id]; !ok {
		storage.mailboxMessages[mailbox.Id] = make(map[string]*Message)
	}

	storage.messages[message.Id] = message
	storage.messageMailboxes[message.Id] = mailbox
	storage.mailboxMessages[mailbox.Id][message.Id] = message

	storage.Backend.Add(message)
}

// Returns the number of known mailboxes.
func (storage *Storage) CountMailboxes() int {
	return len(storage.mailboxes)
}

// Returns the number of known messages bound to specific mailbox.
func (storage *Storage) CountMessages(mailboxId string) int {
	if messages, ok := storage.mailboxMessages[mailboxId]; ok {
		return len(messages)
	}

	return 0
}

// Returns all known mailboxes.
func (storage *Storage) GetMailboxes() []*Mailbox {
	mailboxes := make([]*Mailbox, 0, len(storage.mailboxes))

	for _, mailbox := range storage.mailboxes {
		mailboxes = append(mailboxes, mailbox)
	}

	return mailboxes
}

// Returns all known messages bound to specific mailbox.
func (storage *Storage) GetMessages(mailboxId string) []*Message {
	if mailboxMessages, ok := storage.mailboxMessages[mailboxId]; ok {
		messages := make([]*Message, 0, len(mailboxMessages))

		for _, message := range mailboxMessages {
			messages = append(messages, message)
		}

		return messages
	}

	return make([]*Message, 0)
}

// Registers mailbox name and returns corresponding mailbox.
//
// When mailbox name is already registered just returns mailbox.
func (storage *Storage) addMailbox(name string) *Mailbox {
	if mailbox, ok := storage.mailboxesByName[name]; ok {
		return mailbox
	}

	mailbox := NewMailbox(name)

	storage.mailboxes[mailbox.Id] = mailbox
	storage.mailboxesByName[mailbox.Name] = mailbox

	return mailbox
}

// Creates new central storage.
func NewStorage(backend StorageBackend) *Storage {
	return &Storage{
		Backend:          backend,
		mailboxes:        make(map[string]*Mailbox),
		mailboxesByName:  make(map[string]*Mailbox),
		mailboxMessages:  make(map[string]map[string]*Message),
		messages:         make(map[string]*Message),
		messageMailboxes: make(map[string]*Mailbox),
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

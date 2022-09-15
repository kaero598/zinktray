package storage

import (
	"go-fake-smtp/app/mailbox"
	"go-fake-smtp/app/message"
)

// Central storage for everything mail.
type Storage struct {
	// Index of all known mailboxes.
	//
	// The key is unique mailbox ID.
	mailboxes map[string]*mailbox.Mailbox

	// Index of all known mailboxes.
	//
	// The key is mailbox name.
	mailboxesByName map[string]*mailbox.Mailbox

	// Index of all messages bound to specific mailbox.
	//
	// The key is unique mailbox ID. The key of nested map is unique message ID.
	mailboxMessages map[string]map[string]*message.Message

	// Index of all known messages.
	//
	// The key is unique message ID.
	messages map[string]*message.Message

	// Index of mailboxes to which messages are bound.
	//
	// The key is unique message ID.
	messageMailboxes map[string]*mailbox.Mailbox
}

// Adds new message to the storage.
//
// Requires name of the mailbox (not unique mailbox ID) to append message to.
func (storage *Storage) Add(msg *message.Message, mailboxName string) {
	mailbox := storage.addMailbox(mailboxName)

	if _, ok := storage.mailboxMessages[mailbox.Id]; !ok {
		storage.mailboxMessages[mailbox.Id] = make(map[string]*message.Message)
	}

	storage.messages[msg.Id] = msg
	storage.messageMailboxes[msg.Id] = mailbox
	storage.mailboxMessages[mailbox.Id][msg.Id] = msg
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

// Deletes mailbox along with all its messages.
func (storage *Storage) DeleteMailbox(mailboxId string) {
	if mailbox, ok := storage.mailboxes[mailboxId]; ok {
		for _, message := range storage.mailboxMessages[mailboxId] {
			storage.pruneMessage(message, mailbox.Id)
		}

		delete(storage.mailboxMessages, mailbox.Id)

		storage.pruneMailbox(mailbox)
	}
}

// Deletes message.
func (storage *Storage) DeleteMessage(messageId string) {
	if message, ok := storage.messages[messageId]; !ok {
		mailbox := storage.messageMailboxes[messageId]

		storage.pruneMessage(message, mailbox.Id)
		storage.pruneMailbox(mailbox)
	}
}

// Returns mailbox.
//
// Returns nil for unknown mailbox.
func (storage *Storage) GetMailbox(mailboxId string) *mailbox.Mailbox {
	if mailbox, ok := storage.mailboxes[mailboxId]; ok {
		return mailbox
	}

	return nil
}

// Returns all known mailboxes.
func (storage *Storage) GetMailboxes() []*mailbox.Mailbox {
	mailboxes := make([]*mailbox.Mailbox, 0, len(storage.mailboxes))

	for _, mailbox := range storage.mailboxes {
		mailboxes = append(mailboxes, mailbox)
	}

	return mailboxes
}

// Returns message.
//
// Returns nil for unknown message.
func (storage *Storage) GetMessage(messageId string) *message.Message {
	if message, ok := storage.messages[messageId]; ok {
		return message
	}

	return nil
}

// Returns all known messages bound to specific mailbox.
func (storage *Storage) GetMessages(mailboxId string) []*message.Message {
	if mailboxMessages, ok := storage.mailboxMessages[mailboxId]; ok {
		messages := make([]*message.Message, 0, len(mailboxMessages))

		for _, message := range mailboxMessages {
			messages = append(messages, message)
		}

		return messages
	}

	return make([]*message.Message, 0)
}

// Registers mailbox name and returns corresponding mailbox.
//
// When mailbox name is already registered just returns mailbox.
func (storage *Storage) addMailbox(name string) *mailbox.Mailbox {
	if mailbox, ok := storage.mailboxesByName[name]; ok {
		return mailbox
	}

	mailbox := mailbox.NewMailbox(name)

	storage.mailboxes[mailbox.Id] = mailbox
	storage.mailboxesByName[mailbox.Name] = mailbox

	return mailbox
}

// Deletes all traces of message.
func (storage *Storage) pruneMessage(message *message.Message, mailboxId string) {
	delete(storage.messages, message.Id)
	delete(storage.messageMailboxes, message.Id)
	delete(storage.mailboxMessages[mailboxId], message.Id)
}

// Deletes all traces of mailbox.
func (storage *Storage) pruneMailbox(mailbox *mailbox.Mailbox) {
	if len(storage.mailboxMessages[mailbox.Id]) == 0 {
		delete(storage.mailboxes, mailbox.Id)
		delete(storage.mailboxesByName, mailbox.Name)
		delete(storage.mailboxMessages, mailbox.Id)
	}
}

// Creates new central storage.
func NewStorage() *Storage {
	return &Storage{
		mailboxes:        make(map[string]*mailbox.Mailbox),
		mailboxesByName:  make(map[string]*mailbox.Mailbox),
		mailboxMessages:  make(map[string]map[string]*message.Message),
		messages:         make(map[string]*message.Message),
		messageMailboxes: make(map[string]*mailbox.Mailbox),
	}
}

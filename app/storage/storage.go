package storage

import (
	"container/list"
	"errors"
	"sync"
	"zinktray/app/mailbox"
	"zinktray/app/message"
)

// ErrDuplicate error can be returned upon adding a message when another message with such ID is aready present
// in the storage in any mailbox.
var ErrDuplicate = errors.New("message with such ID already exists")

// ErrMailboxNotRegistered can be returned upon adding a message to a mailbox that is not yet registered.
var ErrMailboxNotRegistered = errors.New("mailbox is not registered")

// Storage represents central storage for everything mail.
type Storage struct {
	mailboxMutex sync.RWMutex
	messageMutex sync.RWMutex

	// Contains list of all registered mailboxes.
	mailboxList *list.List

	// Maps mailbox ID to its respective list element.
	mailboxElements map[string]*list.Element

	// Contains list of all registered messages.
	messageList *list.List

	// Maps message ID to its respective list element.
	messageElements map[string]*list.Element

	// Maps mailbox ID to a list of message IDs belonging to this mailbox.
	mailboxMessageIDs map[string]*list.List

	// Maps message ID to its respective list element inside mailbox it belongs to.
	mailboxMessageIDElements map[string]*list.Element

	// Maps message ID to ID of mailbox it belongs to.
	messageMailboxIDs map[string]string
}

// AddMailbox registers mailbox ID and returns corresponding mailbox.
//
// When mailbox ID is already registered returns that mailbox.
func (storage *Storage) AddMailbox(mailboxId string) *mailbox.Mailbox {
	storage.mailboxMutex.Lock()

	defer storage.mailboxMutex.Unlock()

	var mbx *mailbox.Mailbox

	if element, ok := storage.mailboxElements[mailboxId]; ok {
		if m, ok := element.Value.(*mailbox.Mailbox); ok {
			return m
		}
	}

	mbx = mailbox.NewMailbox(mailboxId)

	storage.mailboxElements[mbx.ID] = storage.mailboxList.PushBack(mbx)
	storage.mailboxMessageIDs[mailboxId] = list.New()

	return mbx
}

// AddMessage stores new message and binds it to mailbox with provided ID.
// Returns ErrDuplicate error upon adding message with an ID that is already present in the storage in any mailbox.
func (storage *Storage) AddMessage(msg *message.Message, mailboxID string) error {
	storage.messageMutex.Lock()
	storage.mailboxMutex.RLock()

	defer storage.mailboxMutex.RUnlock()
	defer storage.messageMutex.Unlock()

	element, ok := storage.mailboxElements[mailboxID]
	if !ok {
		return ErrMailboxNotRegistered
	}

	mbx, ok := element.Value.(*mailbox.Mailbox)
	if !ok {
		return ErrMailboxNotRegistered
	}

	if _, ok := storage.messageElements[msg.ID]; ok {
		return ErrDuplicate
	}

	storage.messageElements[msg.ID] = storage.messageList.PushFront(msg)
	storage.mailboxMessageIDElements[msg.ID] = storage.mailboxMessageIDs[mailboxID].PushFront(msg.ID)
	storage.messageMailboxIDs[msg.ID] = mbx.ID

	return nil
}

// CountMailboxes returns the number of registered mailboxes.
func (storage *Storage) CountMailboxes() int {
	return storage.mailboxList.Len()
}

// CountMessages returns the number of stored messages bound to specified mailbox.
func (storage *Storage) CountMessages(mailboxId string) int {
	storage.mailboxMutex.RLock()

	defer storage.mailboxMutex.RUnlock()

	if element, ok := storage.mailboxElements[mailboxId]; ok {
		if mbx, ok := element.Value.(*mailbox.Mailbox); ok {
			return storage.mailboxMessageIDs[mbx.ID].Len()
		}
	}

	return 0
}

// DeleteMailbox deletes registered mailbox along with all its messages.
func (storage *Storage) DeleteMailbox(mailboxID string) {
	storage.messageMutex.Lock()
	storage.mailboxMutex.Lock()

	defer storage.messageMutex.Unlock()
	defer storage.mailboxMutex.Unlock()

	if msgIDList, ok := storage.mailboxMessageIDs[mailboxID]; ok {
		next := msgIDList.Front()

		for next != nil {
			if messageID, ok := next.Value.(string); ok {
				delete(storage.mailboxMessageIDElements, messageID)
				delete(storage.messageMailboxIDs, messageID)

				if element, ok := storage.messageElements[messageID]; ok {
					storage.messageList.Remove(element)

					delete(storage.messageElements, messageID)
				}
			}

			next = next.Next()
		}

		delete(storage.mailboxMessageIDs, mailboxID)
	}

	if element, ok := storage.mailboxElements[mailboxID]; ok {
		storage.mailboxList.Remove(element)

		delete(storage.mailboxElements, mailboxID)
	}
}

// DeleteMessage deletes stored message.
func (storage *Storage) DeleteMessage(messageID string) {
	storage.messageMutex.Lock()
	storage.mailboxMutex.Lock()

	defer storage.messageMutex.Unlock()
	defer storage.mailboxMutex.Unlock()

	if mailboxID, ok := storage.messageMailboxIDs[messageID]; ok {
		if element, ok := storage.mailboxMessageIDElements[mailboxID]; ok {
			storage.mailboxMessageIDs[mailboxID].Remove(element)

			delete(storage.mailboxMessageIDElements, messageID)
		}

		delete(storage.messageMailboxIDs, messageID)
	}

	if element, ok := storage.messageElements[messageID]; ok {
		storage.messageList.Remove(element)

		delete(storage.messageElements, messageID)
	}
}

// GetMailbox returns registered mailbox.
//
// Returns nil for unregistered mailboxes.
func (storage *Storage) GetMailbox(mailboxId string) *mailbox.Mailbox {
	storage.mailboxMutex.RLock()

	defer storage.mailboxMutex.RUnlock()

	if element, ok := storage.mailboxElements[mailboxId]; ok {
		if m, ok := element.Value.(*mailbox.Mailbox); ok {
			return m
		}
	}

	return nil
}

// GetMailboxes returns a slice of all registered mailboxes.
func (storage *Storage) GetMailboxes() []*mailbox.Mailbox {
	storage.mailboxMutex.RLock()

	defer storage.mailboxMutex.RUnlock()

	var mailboxes []*mailbox.Mailbox

	if storage.mailboxList.Len() > 0 {
		mailboxes = make([]*mailbox.Mailbox, 0, storage.mailboxList.Len())
		next := storage.mailboxList.Front()

		for next != nil {
			if mbx, ok := next.Value.(*mailbox.Mailbox); ok {
				mailboxes = append(mailboxes, mbx)
			}

			next = next.Next()
		}
	}

	return mailboxes
}

// GetMessage returns stored message.
//
// Returns nil for unknown message.
func (storage *Storage) GetMessage(messageId string) *message.Message {
	storage.messageMutex.RLock()

	defer storage.messageMutex.RUnlock()

	if element, ok := storage.messageElements[messageId]; ok {
		if m, ok := element.Value.(*message.Message); ok {
			return m
		}
	}

	return nil
}

// GetMessages returns a list of all known messages bound to specified mailbox.
func (storage *Storage) GetMessages(mailboxId string) []*message.Message {
	storage.mailboxMutex.RLock()
	storage.messageMutex.RLock()

	defer storage.mailboxMutex.RUnlock()
	defer storage.messageMutex.RUnlock()

	var result []*message.Message

	if l, ok := storage.mailboxMessageIDs[mailboxId]; ok {
		result = make([]*message.Message, 0, l.Len())
		next := l.Front()

		for next != nil {
			if msgID, ok := next.Value.(string); ok {
				if element, ok := storage.messageElements[msgID]; ok {
					if msg, ok := element.Value.(*message.Message); ok {
						result = append(result, msg)
					}
				}
			}

			next = next.Next()
		}
	}

	return result
}

// NewStorage creates new central storage structure.
func NewStorage() *Storage {
	return &Storage{
		mailboxMutex: sync.RWMutex{},
		messageMutex: sync.RWMutex{},

		mailboxList:     list.New(),
		mailboxElements: make(map[string]*list.Element),

		messageList:     list.New(),
		messageElements: make(map[string]*list.Element),

		mailboxMessageIDs:        make(map[string]*list.List),
		mailboxMessageIDElements: make(map[string]*list.Element),

		messageMailboxIDs: make(map[string]string),
	}
}

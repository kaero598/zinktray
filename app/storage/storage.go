package storage

import (
	"container/list"
	"sync"
	"zinktray/app/mailbox"
	"zinktray/app/message"
)

// Storage represents central storage for everything mail.
type Storage struct {
	mailboxMutex sync.RWMutex
	messageMutex sync.RWMutex

	// Contains list of all registered mailboxes.
	mailboxList *list.List

	// Maps mailbox IDs to their respective list elements.
	mailboxElements map[string]*list.Element

	// Maps mailbox IDs to the list of messages belonging to them.
	messageList map[string]*list.List

	// Maps message IDs to their respective list elements.
	messageElements map[string]*list.Element

	// Maps message IDs to IDs of mailboxes they belong to.
	messageMailboxIds map[string]string
}

// Add stores new message and binds it to mailbox with provided ID.
func (storage *Storage) Add(msg *message.Message, mailboxID string) {
	storage.messageMutex.Lock()

	defer storage.messageMutex.Unlock()

	mbx := storage.registerMailbox(mailboxID)

	storage.messageMailboxIds[msg.ID] = mbx.ID
	storage.messageElements[msg.ID] = storage.messageList[mailboxID].PushFront(msg)
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
			return storage.messageList[mbx.ID].Len()
		}
	}

	return 0
}

// DeleteMailbox deletes registered mailbox along with all its messages.
func (storage *Storage) DeleteMailbox(mailboxId string) {
	storage.messageMutex.Lock()
	storage.mailboxMutex.Lock()

	defer storage.messageMutex.Unlock()
	defer storage.mailboxMutex.Unlock()

	if element, ok := storage.mailboxElements[mailboxId]; ok {
		if messages, ok := storage.messageList[mailboxId]; ok {
			next := messages.Front()

			for next != nil {
				if m, ok := next.Value.(*message.Message); ok {
					storage.purgeMessage(m)
				}

				next = next.Next()
			}
		}

		if mbx, ok := element.Value.(*mailbox.Mailbox); ok {
			storage.purgeMailbox(mbx)
		}
	}
}

// DeleteMessage deletes stored message.
func (storage *Storage) DeleteMessage(messageId string) {
	storage.messageMutex.Lock()
	storage.mailboxMutex.Lock()

	defer storage.messageMutex.Unlock()
	defer storage.mailboxMutex.Unlock()

	if element, ok := storage.messageElements[messageId]; ok {
		if msg, ok := element.Value.(*message.Message); ok {
			mbxId := storage.messageMailboxIds[msg.ID]

			if element, ok := storage.mailboxElements[mbxId]; ok {
				if mbx, ok := element.Value.(*mailbox.Mailbox); ok {
					storage.purgeMessage(msg)
					storage.purgeMailbox(mbx)
				}
			}
		}
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

// GetMailboxes returns a list of all registered mailboxes.
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

	if v, ok := storage.messageElements[messageId]; ok {
		if m, ok := v.Value.(*message.Message); ok {
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

	if messages, ok := storage.messageList[mailboxId]; ok {
		result = make([]*message.Message, 0, messages.Len())
		next := messages.Front()

		for next != nil {
			if m, ok := next.Value.(*message.Message); ok {
				result = append(result, m)
			}

			next = next.Next()
		}
	}

	return result
}

// registerMailbox registers mailbox ID and returns corresponding mailbox.
//
// When mailbox ID is already registered returns that mailbox.
func (storage *Storage) registerMailbox(mailboxId string) *mailbox.Mailbox {
	var mbx *mailbox.Mailbox

	if element, ok := storage.mailboxElements[mailboxId]; ok {
		if m, ok := element.Value.(*mailbox.Mailbox); ok {
			return m
		}
	}

	mbx = mailbox.NewMailbox(mailboxId)

	storage.mailboxElements[mbx.ID] = storage.mailboxList.PushBack(mbx)
	storage.messageList[mailboxId] = list.New()

	return mbx
}

// purgeMessage deletes all traces of specified message.
func (storage *Storage) purgeMessage(message *message.Message) {
	mbxId := storage.messageMailboxIds[message.ID]

	delete(storage.messageMailboxIds, message.ID)

	if element, ok := storage.messageElements[message.ID]; ok {
		storage.messageList[mbxId].Remove(element)

		delete(storage.messageElements, message.ID)
	}
}

// purgeMailbox deletes all traces of specified mailbox provided it has no registered messages.
func (storage *Storage) purgeMailbox(mbx *mailbox.Mailbox) {
	if messages, ok := storage.messageList[mbx.ID]; ok && messages.Len() > 0 {
		return
	}

	delete(storage.messageList, mbx.ID)

	if element, ok := storage.mailboxElements[mbx.ID]; ok {
		storage.mailboxList.Remove(element)

		delete(storage.mailboxElements, mbx.ID)
	}
}

// NewStorage creates new central storage structure.
func NewStorage() *Storage {
	return &Storage{
		mailboxMutex: sync.RWMutex{},
		messageMutex: sync.RWMutex{},

		mailboxList:     list.New(),
		mailboxElements: make(map[string]*list.Element),

		messageList:     make(map[string]*list.List),
		messageElements: make(map[string]*list.Element),

		messageMailboxIds: make(map[string]string),
	}
}

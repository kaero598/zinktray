package storage

import (
	"container/list"
	"go-fake-smtp/app/mailbox"
	"go-fake-smtp/app/message"
	"sync"
)

// Storage represents central storage for everything mail.
type Storage struct {
	mbxLock sync.RWMutex
	msgLock sync.RWMutex

	// Contains list of all registered mailboxes.
	mbxList *list.List

	// Maps mailbox ID to their respective list element.
	mbxElems map[string]*list.Element

	// Maps message ID to ID of the mailbox it is bound to.
	msgMbxIndex map[string]string

	// Maps mailbox ID to the list of messages belonging to it.
	msgs map[string]*list.List

	// Maps message IDs to their respective list elements.
	msgElems map[string]*list.Element
}

// Add stores new message and binds it to mailbox with provided ID.
func (storage *Storage) Add(msg *message.Message, mailboxID string) {
	storage.msgLock.Lock()

	mbx := storage.registerMailbox(mailboxID)

	if _, ok := storage.msgs[mailboxID]; !ok {
		storage.msgMbxIndex[msg.ID] = mbx.ID
		storage.msgElems[msg.ID] = storage.msgs[mailboxID].PushFront(msg)
	}

	storage.msgLock.Unlock()
}

// CountMailboxes returns the number of registered mailboxes.
func (storage *Storage) CountMailboxes() int {
	return storage.mbxList.Len()
}

// CountMessages returns the number of stored messages bound to specified mailbox.
func (storage *Storage) CountMessages(mailboxId string) int {
	storage.mbxLock.RLock()

	var length int

	if element, ok := storage.mbxElems[mailboxId]; ok {
		if mbx, ok := element.Value.(mailbox.Mailbox); ok {
			length = storage.msgs[mbx.ID].Len()
		}
	}

	storage.mbxLock.RUnlock()

	return length
}

// DeleteMailbox deletes registered mailbox along with all its messages.
func (storage *Storage) DeleteMailbox(mailboxId string) {
	storage.msgLock.Lock()
	storage.mbxLock.Lock()

	if element, ok := storage.mbxElems[mailboxId]; ok {
		if messages, ok := storage.msgs[mailboxId]; ok {
			next := messages.Front()

			for next != nil {
				if m, ok := next.Value.(*message.Message); ok {
					storage.purgeMessage(m)
				}
			}
		}

		if mbx, ok := element.Value.(*mailbox.Mailbox); ok {
			storage.purgeMailbox(mbx)
		}
	}

	storage.mbxLock.Unlock()
	storage.msgLock.Unlock()
}

// DeleteMessage deletes stored message.
func (storage *Storage) DeleteMessage(messageId string) {
	storage.msgLock.Lock()
	storage.mbxLock.Lock()

	if element, ok := storage.msgElems[messageId]; ok {
		if msg, ok := element.Value.(*message.Message); ok {
			mbxId := storage.msgMbxIndex[msg.ID]

			if element, ok := storage.mbxElems[mbxId]; ok {
				if mbx, ok := element.Value.(*mailbox.Mailbox); ok {
					storage.purgeMessage(msg)
					storage.purgeMailbox(mbx)
				}
			}
		}
	}

	storage.mbxLock.Unlock()
	storage.msgLock.Unlock()
}

// GetMailbox returns registered mailbox.
//
// Returns nil for unregistered mailboxes.
func (storage *Storage) GetMailbox(mailboxId string) *mailbox.Mailbox {
	storage.mbxLock.RLock()

	var mbx *mailbox.Mailbox

	if element, ok := storage.mbxElems[mailboxId]; ok {
		if m, ok := element.Value.(*mailbox.Mailbox); ok {
			mbx = m
		}
	}

	storage.mbxLock.RUnlock()

	return mbx
}

// GetMailboxes returns a list of all registered mailboxes.
func (storage *Storage) GetMailboxes() []*mailbox.Mailbox {
	storage.mbxLock.RLock()

	mailboxes := make([]*mailbox.Mailbox, 0, storage.mbxList.Len())
	next := storage.mbxList.Front()

	for next != nil {
		if mbx, ok := next.Value.(*mailbox.Mailbox); ok {
			mailboxes = append(mailboxes, mbx)
		}

		next = next.Next()
	}

	storage.mbxLock.RUnlock()

	return mailboxes
}

// GetMessage returns stored message.
//
// Returns nil for unknown message.
func (storage *Storage) GetMessage(messageId string) *message.Message {
	storage.msgLock.RLock()

	var msg *message.Message

	if v, ok := storage.msgElems[messageId]; ok {
		if m, ok := v.Value.(*message.Message); ok {
			msg = m
		}
	}

	storage.msgLock.RUnlock()

	return msg
}

// GetMessages returns a list of all known messages bound to specified mailbox.
func (storage *Storage) GetMessages(mailboxId string) []*message.Message {
	storage.mbxLock.RLock()
	storage.msgLock.RLock()

	var result []*message.Message

	if messages, ok := storage.msgs[mailboxId]; ok {
		result = make([]*message.Message, 0, messages.Len())
		next := messages.Front()

		for next != nil {
			if m, ok := next.Value.(*message.Message); ok {
				result = append(result, m)
			}

			next = next.Next()
		}
	} else {
		result = make([]*message.Message, 0)
	}

	storage.msgLock.RUnlock()
	storage.mbxLock.RUnlock()

	return result
}

// registerMailbox registers mailbox ID and returns corresponding mailbox.
//
// When mailbox ID is already registered returns that mailbox.
func (storage *Storage) registerMailbox(mailboxId string) *mailbox.Mailbox {
	var mbx *mailbox.Mailbox

	if element, ok := storage.mbxElems[mailboxId]; ok {
		if m, ok := element.Value.(*mailbox.Mailbox); ok {
			mbx = m
		}
	} else {
		mbx = mailbox.NewMailbox(mailboxId)

		storage.mbxElems[mbx.ID] = storage.mbxList.PushBack(mbx)
		storage.msgs[mailboxId] = list.New()
	}

	return mbx
}

// purgeMessage deletes all traces of specified message.
func (storage *Storage) purgeMessage(message *message.Message) {
	mbxId := storage.msgMbxIndex[message.ID]

	delete(storage.msgMbxIndex, message.ID)

	if element, ok := storage.msgElems[message.ID]; ok {
		storage.msgs[mbxId].Remove(element)

		delete(storage.msgElems, message.ID)
	}
}

// purgeMailbox deletes all traces of specified mailbox provided it has no registered messages.
func (storage *Storage) purgeMailbox(mbx *mailbox.Mailbox) {
	if messages, ok := storage.msgs[mbx.ID]; ok && messages.Len() > 0 {
		return
	}

	delete(storage.msgs, mbx.ID)

	if element, ok := storage.mbxElems[mbx.ID]; ok {
		storage.mbxList.Remove(element)

		delete(storage.mbxElems, mbx.ID)
	}
}

// NewStorage creates new central storage structure.
func NewStorage() *Storage {
	return &Storage{
		mbxLock: sync.RWMutex{},
		msgLock: sync.RWMutex{},

		mbxList:  list.New(),
		mbxElems: make(map[string]*list.Element),

		msgMbxIndex: make(map[string]string),

		msgs:     make(map[string]*list.List),
		msgElems: make(map[string]*list.Element),
	}
}

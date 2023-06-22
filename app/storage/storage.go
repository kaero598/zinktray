package storage

import (
	"go-fake-smtp/app/mailbox"
	"go-fake-smtp/app/message"
	"sync"
)

// Storage represents central storage for everything mail.
type Storage struct {
	mbxLock sync.RWMutex
	msgLock sync.RWMutex

	// Contains list of all registered mailboxes.
	mbxList []*mailbox.Mailbox

	// Maps mailbox ID to a list of IDs of all messages bound to that mailbox.
	mbxMsgList map[string][]string

	// Maps mailbox ID to mailbox itself.
	mbxIdIndex map[string]*mailbox.Mailbox

	// Maps message ID to message itself.
	msgIdIndex map[string]*message.Message

	// Maps message ID to ID of the mailbox it is bound to.
	msgMbxIndex map[string]string
}

// Add stores new message and binds it to mailbox with provided ID.
func (storage *Storage) Add(msg *message.Message, mailboxID string) {
	storage.msgLock.Lock()

	mbx := storage.registerMailbox(mailboxID)

	if _, ok := storage.msgIdIndex[msg.ID]; !ok {
		storage.mbxMsgList[mbx.ID] = append(storage.mbxMsgList[mbx.ID], msg.ID)
	}

	storage.msgIdIndex[msg.ID] = msg
	storage.msgMbxIndex[msg.ID] = mbx.ID

	storage.msgLock.Unlock()
}

// CountMailboxes returns the number of registered mailboxes.
func (storage *Storage) CountMailboxes() int {
	return len(storage.mbxList)
}

// CountMessages returns the number of stored messages bound to specified mailbox.
func (storage *Storage) CountMessages(mailboxId string) int {
	storage.mbxLock.RLock()

	var length int

	if mbx, ok := storage.mbxIdIndex[mailboxId]; ok {
		length = len(storage.mbxMsgList[mbx.ID])
	}

	storage.mbxLock.RUnlock()

	return length
}

// DeleteMailbox deletes registered mailbox along with all its messages.
func (storage *Storage) DeleteMailbox(mailboxId string) {
	storage.msgLock.Lock()
	storage.mbxLock.Lock()

	if mbx, ok := storage.mbxIdIndex[mailboxId]; ok {
		for _, msgId := range storage.mbxMsgList[mailboxId] {
			msg := storage.msgIdIndex[msgId]

			storage.purgeMessage(msg)
		}

		storage.purgeMailbox(mbx)
	}

	storage.mbxLock.Unlock()
	storage.msgLock.Unlock()
}

// DeleteMessage deletes stored message.
func (storage *Storage) DeleteMessage(messageId string) {
	storage.msgLock.Lock()
	storage.mbxLock.Lock()

	if msg, ok := storage.msgIdIndex[messageId]; ok {
		mbxId := storage.msgMbxIndex[msg.ID]
		mbx := storage.mbxIdIndex[mbxId]

		storage.purgeMessage(msg)
		storage.purgeMailbox(mbx)
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

	if v, ok := storage.mbxIdIndex[mailboxId]; ok {
		mbx = v
	}

	storage.mbxLock.RUnlock()

	return mbx
}

// GetMailboxes returns a list of all registered mailboxes.
func (storage *Storage) GetMailboxes() []*mailbox.Mailbox {
	storage.mbxLock.RLock()

	mailboxes := make([]*mailbox.Mailbox, 0, len(storage.mbxList))
	mailboxes = append(mailboxes, storage.mbxList...)

	storage.mbxLock.RUnlock()

	return mailboxes
}

// GetMessage returns stored message.
//
// Returns nil for unknown message.
func (storage *Storage) GetMessage(messageId string) *message.Message {
	storage.msgLock.RLock()

	var msg *message.Message

	if v, ok := storage.msgIdIndex[messageId]; ok {
		msg = v
	}

	storage.msgLock.RUnlock()

	return msg
}

// GetMessages returns a list of all known messages bound to specified mailbox.
func (storage *Storage) GetMessages(mailboxId string) []*message.Message {
	storage.mbxLock.RLock()
	storage.msgLock.RLock()

	var result []*message.Message

	if msgIdList, ok := storage.mbxMsgList[mailboxId]; ok {
		result = make([]*message.Message, 0, len(msgIdList))

		for _, msgId := range msgIdList {
			result = append(result, storage.msgIdIndex[msgId])
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

	if v, ok := storage.mbxIdIndex[mailboxId]; ok {
		mbx = v
	} else {
		mbx = mailbox.NewMailbox(mailboxId)

		storage.mbxList = append(storage.mbxList, mbx)
		storage.mbxIdIndex[mbx.ID] = mbx

		storage.mbxMsgList[mbx.ID] = make([]string, 0, 1)
	}

	return mbx
}

// purgeMessage deletes all traces of specified message.
func (storage *Storage) purgeMessage(message *message.Message) {
	mbxId := storage.msgMbxIndex[message.ID]

	delete(storage.msgIdIndex, message.ID)
	delete(storage.msgMbxIndex, message.ID)

	for k, v := range storage.mbxMsgList[mbxId] {
		if v == message.ID {
			updatedMsgList := make([]string, 0, len(storage.mbxMsgList[mbxId])-1)
			updatedMsgList = append(updatedMsgList, storage.mbxMsgList[mbxId][:k]...)
			updatedMsgList = append(updatedMsgList, storage.mbxMsgList[mbxId][k+1:]...)

			storage.mbxMsgList[mbxId] = updatedMsgList

			break
		}
	}
}

// purgeMailbox deletes all traces of specified mailbox provided it has no registered messages.
func (storage *Storage) purgeMailbox(mbx *mailbox.Mailbox) {
	if len(storage.mbxMsgList[mbx.ID]) > 0 {
		return
	}

	delete(storage.mbxIdIndex, mbx.ID)
	delete(storage.mbxMsgList, mbx.ID)

	for k, v := range storage.mbxList {
		if v.ID == mbx.ID {
			updatedMbxList := make([]*mailbox.Mailbox, 0, len(storage.mbxList)-1)
			updatedMbxList = append(updatedMbxList, storage.mbxList[:k]...)
			updatedMbxList = append(updatedMbxList, storage.mbxList[k+1:]...)

			storage.mbxList = updatedMbxList

			break
		}
	}
}

// NewStorage creates new central storage structure.
func NewStorage() *Storage {
	return &Storage{
		mbxLock: sync.RWMutex{},
		msgLock: sync.RWMutex{},

		mbxList:    make([]*mailbox.Mailbox, 0),
		mbxMsgList: make(map[string][]string),

		mbxIdIndex: make(map[string]*mailbox.Mailbox),

		msgIdIndex:  make(map[string]*message.Message),
		msgMbxIndex: make(map[string]string),
	}
}

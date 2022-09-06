package app

// Index of available mailboxes.
type MailboxIndex struct {
	mailboxes map[string][]*Message
}

// Adds new mailbox ID to the index.
func (index *MailboxIndex) Add(mailboxId string, message *Message) {
	if _, ok := index.mailboxes[mailboxId]; !ok {
		index.mailboxes[mailboxId] = make([]*Message, 0, 1)
	}

	index.mailboxes[mailboxId] = append(index.mailboxes[mailboxId], message)
}

// Returns the number of mailboxes in the index.
func (index *MailboxIndex) Count() int {
	return len(index.mailboxes)
}

// Returns a slice of all mailbox IDs in the index.
func (index *MailboxIndex) GetAll() []string {
	keys := make([]string, 0, len(index.mailboxes))

	for k := range index.mailboxes {
		keys = append(keys, k)
	}

	return keys
}

// Returns all messages of specified mailbox.
func (index *MailboxIndex) GetMessages(mailboxId string) []*Message {
	if _, ok := index.mailboxes[mailboxId]; ok {
		return index.mailboxes[mailboxId]
	}

	return make([]*Message, 0)
}

package app

// This empty struct is a zero-width value for mailbox index.
type void struct{}

// Index of available mailboxes.
type MailboxIndex struct {
	mailboxes map[string]void
}

// Adds new mailbox ID to the index.
func (index *MailboxIndex) Add(mailboxId string) {
	if _, ok := index.mailboxes[mailboxId]; !ok {
		index.mailboxes[mailboxId] = void{}
	}
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

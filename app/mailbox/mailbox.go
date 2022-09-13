package mailbox

import "go-fake-smtp/app/id"

type Mailbox struct {
	// Unique mailbox ID
	Id string

	// Name of the mailbox
	//
	// This does not serve any real purpose other than for API users to distinguish mailboxes with accessible way.
	//
	// Anonymous mailbox has empty name.
	Name string
}

// Creates new mailbox.
func NewMailbox(name string) *Mailbox {
	return &Mailbox{
		Id:   id.NewId(),
		Name: name,
	}
}

package app

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

func NewMailbox(name string) *Mailbox {
	return &Mailbox{
		Id:   NewId(),
		Name: name,
	}
}

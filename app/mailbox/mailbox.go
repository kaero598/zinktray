package mailbox

// Mailbox structure represents information on individual mailbox.
type Mailbox struct {
	// ID contains unique mailbox identifier.
	//
	// Usually this is the username provided during authentication.
	ID string
}

// NewMailbox creates new mailbox structure.
func NewMailbox(mailboxId string) *Mailbox {
	return &Mailbox{
		ID: mailboxId,
	}
}

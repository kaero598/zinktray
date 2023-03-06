package mailbox

// Essential mailbox details.
type MailboxDetails struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	IsAnonymous bool   `json:"isAnonymous"`
}

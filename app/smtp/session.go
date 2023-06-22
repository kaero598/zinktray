package smtp

import (
	"go-fake-smtp/app/message"
	"go-fake-smtp/app/storage"
	"io"

	"github.com/emersion/go-smtp"
)

// smtpSession represents information on individual SMTP session.
type smtpSession struct {
	// store provides central message storage.
	store *storage.Storage

	// mailboxName contains name of the mailbox in use.
	//
	// This is a username string SMTP client has authenticated with
	// or mailbox.Anonymous for anonymous session.
	mailboxName string
}

func (session *smtpSession) Mail(_ string, _ smtp.MailOptions) error {
	// Allow any "FROM" address (even malformed) since it is not used in any way.
	return nil
}

func (session *smtpSession) Rcpt(_ string) error {
	// Allow any "RCPT" address (even malformed) since it is not used in any way.
	return nil
}

func (session *smtpSession) Data(reader io.Reader) error {
	if buffer, err := io.ReadAll(reader); err != nil {
		return err
	} else {
		msg := message.NewMessage(string(buffer))

		session.store.Add(msg, session.mailboxName)
	}

	return nil
}

func (session *smtpSession) Reset() {
	// Nothing to reset.
}

func (session *smtpSession) Logout() error {
	// Nothing to clean-up.
	return nil
}

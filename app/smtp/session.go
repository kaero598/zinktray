package smtp

import (
	"go-fake-smtp/app/message"
	"go-fake-smtp/app/storage"
	"io"

	"github.com/emersion/go-smtp"
)

// Information on SMTP session.
type smtpSession struct {
	// Central message storage.
	storage *storage.Storage

	// Name of the mailbox.
	//
	// Contains authenticated username or empty string for anonymous session.
	mailboxName string
}

func (session *smtpSession) Mail(from string, opts smtp.MailOptions) error {
	// Allow any "FROM" address (even malformed) since it is not used in any way.
	return nil
}

func (session *smtpSession) Rcpt(rcpt string) error {
	// Allow any "RCPT" address (even malformed) since it is not used in any way.
	return nil
}

func (session *smtpSession) Data(reader io.Reader) error {
	if buffer, err := io.ReadAll(reader); err != nil {
		return err
	} else {
		message := message.NewMessage(string(buffer))

		session.storage.Add(message, session.mailboxName)
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

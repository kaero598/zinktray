package smtp

import (
	"errors"
	"io"
	"log"
	"zinktray/app/message"
	"zinktray/app/storage"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

var errInternal = errors.New("internal error")

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

func (session *smtpSession) Mail(_ string, _ *smtp.MailOptions) error {
	// Allow any "FROM" address (even malformed) since it is not used in any way.
	return nil
}

func (session *smtpSession) Rcpt(_ string, _ *smtp.RcptOptions) error {
	// Allow any "RCPT" address (even malformed) since it is not used in any way.
	return nil
}

func (session *smtpSession) Data(reader io.Reader) error {
	if buffer, err := io.ReadAll(reader); err != nil {
		return err
	} else {
		mbox := session.store.AddMailbox(session.mailboxName)
		msg := message.NewMessage(string(buffer))

		if err := session.store.AddMessage(msg, mbox.ID); err != nil {
			log.Printf("Cannot store message: %s", err)

			return errInternal
		}
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

func (session *smtpSession) Auth(mech string) (sasl.Server, error) {
	var authenticator = func(identity string, username string, password string) error {
		session.mailboxName = username

		return nil
	}

	return sasl.NewPlainServer(authenticator), nil
}

func (session *smtpSession) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

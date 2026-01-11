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

var errAuthenticationRequired = errors.New("authentication is required")
var errInternal = errors.New("internal error")
var errEmptyUsername = errors.New("username is mandatory")

// smtpSession represents information on individual SMTP session.
type smtpSession struct {
	// store provides central message storage.
	store *storage.Storage

	// mailboxID contains ID of the mailbox in use.
	mailboxID string
}

func (session *smtpSession) Mail(_ string, _ *smtp.MailOptions) error {
	if session.mailboxID == "" {
		return errAuthenticationRequired
	}

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
		mbox := session.store.AddMailbox(session.mailboxID)
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
		if username == "" {
			return errEmptyUsername
		}

		mbox := session.store.GetMailbox(username)
		if mbox == nil {
			mbox = session.store.AddMailbox(username)
		}

		session.mailboxID = mbox.ID

		return nil
	}

	return sasl.NewPlainServer(authenticator), nil
}

func (session *smtpSession) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

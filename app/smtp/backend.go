package smtp

import (
	"errors"
	"go-fake-smtp/app/storage"

	"github.com/emersion/go-smtp"
)

// Application SMTP backend.
type smtpBackend struct {
	// Central message storage
	storage *storage.Storage
}

func (backend *smtpBackend) AnonymousLogin(_ *smtp.ConnectionState) (smtp.Session, error) {
	// Allow anonymous login and store all messages in anonymous mailbox.
	return &smtpSession{storage: backend.storage}, nil
}

func (backend *smtpBackend) Login(_ *smtp.ConnectionState, username string, password string) (smtp.Session, error) {
	// Forbid empty login to prevent confusion with anonymous mailbox.
	if username == "" {
		return nil, errors.New("empty username is forbidden")
	}

	// Allow any other login since there are no mechanics to utilize logins yet.
	return &smtpSession{storage: backend.storage, mailboxName: username}, nil
}

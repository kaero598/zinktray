package smtp

import (
	"zinktray/app/mailbox"
	"zinktray/app/storage"

	"github.com/emersion/go-smtp"
)

// smtpBackend structure represents an application SMTP backend.
type smtpBackend struct {
	// store provides central message storage.
	store *storage.Storage
}

// AnonymousLogin creates a new session for clients without authentication.
//
// Anonymous session is operating a built-in anonymous mailbox.
func (backend *smtpBackend) AnonymousLogin(_ *smtp.ConnectionState) (smtp.Session, error) {
	return &smtpSession{
		store:       backend.store,
		mailboxName: mailbox.Anonymous,
	}, nil
}

// Login creates a new session for authenticated clients.
//
// Authenticated username is used as a name for operated mailbox.
func (backend *smtpBackend) Login(_ *smtp.ConnectionState, username string, _ string) (smtp.Session, error) {
	return &smtpSession{
		store:       backend.store,
		mailboxName: username,
	}, nil
}

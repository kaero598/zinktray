package smtp

import (
	"zinktray/app/storage"

	"github.com/emersion/go-smtp"
)

// smtpBackend structure represents an application SMTP backend.
type smtpBackend struct {
	// store provides central message storage.
	store *storage.Storage
}

func (b *smtpBackend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &smtpSession{
		store: b.store,
	}, nil
}

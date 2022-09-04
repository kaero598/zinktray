package app

import (
	"context"
	"io"
	"log"
	"sync"
	"time"

	"github.com/emersion/go-smtp"
)

// Information on SMTP session.
type SmtpSession struct {
	storage *Storage
}

func (session *SmtpSession) Mail(from string, opts smtp.MailOptions) error {
	// Allow any "FROM" address (even malformed) since it is not used in any way.
	return nil
}

func (session *SmtpSession) Rcpt(rcpt string) error {
	// Allow any "RCPT" address (even malformed) since it is not used in any way.
	return nil
}

func (session *SmtpSession) Data(reader io.Reader) error {
	if buffer, err := io.ReadAll(reader); err != nil {
		return err
	} else {
		message := Message{
			RawData: string(buffer),
		}

		session.storage.Backend.Add(message)
	}

	return nil
}

func (session *SmtpSession) Reset() {
	// Nothing to reset.
}

func (session *SmtpSession) Logout() error {
	// Nothing to clean-up.
	return nil
}

// Custom SMTP backend.
type SmtpBackend struct {
	storage *Storage
}

func (backend *SmtpBackend) AnonymousLogin(_ *smtp.ConnectionState) (smtp.Session, error) {
	// Anonymous login is forbidden until we decide how to deal with it.
	return nil, smtp.ErrAuthRequired
}

func (backend *SmtpBackend) Login(_ *smtp.ConnectionState, username string, password string) (smtp.Session, error) {
	// Allow any login since there are no mechanics to utilize it yet.
	return &SmtpSession{storage: backend.storage}, nil
}

// SMTP server.
//
// Starts SMTP backend and handles it's termination.
type SmtpServer struct {
	// Storage for received messages.
	storage *Storage
}

// Wires-up SMTP server.
//
// The backend is terminated as soon as ctx is cancelled.
func (srv *SmtpServer) Start(ctx context.Context, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	backend := &SmtpBackend{
		storage: srv.storage,
	}

	server := smtp.NewServer(backend)

	server.Addr = ":2525"
	server.Domain = "fake"
	server.ReadTimeout = 30 * time.Second
	server.WriteTimeout = 30 * time.Second
	server.AllowInsecureAuth = true
	server.MaxMessageBytes = 1024 * 1024
	server.MaxRecipients = 50

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("SMTP server failed to start: %s", err)
		}
	}()

	<-ctx.Done()

	if err := server.Close(); err != nil {
		log.Fatalf("Cannot shutdown SMTP server: %s", err)
	}
}

// Creates new SMTP server
func NewSmtpServer(storage *Storage) *SmtpServer {
	return &SmtpServer{
		storage: storage,
	}
}

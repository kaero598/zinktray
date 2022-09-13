package smtp

import (
	"context"
	"go-fake-smtp/app/storage"
	"log"
	"sync"
	"time"

	"github.com/emersion/go-smtp"
)

// SMTP server.
//
// Starts SMTP backend and handles it's termination.
type SmtpServer struct {
	// Storage for received messages.
	storage *storage.Storage
}

// Wires-up SMTP server.
//
// The backend is terminated as soon as ctx is cancelled.
func (srv *SmtpServer) Start(ctx context.Context, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	backend := &smtpBackend{
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
func NewServer(storage *storage.Storage) *SmtpServer {
	return &SmtpServer{
		storage: storage,
	}
}

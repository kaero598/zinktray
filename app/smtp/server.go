package smtp

import (
	"context"
	"go-fake-smtp/app/storage"
	"log"
	"sync"
	"time"

	"github.com/emersion/go-smtp"
)

// SmtpServer structure represents an SMTP server implementation.
//
// Handles start and termination of SMTP backend.
type SmtpServer struct {
	// store provides central message storage.
	store *storage.Storage
}

// Start wires-up SMTP server.
//
// The backend is terminated as soon as ctx is cancelled.
func (srv *SmtpServer) Start(ctx context.Context, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	backend := &smtpBackend{
		store: srv.store,
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

// NewServer creates new SMTP server structure.
func NewServer(storage *storage.Storage) *SmtpServer {
	return &SmtpServer{
		store: storage,
	}
}

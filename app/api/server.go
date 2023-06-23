package api

import (
	"context"
	"errors"
	context2 "go-fake-smtp/app/api/context"
	"go-fake-smtp/app/api/mailbox"
	"go-fake-smtp/app/api/message"
	"go-fake-smtp/app/storage"
	"log"
	"net/http"
	"sync"
	"time"
)

// Server structure represents HTTP API server.
type Server struct {
	storage *storage.Storage
}

// Start wires-up HTTP API server.
//
// The server is terminated as soon as ctx is cancelled.
func (srv *Server) Start(ctx context.Context, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	srv.addHandlers()

	server := &http.Server{
		Addr: "127.0.0.1:8080",
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("HTTP server failed to start: %s", err)
			}
		}
	}()

	<-ctx.Done()

	shutdownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		log.Fatal("Cannot shutdown HTTP server", err)
	}
}

// addHandlers registers HTTP API endpoints and their handlers.
func (srv *Server) addHandlers() {
	requestHandlerContext := &context2.RequestHandlerContext{
		Store: srv.storage,
	}

	http.Handle("/api/mailboxes/delete", mailbox.DeleteMailboxHandler(requestHandlerContext))
	http.Handle("/api/mailboxes/list", mailbox.GetMailboxListHandler(requestHandlerContext))

	http.Handle("/api/messages/delete", message.DeleteMessageHandler(requestHandlerContext))
	http.Handle("/api/messages/list", message.GetMessageListHandler(requestHandlerContext))
	http.Handle("/api/messages/details", message.GetMessageDetailsHandler(requestHandlerContext))
}

// NewServer creates new HTTP API server structure.
func NewServer(storage *storage.Storage) *Server {
	return &Server{
		storage: storage,
	}
}

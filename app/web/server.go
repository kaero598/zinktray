package web

import (
	"context"
	"go-fake-smtp/app/storage"
	"log"
	"net/http"
	"sync"
	"time"
)

// HTTP server.
//
// Serves API endpoints.
type WebServer struct {
	storage *storage.Storage
}

// Wires-up HTTP server.
//
// The server is terminated as soon as ctx is cancelled.
func (srv *WebServer) Start(ctx context.Context, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	srv.addHandlers()

	server := &http.Server{
		Addr: "127.0.0.1:8080",
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
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

// Registers API endpoints and their handlers.
func (srv *WebServer) addHandlers() {
	handler := &requestHandler{
		storage: srv.storage,
	}

	http.Handle("/", http.FileServer(http.Dir("public")))

	http.Handle("/api/mailboxes/delete", http.HandlerFunc(handler.deleteMailbox))
	http.Handle("/api/mailboxes/list", http.HandlerFunc(handler.getMailboxList))

	http.Handle("/api/messages/delete", http.HandlerFunc(handler.deleteMessage))
	http.Handle("/api/messages/list", http.HandlerFunc(handler.getMessageList))
	http.Handle("/api/messages/raw", http.HandlerFunc(handler.getMessageRawContents))
}

// Creates new HTTP server
func NewServer(storage *storage.Storage) *WebServer {
	return &WebServer{
		storage: storage,
	}
}

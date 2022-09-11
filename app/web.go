package app

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

// Essential message data to be returned via API.
type PublishedMessage struct {
	// Raw message contents along with body and headers.
	RawData string `json:"rawData"`
}

// Custom HTTP request handler.
type RequestHandler struct {
	// Received messages storage.
	storage *Storage
}

// Deletes mailbox along with all its messages.
//
// Expects "mailbox_id" form parameter. Returns HTTP 404 for unknown mailbox.
func (handler *RequestHandler) DeleteMailbox(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		response.WriteHeader(405)
		return
	}

	mailboxId := request.FormValue("mailbox_id")
	mailbox := handler.storage.GetMailbox(mailboxId)

	if mailbox != nil {
		handler.storage.DeleteMailbox(mailbox.Id)
	} else {
		response.WriteHeader(404)
	}
}

// Deletes message.
//
// Expects "message_id" form parameter. Returns HTTP 404 for unknown message.
func (handler *RequestHandler) DeleteMessage(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		response.WriteHeader(405)
		return
	}

	messageId := request.FormValue("message_id")
	message := handler.storage.GetMessage(messageId)

	if message != nil {
		handler.storage.DeleteMessage(message.Id)
	} else {
		response.WriteHeader(404)
	}
}

// Returns JSON-formatted list of IDs of all available mailboxes.
func (handler *RequestHandler) GetMailboxList(response http.ResponseWriter, request *http.Request) {
	publishList := make([]string, 0, handler.storage.CountMailboxes())

	for _, mailbox := range handler.storage.GetMailboxes() {
		publishList = append(publishList, mailbox.Id)
	}

	if encoded, err := json.Marshal(publishList); err != nil {
		log.Printf("Cannot encode mailbox list: %s\n", err)

		response.WriteHeader(http.StatusInternalServerError)
	} else {
		response.Header().Add("Content-Type", "application/json")
		response.Write(encoded)
	}
}

// Returns JSON-formatted list of all stored messages.
func (handler *RequestHandler) GetMessageList(response http.ResponseWriter, request *http.Request) {
	mailboxId := request.FormValue("mailbox_id")
	publishList := make([]PublishedMessage, 0, handler.storage.CountMessages(mailboxId))

	for _, msg := range handler.storage.GetMessages(mailboxId) {
		publishList = append(publishList, PublishedMessage{
			RawData: msg.RawData,
		})
	}

	if encoded, err := json.Marshal(publishList); err != nil {
		log.Printf("Cannot encode message list: %s\n", err)

		response.WriteHeader(http.StatusInternalServerError)
	} else {
		response.Header().Add("Content-Type", "application/json")
		response.Write(encoded)
	}
}

// HTTP server.
//
// Serves API endpoints.
type WebServer struct {
	storage *Storage
}

// Registers API endpoints and their handlers.
func (srv *WebServer) AddHandlers() {
	handler := &RequestHandler{
		storage: srv.storage,
	}

	http.Handle("/api/mailboxes/delete", http.HandlerFunc(handler.DeleteMailbox))
	http.Handle("/api/mailboxes/list", http.HandlerFunc(handler.GetMailboxList))

	http.Handle("/api/messages/delete", http.HandlerFunc(handler.DeleteMessage))
	http.Handle("/api/messages/list", http.HandlerFunc(handler.GetMessageList))
}

// Wires-up HTTP server.
//
// The server is terminated as soon as ctx is cancelled.
func (srv *WebServer) Start(ctx context.Context, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	srv.AddHandlers()

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

// Creates new HTTP server
func NewWebServer(storage *Storage) *WebServer {
	return &WebServer{
		storage: storage,
	}
}

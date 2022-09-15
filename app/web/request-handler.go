package web

import (
	"encoding/json"
	"go-fake-smtp/app/message"
	"go-fake-smtp/app/storage"
	"log"
	"net/http"
)

// Custom HTTP request handler.
type requestHandler struct {
	// Received messages storage.
	storage *storage.Storage
}

// Deletes mailbox along with all its messages.
//
// Expects "mailbox_id" form parameter. Returns HTTP 404 for unknown mailbox.
func (handler *requestHandler) deleteMailbox(response http.ResponseWriter, request *http.Request) {
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
func (handler *requestHandler) deleteMessage(response http.ResponseWriter, request *http.Request) {
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
func (handler *requestHandler) getMailboxList(response http.ResponseWriter, request *http.Request) {
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
func (handler *requestHandler) getMessageList(response http.ResponseWriter, request *http.Request) {
	mailboxId := request.FormValue("mailbox_id")
	publishList := make([]*message.MessageDetails, 0, handler.storage.CountMessages(mailboxId))

	for _, msg := range handler.storage.GetMessages(mailboxId) {
		publishList = append(publishList, message.Parse(msg))
	}

	if encoded, err := json.Marshal(publishList); err != nil {
		log.Printf("Cannot encode message list: %s\n", err)

		response.WriteHeader(http.StatusInternalServerError)
	} else {
		response.Header().Add("Content-Type", "application/json")
		response.Write(encoded)
	}
}

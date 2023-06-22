package web

import (
	"encoding/json"
	"go-fake-smtp/app/message/parse"
	"go-fake-smtp/app/storage"
	webmailbox "go-fake-smtp/app/web/mailbox"
	webmessage "go-fake-smtp/app/web/message"
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
		handler.storage.DeleteMailbox(mailbox.ID)
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
		handler.storage.DeleteMessage(message.ID)
	} else {
		response.WriteHeader(404)
	}
}

// Returns JSON-formatted list of all available mailboxes.
func (handler *requestHandler) getMailboxList(response http.ResponseWriter, request *http.Request) {
	publishList := make([]webmailbox.EssentialMailboxInfo, 0, handler.storage.CountMailboxes())

	for _, mbx := range handler.storage.GetMailboxes() {
		publishList = append(publishList, webmailbox.EssentialMailboxInfo{
			ID: mbx.ID,
		})
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
	publishList := make([]webmessage.EssentialMessageInfo, 0, handler.storage.CountMessages(mailboxId))

	for _, msg := range handler.storage.GetMessages(mailboxId) {
		if messageInfo, err := parse.ReadBasic(msg.GetRawData()); err == nil {
			publishList = append(publishList, webmessage.EssentialMessageInfo{
				ID:         msg.ID,
				From:       messageInfo.From,
				To:         messageInfo.To,
				Subject:    messageInfo.Subject,
				ReceivedAt: msg.ReceivedAt.Unix(),
			})
		} else {
			log.Printf("Cannot extract basic message info: %s\n", err)

			response.WriteHeader(http.StatusInternalServerError)

			return
		}
	}

	if encoded, err := json.Marshal(publishList); err != nil {
		log.Printf("Cannot encode message list: %s\n", err)

		response.WriteHeader(http.StatusInternalServerError)
	} else {
		response.Header().Add("Content-Type", "application/json")
		response.Write(encoded)
	}
}

// Returns JSON-formatted details of stored message.
func (handler *requestHandler) getMessageDetails(response http.ResponseWriter, request *http.Request) {
	messageId := request.FormValue("message_id")

	if msg := handler.storage.GetMessage(messageId); msg == nil {
		log.Printf("Message \"%s\" not found\n", messageId)

		response.WriteHeader(http.StatusNotFound)
	} else {
		var messageInfo *parse.BasicInfo
		var messageContent *parse.ContentInfo
		var err error

		messageInfo, err = parse.ReadBasic(msg.GetRawData())

		if err != nil {
			log.Printf("Cannot extract basic message info: %s\n", err)

			response.WriteHeader(http.StatusInternalServerError)
		}

		messageContent, err = parse.ReadContents(msg.GetRawData())

		if err != nil {
			log.Printf("Cannot extract message content: %s\n", err)

			response.WriteHeader(http.StatusInternalServerError)
		}

		publishInfo := webmessage.DetailedMessageInfo{
			ID:         msg.ID,
			From:       messageInfo.From,
			To:         messageInfo.To,
			Subject:    messageInfo.Subject,
			ReceivedAt: msg.ReceivedAt.Unix(),
			Content: webmessage.MessageContent{
				Raw:  msg.GetRawData(),
				Html: messageContent.Html,
				Text: messageContent.Plain,
			},
		}

		if encoded, err := json.Marshal(publishInfo); err != nil {
			log.Printf("Cannot encode message: %s\n", err)

			response.WriteHeader(http.StatusInternalServerError)
		} else {
			response.Header().Add("Content-Type", "application/json")
			response.Write(encoded)
		}
	}
}

// Returns raw message contents as received via SMTP session.
func (handler *requestHandler) getMessageRawContents(response http.ResponseWriter, request *http.Request) {
	messageId := request.FormValue("message_id")

	if message := handler.storage.GetMessage(messageId); message == nil {
		log.Printf("Message \"%s\" not found\n", messageId)

		response.WriteHeader(http.StatusNotFound)
	} else {
		response.Header().Add("Content-Type", "text/plain")
		response.Write([]byte(message.GetRawData()))
	}
}

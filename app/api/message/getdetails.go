package message

import (
	"encoding/json"
	"log"
	"net/http"
	"zinktray/app/api/context"
	"zinktray/app/message/parse"
)

// Returns JSON-formatted details of stored message.

// GetMessageDetailsHandler creates handler for detailed message information retrieval API.
//
// Detailed message information contains essential message information as returned by message list retrieval API,
// coupled with message contents.
//
// Expects "message_id" form parameter. Returns HTTP 404 Not Found for unknown message.
func GetMessageDetailsHandler(context *context.RequestHandlerContext) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		messageId := request.FormValue("message_id")

		if msg := context.Store.GetMessage(messageId); msg == nil {
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

			publishInfo := detailedMessageInfo{
				ID:         msg.ID,
				From:       messageInfo.From,
				To:         messageInfo.To,
				Subject:    messageInfo.Subject,
				ReceivedAt: msg.ReceivedAt.Unix(),
				Content: content{
					Raw:  msg.GetRawData(),
					Html: messageContent.Html,
					Text: messageContent.Plain,
				},
			}

			if encoded, err := json.Marshal(publishInfo); err != nil {
				log.Printf("Cannot encode message: %s\n", err)

				response.WriteHeader(http.StatusInternalServerError)
			} else {
				response.Header().Add("content-Type", "application/json")
				response.Write(encoded)
			}
		}
	}
}

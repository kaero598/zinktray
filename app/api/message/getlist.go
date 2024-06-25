package message

import (
	"encoding/json"
	"log"
	"net/http"
	"zinktray/app/api/context"
	"zinktray/app/message/parse"
)

// GetMessageListHandler creates handler for message list retrieval API.
//
// Message list contains essential information on each message stored in provided mailbox. To retrieve message contents
// use detailed message information retrieval API.
//
// Expects "mailbox_id" form parameter. Returns empty message list for unknown mailbox.
func GetMessageListHandler(context *context.RequestHandlerContext) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		mailboxId := request.FormValue("mailbox_id")
		publishList := make([]essentialMessageInfo, 0, context.Store.CountMessages(mailboxId))

		for _, msg := range context.Store.GetMessages(mailboxId) {
			if messageInfo, err := parse.ReadBasic(msg.GetRawData()); err == nil {
				publishList = append(publishList, essentialMessageInfo{
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
			response.Header().Add("content-Type", "application/json")
			response.Write(encoded)
		}
	}
}

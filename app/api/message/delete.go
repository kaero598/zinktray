package message

import (
	"net/http"
	"zinktray/app/api/context"
)

// DeleteMessageHandler creates handler for message deletion API.
//
// Expects "message_id" form parameter. Returns HTTP 404 Not Found for unknown message.
func DeleteMessageHandler(context *context.RequestHandlerContext) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			response.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		messageId := request.FormValue("message_id")
		message := context.Store.GetMessage(messageId)

		if message != nil {
			context.Store.DeleteMessage(message.ID)
		} else {
			response.WriteHeader(http.StatusNotFound)
		}
	}
}

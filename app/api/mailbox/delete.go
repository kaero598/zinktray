package mailbox

import (
	"net/http"
	"zinktray/app/api/context"
)

// DeleteMailboxHandler creates handler for mailbox deletion API.
//
// Expects "mailbox_id" form parameter. Returns HTTP 404 Not Found for unknown mailbox.
func DeleteMailboxHandler(context *context.RequestHandlerContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		mailboxId := request.FormValue("mailbox_id")
		mailbox := context.Store.GetMailbox(mailboxId)

		if mailbox != nil {
			context.Store.DeleteMailbox(mailbox.ID)
		} else {
			writer.WriteHeader(http.StatusNotFound)
		}
	}
}

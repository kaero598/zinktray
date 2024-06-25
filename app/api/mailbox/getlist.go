package mailbox

import (
	"encoding/json"
	"log"
	"net/http"
	"zinktray/app/api/context"
)

func GetMailboxListHandler(context *context.RequestHandlerContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		publishList := make([]essentialMailboxInfo, 0, context.Store.CountMailboxes())

		for _, mbx := range context.Store.GetMailboxes() {
			publishList = append(publishList, essentialMailboxInfo{
				ID: mbx.ID,
			})
		}

		if encoded, err := json.Marshal(publishList); err != nil {
			log.Printf("Cannot encode mailbox list: %s\n", err)

			writer.WriteHeader(http.StatusInternalServerError)
		} else {
			writer.Header().Add("Content-Type", "application/json")
			writer.Write(encoded)
		}
	}
}

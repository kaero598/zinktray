package message

import (
	"io"
	"log"
	"net/mail"
	"strings"
)

// Parses raw message contents into details structure.
func Parse(msg *Message) *MessageDetails {
	message, err := mail.ReadMessage(strings.NewReader(msg.GetRawData()))

	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(message.Body)

	if err != nil {
		log.Printf("Cannot read body of message %s", msg.Id)
	}

	return &MessageDetails{
		Id:         msg.Id,
		Subject:    message.Header.Get("Subject"),
		From:       extractAddressList(message, "From"),
		To:         extractAddressList(message, "To"),
		ReceivedAt: msg.ReceivedAt.Unix(),
		RawBody:    string(body),
	}
}

// Extracts addresses from message header.
func extractAddressList(message *mail.Message, headerKey string) []string {
	addressList, err := message.Header.AddressList(headerKey)

	if err != nil {
		log.Printf(
			"Cannot parse address list: %s. Raw header (%s): %s\n",
			err,
			headerKey,
			message.Header.Get(headerKey),
		)

		return make([]string, 0)
	}

	result := make([]string, 0)

	for _, address := range addressList {
		formattedAddress := "<" + address.Address + ">"

		if address.Name != "" {
			formattedAddress = address.Name + " " + formattedAddress
		}

		result = append(result, formattedAddress)
	}

	return result
}

package message

import "go-fake-smtp/app/id"

// Information on individual message.
type Message struct {
	// Unique message ID.
	Id string

	// Raw message contents along with body and headers as received via SMTP session.
	RawData string
}

// Creates new message.
func NewMessage(rawData string) *Message {
	return &Message{
		Id:      id.NewId(),
		RawData: rawData,
	}
}
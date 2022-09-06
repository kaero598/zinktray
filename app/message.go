package app

// Information on individual message.
type Message struct {
	// Unique message ID.
	Id string

	// Raw message contents along with body and headers as received via SMTP session.
	RawData string
}

func NewMessage(rawData string) *Message {
	return &Message{
		Id:      NewId(),
		RawData: rawData,
	}
}

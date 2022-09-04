package app

// Information on individual message.
type Message struct {
	// Raw message contents along with body and headers as received via SMTP session.
	RawData string
}

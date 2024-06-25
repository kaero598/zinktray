package message

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"time"
	"zinktray/app/id"
)

// Message structure represents information on individual message.
type Message struct {
	// ID contains unique message identifier.
	ID string

	// ReceivedAt contains time message has been received at.
	ReceivedAt time.Time

	// rawData contains raw message contents along with body and headers.
	//
	// Contents of rawData is compressed. Use GetRawData to read and SetRawData to write
	// uncompressed message contents.
	rawData string
}

// GetRawData reads uncompressed raw message contents.
func (msg *Message) GetRawData() string {
	rd, err := gzip.NewReader(strings.NewReader(msg.rawData))

	if err != nil {
		panic(err)
	}

	rawData, err := io.ReadAll(rd)

	if err != nil {
		panic(err)
	}

	rd.Close()

	return string(rawData)
}

// SetRawData writes uncompressed raw message contents.
func (msg *Message) SetRawData(rawData string) {
	var out bytes.Buffer

	writer := gzip.NewWriter(&out)

	_, err := writer.Write([]byte(rawData))

	if err != nil {
		panic(err)
	}

	writer.Close()

	msg.rawData = out.String()
}

// NewMessage creates new message structure.
func NewMessage(rawData string) *Message {
	msg := &Message{
		ID:         id.NewId(),
		ReceivedAt: time.Now(),
	}

	msg.SetRawData(rawData)

	return msg
}

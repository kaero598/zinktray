package message

import (
	"bytes"
	"compress/gzip"
	"go-fake-smtp/app/id"
	"io"
	"strings"
)

// Information on individual message.
type Message struct {
	// Unique message ID.
	Id string

	// Raw message contents along with body and headers.
	//
	// Compressed with gzip.
	rawData string
}

// Reads raw message contents.
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

// Writes raw message contents.
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

// Creates new message.
func NewMessage(rawData string) *Message {
	msg := &Message{
		Id: id.NewId(),
	}

	msg.SetRawData(rawData)

	return msg
}

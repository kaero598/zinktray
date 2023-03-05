package parse

import (
	"io"
	"mime/multipart"
	"net/mail"
)

type Part interface {
	GetHeader(name string) string

	GetReader() io.Reader
}

type MessagePart struct {
	msg *mail.Message
}

func (part *MessagePart) GetHeader(name string) string {
	return part.msg.Header.Get(name)
}

func (part *MessagePart) GetReader() io.Reader {
	return part.msg.Body
}

type MultipartPart struct {
	part *multipart.Part
}

func (part *MultipartPart) GetHeader(name string) string {
	return part.part.Header.Get(name)
}

func (part *MultipartPart) GetReader() io.Reader {
	return part.part
}

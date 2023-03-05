package parse

import (
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
)

// ReadBasic extracts basic message information from its raw body.
func ReadBasic(body string) (*BasicInfo, error) {
	msg, err := parseMessage(body)

	if err != nil {
		return nil, err
	}

	return &BasicInfo{
		Subject: msg.Header.Get("Subject"),
		From:    extractAddressList(msg, "From"),
		To:      extractAddressList(msg, "To"),
	}, nil
}

// ReadContents extract information about message contents from its raw body.
func ReadContents(body string) (*ContentInfo, error) {
	msg, err := parseMessage(body)

	if err != nil {
		return nil, err
	}

	contentPlain, contentHtml, err := readStructure(msg)

	if err != nil {
		return nil, err
	}

	if contentHtml == nil && contentPlain == nil {
		// Message must always have readable content, even if it is empty
		contentPlain = new(string)
	}

	return &ContentInfo{
		Html:  contentHtml,
		Plain: contentPlain,
		Raw:   body,
	}, nil
}

// parseMessage parses raw message body into readable structure.
func parseMessage(body string) (*mail.Message, error) {
	if msg, err := mail.ReadMessage(strings.NewReader(body)); err != nil {
		return nil, fmt.Errorf("cannot parse message body: %w", err)
	} else {
		return msg, nil
	}
}

func parsePart(part Part) (*string, *string, error) {
	var contentPlain, contentHtml *string

	partIndexStack := make([]uint, 0, 3)
	readersStack := make([]*multipart.Reader, 0, 3)

	var currentPart = part
	var currentReader *multipart.Reader
	var level uint
	var partIndex uint

ReadParts:
	for {
		if currentPart != nil {
			mediaType, boundary, err := extractMediaType(currentPart.GetHeader("Content-Type"))

			if err != nil {
				return nil, nil, fmt.Errorf(
					"cannot extract media type out of part %d, level %d: %w",
					partIndex,
					level,
					err,
				)
			}

			if isMultipart(mediaType) && boundary != "" {
				nextReader := multipart.NewReader(currentPart.GetReader(), boundary)

				if currentReader != nil {
					readersStack = append(readersStack, currentReader)
					partIndexStack = append(partIndexStack, partIndex)
				}

				currentReader = nextReader
				level++
				partIndex = 0
			} else if isHumanReadable(mediaType) {
				if content, err := io.ReadAll(currentPart.GetReader()); err != nil {
					return nil, nil, fmt.Errorf(
						"cannot read content of part %d, level %d: %w",
						partIndex,
						level,
						err,
					)
				} else if isHtml(mediaType) {
					contentHtml = appendString(contentHtml, string(content))
				} else {
					contentPlain = appendString(contentPlain, string(content))
				}
			} else {
				// todo: Probably we can treat this part as attachment
			}
		}

		if currentReader != nil {
			for {
				if nextPart, err := currentReader.NextPart(); err == io.EOF {
					if len(readersStack) == 0 {
						break ReadParts
					}

					currentReader = readersStack[len(readersStack)-1]
					partIndex = partIndexStack[len(partIndexStack)-1]

					readersStack = readersStack[:len(readersStack)-1]
					partIndexStack = partIndexStack[:len(partIndexStack)-1]

					level--
				} else if err != nil {
					return nil, nil, fmt.Errorf(
						"cannot read multipart contents of part %d, level %d: %w",
						partIndex,
						level,
						err,
					)
				} else {
					currentPart = &MultipartPart{part: nextPart}

					partIndex++

					break
				}
			}
		} else {
			break
		}
	}

	return contentPlain, contentHtml, nil
}

// readStructure reads message structure and returns aggregated contents. These include readable plain-text and HTML
// content.
func readStructure(msg *mail.Message) (*string, *string, error) {
	part := &MessagePart{
		msg: msg,
	}

	return parsePart(part)
}

// appendString appends an appendWhat string to a string pointed to by appendTo.
func appendString(appendTo *string, appendWhat string) *string {
	if appendTo == nil {
		return &appendWhat
	}

	*appendTo += appendWhat

	return appendTo
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

// extractMediaType retrieves media type and value of boundary parameter from Content-Type header value.
//
// Returned media type is normalized to lowercase.
func extractMediaType(contentType string) (string, string, error) {
	mediaType, params, err := mime.ParseMediaType(contentType)

	if err != nil {
		return "", "", fmt.Errorf("cannot extract boundary: %w", err)
	}

	return strings.ToLower(mediaType), params["boundary"], nil
}

// isHumanReadable tests whether media could be read by humans.
//
// mediaType is expected to be normalized to lowercase.
func isHumanReadable(mediaType string) bool {
	return isHtml(mediaType) || isPlainText(mediaType)
}

func isMultipart(mediaType string) bool {
	return strings.HasPrefix(mediaType, "multipart/")
}

// isHtml tests whether media is HTML source.
//
// mediaType is expected to be normalized to lowercase.
func isHtml(mediaType string) bool {
	return mediaType == "text/html"
}

// isPlainText tests whether media is plain text.
//
// mediaType is expected to be normalized to lowercase.
func isPlainText(mediaType string) bool {
	return mediaType == "text/plain"
}

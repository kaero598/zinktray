package message

// Essential message information.
type EssentialMessageInfo struct {
	Id         string   `json:"id"`
	From       []string `json:"from"`
	To         []string `json:"to"`
	Subject    string   `json:"subject"`
	ReceivedAt int64    `json:"receivedAt"`
}

// Detailed message information.
type DetailedMessageInfo struct {
	Id         string         `json:"id"`
	From       []string       `json:"from"`
	To         []string       `json:"to"`
	Subject    string         `json:"subject"`
	ReceivedAt int64          `json:"receivedAt"`
	Content    MessageContent `json:"content"`
}

// Extracted message content
type MessageContent struct {
	Raw  string  `json:"raw"`
	Html *string `json:"html"`
	Text *string `json:"text"`
}

package message

// essentialMessageInfo describes essential information on individual message to be exposed through HTTP API.
type essentialMessageInfo struct {
	ID         string   `json:"id"`
	From       []string `json:"from"`
	To         []string `json:"to"`
	Subject    string   `json:"subject"`
	ReceivedAt int64    `json:"receivedAt"`
}

// detailedMessageInfo describes full information on individual message to be exposed through HTTP API.
type detailedMessageInfo struct {
	ID         string   `json:"id"`
	From       []string `json:"from"`
	To         []string `json:"to"`
	Subject    string   `json:"subject"`
	ReceivedAt int64    `json:"receivedAt"`
	Content    content  `json:"content"`
}

// content describes contents of a message to be exposed through HTTP API.
type content struct {
	Raw  string  `json:"raw"`
	Html *string `json:"html"`
	Text *string `json:"text"`
}

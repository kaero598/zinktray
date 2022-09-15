package message

// Essential message details.
type MessageDetails struct {
	Id         string   `json:"id"`
	From       []string `json:"from"`
	To         []string `json:"to"`
	Subject    string   `json:"subject"`
	ReceivedAt int64    `json:"receivedAt"`
}

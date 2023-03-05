package parse

// BasicInfo contains basic message information.
//
// Intended for message introspection without delving into depths of its body (probably deep and complex).
type BasicInfo struct {
	From    []string
	To      []string
	Subject string
}

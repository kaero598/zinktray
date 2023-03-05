package parse

// ContentInfo contains information required to render message contents.
//
// todo: Implement attachments.
// todo: Implement proxying embedded images.
type ContentInfo struct {
	Html  *string // HTML message contents. nil means message has no HTML contents.
	Plain *string // Plain text message contents. nil means message has no plain-text contents.
	Raw   string  // Raw message body. Headers included.
}

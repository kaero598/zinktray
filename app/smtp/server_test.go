package smtp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/smtp"
	"os"
	"sync"
	"testing"
	"time"
	"zinktray/app/mailbox"
	"zinktray/app/storage"
)

var defaultURL = "127.0.0.1:2525"

var mailboxes = []string{
	"",
	"test1",
	"test2",
}

var messages = []messageSource{
	{mailboxes[0], "testdata/plain-text.txt"},
	{mailboxes[0], "testdata/text-html.txt"},
	{mailboxes[1], "testdata/multipart-alternative.txt"},
	{mailboxes[1], "testdata/multipart-mixed.txt"},
	{mailboxes[2], "testdata/multipart-related.txt"},
}

var messagesCount = map[string]int{
	mailboxes[0]: 2,
	mailboxes[1]: 2,
	mailboxes[2]: 1,
}

type messageSource struct {
	username string
	src      string
}

func TestSmtp(t *testing.T) {
	var storage = storage.NewStorage()
	var cancel = newServer(storage)

	t.Cleanup(cancel)

	var msg messageSource

	for _, msg = range messages {
		msg := msg

		t.Run("source", func(t *testing.T) {
			var err error
			var f *os.File

			f, err = os.Open(msg.src)

			if err != nil {
				t.Fatalf("Cannot open file \"%s\"", msg.src)
			}

			defer f.Close()

			var stat fs.FileInfo
			var size int64 = 4096

			stat, err = f.Stat()

			if err == nil {
				size = stat.Size()
			} else {
				t.Logf("Cannot stat file \"%s\": %s", msg.src, err)
			}

			var client = newClient()

			defer client.Close()

			if msg.username != "" {
				var auth = smtp.PlainAuth("", msg.username, "", "127.0.0.1")

				if err := client.Auth(auth); err != nil {
					t.Fatalf("Cannot authenticate: %s", err)
				}
			}

			err = client.Mail("test@localhost")

			if err != nil {
				t.Fatalf("Cannot issue SMTP MAIL command: %s", err)
			}

			err = client.Rcpt("test@localhost")

			if err != nil {
				t.Fatalf("Cannot issue SMTP RCPT command: %s", err)
			}

			var w io.WriteCloser

			w, err = client.Data()

			if err != nil {
				t.Fatalf("Cannot acquire writer for SMTP data: %s", err)
			}

			defer w.Close()

			var b = bytes.NewBuffer(make([]byte, size))

			_, err = b.ReadFrom(f)

			if err != nil {
				t.Fatalf("Cannot stream from file to SMTP data: %s", err)
			}

			_, err = w.Write(b.Bytes())

			if err != nil {
				t.Fatalf("Cannot write SMTP data: %s", err)
			}
		})
	}

	for mboxID, msgCountExpected := range messagesCount {
		if mboxID == "" {
			mboxID = mailbox.Anonymous
		}

		var msgCount = storage.CountMessages(mboxID)

		if msgCount != msgCountExpected {
			t.Fatalf(
				"Message count is wrong on mailbox \"%s\": got %d, expected %d",
				mboxID,
				msgCount,
				msgCountExpected,
			)
		}
	}
}

func newServer(storage *storage.Storage) context.CancelFunc {
	var ctx, cancel = context.WithCancel(context.Background())
	var server = NewServer(storage)
	var wg = &sync.WaitGroup{}

	wg.Add(1)

	go func() {
		server.Start(ctx, wg)
	}()

	return cancel
}

func newClient() *smtp.Client {
	var client *smtp.Client
	var err error

	for i := 3; i > 0; i-- {
		if client, err = smtp.Dial(defaultURL); err == nil {
			return client
		}

		time.Sleep(50 * time.Microsecond)
	}

	panic(fmt.Sprintf("Cannot dial %s: %s", defaultURL, err))
}

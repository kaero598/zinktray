package storage

import (
	"errors"
	"testing"
	"time"
	"zinktray/app/message"
)

func TestAdd(t *testing.T) {
	var storage = NewStorage()

	var mboxID = "mailbox_1"

	var mboxCount = storage.CountMailboxes()
	var mboxCountExpected = 0

	if mboxCount != mboxCountExpected {
		t.Fatalf("Mailbox count is wrong: got %d, expected %d", mboxCount, mboxCountExpected)
	}

	var msgCount = storage.CountMessages(mboxID)
	var msgCountExpected = 0

	if msgCount != msgCountExpected {
		t.Fatalf("Message count is wrong: got %d, expected %d", msgCount, msgCountExpected)
	}

	var msgExpected = &message.Message{ID: "message_1", ReceivedAt: time.Time{}}

	storage.Add(msgExpected, mboxID)

	mboxCount = storage.CountMailboxes()
	mboxCountExpected = 1

	if mboxCount != mboxCountExpected {
		t.Fatalf("Mailbox count is wrong: got %d, expected %d", mboxCount, mboxCountExpected)
	}

	msgCount = storage.CountMessages(mboxID)
	msgCountExpected = 1

	if msgCount != msgCountExpected {
		t.Fatalf("Message count is wrong: got %d, expected %d", msgCount, msgCountExpected)
	}

	var msg = storage.GetMessage(msgExpected.ID)

	if msg == nil {
		t.Fatalf("Message \"%s\" not found", msgExpected.ID)
	}

	if msg != msgExpected {
		t.Fatalf("Message \"%s\" is not the same as the one added", msgExpected.ID)
	}
}

func TestAddDuplicate(t *testing.T) {
	var storage = NewStorage()

	var mboxID_1 = "mailbox_1"
	var mboxID_2 = "mailbox_2"
	var msgID = "message_1"

	var msgCount = storage.CountMessages(mboxID_1)
	var msgCountExpected = 0

	if msgCount != msgCountExpected {
		t.Fatalf("Message count is wrong: got %d, expected %d", msgCount, msgCountExpected)
	}

	var err error
	var msg_1 = &message.Message{ID: msgID, ReceivedAt: time.Time{}}
	var msg_2 = &message.Message{ID: msgID, ReceivedAt: time.Time{}}

	err = storage.Add(msg_1, mboxID_1)

	if err != nil {
		t.Fatalf("Unexpected error upon adding message: %s", err)
	}

	msgCount = storage.CountMessages(mboxID_1)
	msgCountExpected = 1

	if msgCount != msgCountExpected {
		t.Fatalf("Message count is wrong: got %d, expected %d", msgCount, msgCountExpected)
	}

	err = storage.Add(msg_2, mboxID_2)

	if err == nil {
		t.Fatal("Unexpected adding result: expected error, got nil")
	} else if !errors.Is(err, ErrDuplicate) {
		t.Fatalf("Unexpected error: expected \"%s\", got \"%s\"", ErrDuplicate, err)
	}

	msgCount = storage.CountMessages(mboxID_1)
	msgCountExpected = 1

	if msgCount != msgCountExpected {
		t.Fatalf("Message count is wrong: got %d, expected %d", msgCount, msgCountExpected)
	}

	var msg = storage.GetMessage(msgID)

	if msg == nil {
		t.Fatalf("Message \"%s\" not found", msgID)
	}

	if msg != msg_1 {
		t.Fatalf("Message \"%s\" is not the same as the one added", msgID)
	}
}

func TestDeleteMailbox(t *testing.T) {
	var storage = NewStorage()

	var mboxID = "mailbox_1"
	var msgID = "message_1"

	storage.Add(&message.Message{ID: msgID, ReceivedAt: time.Time{}}, mboxID)
	storage.Add(&message.Message{ID: "message_2", ReceivedAt: time.Time{}}, mboxID)

	var mbox = storage.GetMailbox(mboxID)

	if mbox == nil {
		t.Fatalf("Mailbox \"%s\" not found", mboxID)
	}

	if mbox.ID != mboxID {
		t.Fatalf("Mailbox ID does not match: got \"%s\", expected \"%s\"", mbox.ID, mboxID)
	}

	storage.DeleteMessage(msgID)

	var msgCount = storage.CountMessages(mboxID)
	var msgCountExpected = 1

	if msgCount != msgCountExpected {
		t.Fatalf("Message count is wrong: got %d, expected %d", msgCount, msgCountExpected)
	}

	storage.DeleteMailbox(mboxID)

	mbox = storage.GetMailbox(mboxID)

	if mbox != nil {
		t.Fatalf("Mailbox \"%s\" found after deletion", mboxID)
	}

	msgCount = storage.CountMessages(mboxID)
	msgCountExpected = 0

	if msgCount != msgCountExpected {
		t.Fatalf("Message count is wrong: got %d, expected %d", msgCount, msgCountExpected)
	}
}

func TestDeleteMessage(t *testing.T) {
	var storage = NewStorage()

	var mboxID = "mailbox_1"
	var msgID_1 = "message_1"
	var msgID_2 = "message_2"

	for _, ID := range []string{msgID_1, msgID_2} {
		var msg = storage.GetMessage(ID)

		if msg != nil {
			t.Errorf("Message \"%s\" found before adding", ID)
		}
	}

	storage.Add(&message.Message{ID: msgID_1, ReceivedAt: time.Time{}}, mboxID)
	storage.Add(&message.Message{ID: msgID_2, ReceivedAt: time.Time{}}, mboxID)

	var mbox = storage.GetMailbox(mboxID)

	if mbox == nil {
		t.Fatalf("Mailbox \"%s\" not found", mboxID)
	}

	for _, ID := range []string{msgID_1, msgID_2} {
		storage.DeleteMessage(ID)

		var msg = storage.GetMessage(ID)

		if msg != nil {
			t.Errorf("Message \"%s\" found after deletion", ID)
		}
	}

	mbox = storage.GetMailbox(mboxID)

	if mbox != nil {
		t.Fatalf("Mailbox \"%s\" found after deletion of messages", mboxID)
	}
}

func TestGetMailbox(t *testing.T) {
	var storage = NewStorage()

	var mboxID = "mailbox_1"

	storage.Add(&message.Message{ID: "message_1", ReceivedAt: time.Time{}}, mboxID)

	var mbox = storage.GetMailbox(mboxID)

	if mbox == nil {
		t.Fatalf("Mailbox \"%s\" not found", mboxID)
	}

	if mbox.ID != mboxID {
		t.Fatalf("Mailbox ID does not match: got \"%s\", expected \"%s\"", mbox.ID, mboxID)
	}

	storage.DeleteMailbox(mboxID)

	mbox = storage.GetMailbox(mboxID)

	if mbox != nil {
		t.Fatalf("Mailbox \"%s\" found after deletion", mboxID)
	}
}

func TestGetMailboxes(t *testing.T) {
	var storage = NewStorage()

	var mboxID_1 = "mailbox_1"
	var mboxID_2 = "mailbox_2"

	storage.Add(&message.Message{ID: "message_1", ReceivedAt: time.Time{}}, mboxID_1)
	storage.Add(&message.Message{ID: "message_2", ReceivedAt: time.Time{}}, mboxID_2)

	var mboxCountExpected = 2
	var mboxCount = storage.CountMailboxes()

	if mboxCount != mboxCountExpected {
		t.Fatalf("Mailbox count does not match: got %d, expected %d", mboxCount, mboxCountExpected)
	}

	var mboxes = storage.GetMailboxes()

	if len(mboxes) != mboxCountExpected {
		t.Fatalf("Mailbox count does not match: got %d, expected %d", len(mboxes), mboxCountExpected)
	}

	for i, mboxID := range []string{mboxID_1, mboxID_2} {
		if mboxes[i].ID != mboxID {
			t.Errorf(
				"Mailbox ID does not match at index %d: got \"%s\", expected \"%s\"",
				i,
				mboxes[i].ID,
				mboxID,
			)
		}
	}
}

func TestGetMessages(t *testing.T) {
	var storage = NewStorage()

	var mboxID = "mailbox_1"
	var msgID_1 = "message_1"
	var msgID_2 = "message_2"

	storage.Add(&message.Message{ID: msgID_1, ReceivedAt: time.Time{}}, mboxID)
	storage.Add(&message.Message{ID: msgID_2, ReceivedAt: time.Time{}}, mboxID)

	var msgCountExpected = 2
	var msgCount = storage.CountMessages(mboxID)

	if msgCount != msgCountExpected {
		t.Fatalf("Message count does not match: got %d, expected %d", msgCount, msgCountExpected)
	}

	var msgs = storage.GetMessages(mboxID)

	if len(msgs) != msgCountExpected {
		t.Fatalf("Message count does not match: got %d, expected %d", len(msgs), msgCountExpected)
	}

	for i, msgID := range []string{msgID_2, msgID_1} {
		if msgs[i].ID != msgID {
			t.Errorf(
				"Message ID does not match at index %d: got \"%s\", expected \"%s\"",
				i,
				msgs[i].ID,
				msgID,
			)
		}
	}
}

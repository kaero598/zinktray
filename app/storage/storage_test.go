package storage

import (
	"errors"
	"testing"
	"time"
	"zinktray/app/message"
)

func TestAddMailbox(t *testing.T) {
	var storage = NewStorage()

	var mboxID1 = "test-mailbox-1"
	var mboxID2 = mboxID1
	var mboxID3 = "test-mailbox-3"

	var mbox1 = storage.AddMailbox(mboxID1)

	if mbox1 == nil {
		t.Fatal("Unexpected nil mailbox returned")
	}

	var mbox2 = storage.AddMailbox(mboxID2)

	if mbox1 != mbox2 {
		t.Fatal("Mailboxes are expected to be the same mailbox")
	}

	var mbox3 = storage.AddMailbox(mboxID3)

	if mbox3 == mbox1 {
		t.Fatal("Mailboxes are expected to no be the same mailbox")
	}

	if mboxCount := storage.CountMailboxes(); mboxCount != 2 {
		t.Fatalf("Wrong number of mailboxes: got %d, expected %d", mboxCount, 2)
	}
}

func TestAddMessage(t *testing.T) {
	var storage = NewStorage()

	var mboxID = "test-mailbox"

	var msgExpected = &message.Message{ID: "message_1", ReceivedAt: time.Time{}}

	if err := storage.AddMessage(msgExpected, mboxID); !errors.Is(err, ErrMailboxNotRegistered) {
		t.Fatalf("Operation result is expected to be \"%s\", got \"%s\"", ErrMailboxNotRegistered, err)
	}

	storage.AddMailbox(mboxID)

	if err := storage.AddMessage(msgExpected, mboxID); err != nil {
		t.Fatalf("Unexpected error upon adding message: %s", err)
	}

	if msg := storage.GetMessage(msgExpected.ID); msg == nil {
		t.Fatal("Message not found")
	} else if msg != msgExpected {
		t.Fatal("Message is not the same as the one added")
	}

	var msgCount = storage.CountMessages(mboxID)
	var msgCountExpected = 1

	if msgCount != msgCountExpected {
		t.Fatalf("Message count is wrong: got %d, expected %d", msgCount, msgCountExpected)
	}
}

func TestAddDuplicate(t *testing.T) {
	var storage = NewStorage()

	var mboxID1 = "mailbox_1"
	var mboxID2 = "mailbox_2"
	var msgID = "message_1"

	storage.AddMailbox(mboxID1)
	storage.AddMailbox(mboxID2)

	var msg1 = &message.Message{ID: msgID, ReceivedAt: time.Time{}}
	var msg2 = &message.Message{ID: msgID, ReceivedAt: time.Time{}}

	if err := storage.AddMessage(msg1, mboxID1); err != nil {
		t.Fatalf("Unexpected error upon adding message: %s", err)
	}

	if err := storage.AddMessage(msg2, mboxID2); err == nil {
		t.Fatal("Unexpected adding result: expected error, got nil")
	} else if !errors.Is(err, ErrDuplicate) {
		t.Fatalf("Unexpected error: expected \"%s\", got \"%s\"", ErrDuplicate, err)
	}
}

func TestDeleteMailbox(t *testing.T) {
	var storage = NewStorage()

	var mboxID = "mailbox_1"
	var msgID = "message_1"

	storage.AddMailbox(mboxID)

	_ = storage.AddMessage(&message.Message{ID: msgID, ReceivedAt: time.Time{}}, mboxID)
	_ = storage.AddMessage(&message.Message{ID: "message_2", ReceivedAt: time.Time{}}, mboxID)

	var mbox = storage.GetMailbox(mboxID)

	if mbox == nil {
		t.Fatalf("Mailbox \"%s\" not found", mboxID)
	}

	storage.DeleteMailbox(mboxID)

	if mbox = storage.GetMailbox(mboxID); mbox != nil {
		t.Fatalf("Mailbox \"%s\" found after deletion", mboxID)
	}

	if msgCount := storage.CountMessages(mboxID); msgCount > 0 {
		t.Fatalf("Mailbox appears to still have messages attached")
	}
}

func TestDeleteMessage(t *testing.T) {
	var storage = NewStorage()

	var mboxID = "mailbox_1"
	var msgID1 = "message_1"
	var msgID2 = "message_2"

	storage.AddMailbox(mboxID)

	_ = storage.AddMessage(&message.Message{ID: msgID1, ReceivedAt: time.Time{}}, mboxID)
	_ = storage.AddMessage(&message.Message{ID: msgID2, ReceivedAt: time.Time{}}, mboxID)

	for _, ID := range []string{msgID1, msgID2} {
		storage.DeleteMessage(ID)

		if msg := storage.GetMessage(ID); msg != nil {
			t.Errorf("Message \"%s\" found after deletion", ID)
		}
	}
}

func TestGetMailbox(t *testing.T) {
	var storage = NewStorage()

	var mboxID = "mailbox_1"

	storage.AddMailbox(mboxID)

	var mbox = storage.GetMailbox(mboxID)

	if mbox == nil {
		t.Fatalf("Mailbox \"%s\" not found", mboxID)
	}

	if mbox.ID != mboxID {
		t.Fatalf("Mailbox ID does not match: got \"%s\", expected \"%s\"", mbox.ID, mboxID)
	}
}

func TestGetMailboxes(t *testing.T) {
	var storage = NewStorage()

	var mboxID1 = "mailbox_1"
	var mboxID2 = "mailbox_2"

	storage.AddMailbox(mboxID1)
	storage.AddMailbox(mboxID2)

	var mboxCountExpected = 2
	var mboxCount = storage.CountMailboxes()

	if mboxCount != mboxCountExpected {
		t.Fatalf("Mailbox count does not match: got %d, expected %d", mboxCount, mboxCountExpected)
	}

	var mboxes = storage.GetMailboxes()

	if len(mboxes) != mboxCountExpected {
		t.Fatalf("Mailbox count does not match: got %d, expected %d", len(mboxes), mboxCountExpected)
	}

	for i, mboxID := range []string{mboxID1, mboxID2} {
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
	var msgID1 = "message_1"
	var msgID2 = "message_2"

	storage.AddMailbox(mboxID)

	_ = storage.AddMessage(&message.Message{ID: msgID1, ReceivedAt: time.Time{}}, mboxID)
	_ = storage.AddMessage(&message.Message{ID: msgID2, ReceivedAt: time.Time{}}, mboxID)

	var msgCountExpected = 2

	if msgCount := storage.CountMessages(mboxID); msgCount != msgCountExpected {
		t.Fatalf("Message count does not match: got %d, expected %d", msgCount, msgCountExpected)
	}

	var msgList = storage.GetMessages(mboxID)

	if len(msgList) != msgCountExpected {
		t.Fatalf("Message count does not match: got %d, expected %d", len(msgList), msgCountExpected)
	}

	for i, msgID := range []string{msgID2, msgID1} {
		if msgList[i].ID != msgID {
			t.Errorf(
				"Message ID does not match at index %d: got \"%s\", expected \"%s\"",
				i,
				msgList[i].ID,
				msgID,
			)
		}
	}
}

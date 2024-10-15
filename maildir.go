// The maildir package provides an interface to mailboxes in the Maildir format.
//
// Maildir mailboxes are designed to be safe for concurrent delivery. This
// means that at the same time, multiple processes can deliver to the same
// mailbox. However only one process can receive and read messages stored in
// the Maildir.
package maildir

import (
	"io"

	"github.com/emersion/go-maildir/internal"
)

type KeyError = internal.KeyError
type FlagError = internal.FlagError
type MailfileError = internal.MailfileError
type Flag = internal.Flag

// Message represents a message in a Maildir.
type Message struct {
	wrapped *internal.Message
}

// Filename returns the filesystem path to the message's file.
//
// The filename is not stable, it changes depending on the message flags.
func (msg *Message) Filename() string { return msg.wrapped.Filename() }

// Key returns the stable, unique identifier for the message.
func (msg *Message) Key() string { return msg.wrapped.Key() }

// Flags returns the message flags.
func (msg *Message) Flags() []Flag { return msg.wrapped.Flags() }

// SetFlags sets the message flags.
//
// Any duplicate flags are dropped, and flags are sorted before being saved.
func (msg *Message) SetFlags(flags []Flag) error { return msg.wrapped.SetFlags(flags) }

// Open reads the contents of a message.
func (msg *Message) Open() (io.ReadCloser, error) { return msg.wrapped.Open() }

// Remove deletes a message.
func (msg *Message) Remove() error { return msg.wrapped.Remove() }

// MoveTo moves a message from this Maildir to another one.
//
// The message flags are preserved, but its key might change.
func (msg *Message) MoveTo(target Dir) error { return msg.wrapped.MoveTo(target.wrapped) }

// CopyTo copies a message from this Maildir to another one.
//
// The copied message is returned. Its flags will be identical but its key
// might be different.
func (msg *Message) CopyTo(target Dir) (*Message, error) {
	m, err := msg.wrapped.CopyTo(target.wrapped)
	if err != nil {
		return nil, err
	}

	return fromInternalMessage(m), nil
}

// A Dir represents a single directory in a Maildir mailbox.
//
// Dir is used by programs receiving and reading messages from a Maildir. Only
// one process can perform these operations. Programs which only need to
// deliver new messages to the Maildir should use Delivery.
type Dir struct {
	wrapped internal.Dir
}

func NewDir(d string) *Dir {
	return &Dir{internal.Dir(d)}
}

// Init creates the directory structure for a Maildir.
//
// If the main directory already exists, it tries to create the subdirectories
// in there. If an error occurs while creating one of the subdirectories, this
// function may leave a partially created directory structure.
func (d Dir) Init() error { return d.wrapped.Init() }

// Create inserts a new message into the Maildir.
func (d *Dir) Create(flags []Flag) (*Message, io.WriteCloser, error) {
	m, w, err := d.wrapped.Create(flags, nil)
	if err != nil {
		return nil, nil, err
	}
	return fromInternalMessage(m), w, nil
}

// Clean removes old files from tmp and should be run periodically.
// This does not use access time but modification time for portability reasons.
func (d Dir) Clean() error { return d.wrapped.Clean() }

// Unseen moves messages from new to cur and returns them.
// This means the messages are now known to the application.
func (d Dir) Unseen() ([]*Message, error) {
	m, err := d.wrapped.Unseen()
	if err != nil {
		return nil, err
	}

	return fromInternalMessages(m), nil
}

// UnseenCount returns the number of messages in new without looking at them.
func (d Dir) UnseenCount() (int, error) { return d.wrapped.UnseenCount() }

// Walk calls fn for every message.
//
// If Walk encounters a malformed entry, it accumulates errors and continues
// iterating. If fn returns an error, Walk stops and returns a new error that
// contains fn's error in its tree (and can be checked via errors.Is).
func (d Dir) Walk(fn func(*Message) error) error {
	innerFn := func(msg *internal.Message) error {
		return fn(fromInternalMessage(msg))
	}

	return d.wrapped.Walk(innerFn)
}

// Messages returns a list of all messages in cur.
func (d Dir) Messages() ([]*Message, error) {
	m, err := d.wrapped.Messages()
	if err != nil {
		return nil, err
	}

	return fromInternalMessages(m), nil
}

// MessageByKey finds a message by key.
func (d Dir) MessageByKey(key string) (*Message, error) {
	m, err := d.wrapped.MessageByKey(key)
	if err != nil {
		return nil, err
	}

	return fromInternalMessage(m), nil
}

// Delivery represents an ongoing message delivery to the mailbox. It
// implements the io.WriteCloser interface. On Close the underlying file is
// moved/relinked to new.
//
// Multiple processes can perform a delivery on the same Maildir concurrently.
type Delivery struct {
	wrapped *internal.Delivery
}

// NewDelivery creates a new Delivery.
func NewDelivery(d string) (*Delivery, error) {
	dl, err := internal.NewDelivery(d, nil)
	if err != nil {
		return nil, err
	}

	return &Delivery{wrapped: dl}, nil
}

// fromInternalMessage converts an *internal.Message in a *Message,
// keeping the nil-ness, so:
//
//	fromInternalMessage(nil) == nil
func fromInternalMessage(msg *internal.Message) *Message {
	if msg == nil {
		return nil
	}

	return &Message{wrapped: msg}
}

// fromInternalMessages converts from []*internal.Message to []*Message
func fromInternalMessages(msgs []*internal.Message) []*Message {
	if msgs == nil {
		return nil
	}

	res := make([]*Message, len(msgs))
	for i, msg := range msgs {
		res[i] = fromInternalMessage(msg)
	}

	return res
}

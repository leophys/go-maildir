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
type Message = internal.Message

type Flag = internal.Flag

const (
	FlagPassed  Flag = internal.FlagPassed
	FlagReplied Flag = internal.FlagReplied
	FlagSeen    Flag = internal.FlagSeen
	FlagTrashed Flag = internal.FlagTrashed
	FlagDraft   Flag = internal.FlagDraft
	FlagFlagged Flag = internal.FlagFlagged
)

// A Dir represents a single directory in a Maildir mailbox.
//
// Dir is used by programs receiving and reading messages from a Maildir. Only
// one process can perform these operations. Programs which only need to
// deliver new messages to the Maildir should use Delivery.
type Dir struct {
	internal.Dir
}

func NewDir(d string) *Dir {
	return &Dir{internal.Dir(d)}
}

// Create inserts a new message into the Maildir.
func (d *Dir) Create(flags []Flag) (*Message, io.WriteCloser, error) {
	return d.Dir.Create(flags, nil)
}

// Delivery represents an ongoing message delivery to the mailbox. It
// implements the io.WriteCloser interface. On Close the underlying file is
// moved/relinked to new.
//
// Multiple processes can perform a delivery on the same Maildir concurrently.
type Delivery = internal.Delivery

// NewDelivery creates a new Delivery.
func NewDelivery(d string) (*Delivery, error) {
	return internal.NewDelivery(d, nil)
}

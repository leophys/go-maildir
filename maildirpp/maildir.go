// The maildirpp package mirrors the functionality of the maildir one, introducing
// support for attributes encoded in the name of the mail file.
package maildirpp

import (
	"io"

	"github.com/emersion/go-maildir/internal"
)

type KeyError = internal.KeyError
type FlagError = internal.FlagError
type MailfileError = internal.MailfileError
type Flag = internal.Flag
type Attributes = internal.Attributes
type Message = internal.Message

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
func (d *Dir) Create(flags []Flag, attrs Attributes) (*Message, io.WriteCloser, error) {
	return d.Dir.Create(flags, attrs)
}

// Delivery represents an ongoing message delivery to the mailbox. It
// implements the io.WriteCloser interface. On Close the underlying file is
// moved/relinked to new.
//
// Multiple processes can perform a delivery on the same Maildir concurrently.
type Delivery = internal.Delivery

// NewDelivery creates a new Delivery.
func NewDelivery(d string, attrs Attributes) (*Delivery, error) {
	return internal.NewDelivery(d, attrs)
}

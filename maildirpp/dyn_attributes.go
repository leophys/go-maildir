package maildirpp

import (
	"fmt"
	"sync/atomic"
)

// Computes the message size in bytes and returns it with the key S.
func DovecotMessageSize() DynAttribute {
	return &dynMsgSize{}
}

type dynMsgSize struct {
	count int
}

func (d *dynMsgSize) Write(p []byte) (n int, err error) {
	n = len(p)
	d.count += n
	return
}

func (d *dynMsgSize) Close() error {
	return nil
}

func (d *dynMsgSize) Key() string {
	return "S"
}

func (d *dynMsgSize) Compute() string {
	return fmt.Sprint(d.count)
}

var DovecotMessageRFC5322Size = DovecotMessageRFC822Size

// Computes the message size in bytes as if all the line feeds were CRLF,
// but counting the proper CRLF as 2 bytes, and returns it with the key W.
func DovecotMessageRFC822Size() DynAttribute {
	return &dynMsgRFC822Size{}
}

type dynMsgRFC822Size struct {
	count int
	carry bool
}

const (
	// this is carriage return \r
	CR = 13
	// this is line feed \n
	LF = 10
)

func (d *dynMsgRFC822Size) Write(p []byte) (n int, err error) {
	n = len(p)
	switch n {
	case 0:
		return
	case 1:
		n = 1
		if p[0] == CR {
			d.carry = true
			d.count += 1
			return
		}
		if p[0] == LF && d.carry {
			d.carry = false
			d.count += 1
			return
		}
		if p[0] == LF && !d.carry {
			d.count += 2
			return
		}
		d.count += 1
		return
	}
	// NOTE: here we know that the buffer p is long at least 2 bytes
	for i, b := range p {
		if i == 0 {
			d.count++
			continue
		}
		if b == LF {
			if p[i-1] == CR {
				d.count++
			} else {
				d.count += 2
			}
		} else {
			d.count++
		}
	}
	if p[n-1] == CR {
		d.carry = true
	}
	return
}

func (d *dynMsgRFC822Size) Close() error {
	return nil
}

func (d *dynMsgRFC822Size) Key() string {
	return "W"
}

func (d *dynMsgRFC822Size) Compute() string {
	return fmt.Sprint(d.count)
}

// PanicFinalizingNotClosed is an helper that wraps a DynAttribute
// and panics if its Compute method is called before Close.
// It is safe for concurrent use.
func PanicFinalizingNotClosed(d DynAttribute) DynAttribute {
	return &panicWrapperDynAttr{wrapped: d}
}

type panicWrapperDynAttr struct {
	wrapped DynAttribute
	closed  atomic.Bool
}

func (w *panicWrapperDynAttr) Write(p []byte) (n int, err error) {
	return w.wrapped.Write(p)
}

func (w *panicWrapperDynAttr) Close() error {
	w.closed.Store(true)
	return w.wrapped.Close()
}

func (w *panicWrapperDynAttr) Key() string {
	return w.wrapped.Key()
}

func (w *panicWrapperDynAttr) Compute() string {
	if !w.closed.Load() {
		panic("message not closed: cannot compute dynamic attributes")
	}
	return w.wrapped.Compute()
}

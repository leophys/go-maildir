package maildirpp

import (
	"fmt"
	"io"
	"slices"
	"testing"
)

func TestDovecotMessageSize(t *testing.T) {
	dynAttr := DovecotMessageSize()

	io.WriteString(dynAttr, "ciao")

	dynAttr.Close()

	count := dynAttr.Compute()
	if count != "4" {
		t.Fatalf("Unexpected count: %s", count)
	}

	key := dynAttr.Key()
	if key != "S" {
		t.Fatalf("Unexpected key: %s", key)
	}
}

func TestDovecotMessageRFC822Size(t *testing.T) {
	testCases := map[string]struct {
		msg      string
		expected int
	}{
		"zero":            {msg: "", expected: 0},
		"LF":              {msg: "\n", expected: 2},
		"CR":              {msg: "\r", expected: 1},
		"CRLF":            {msg: "\r\n", expected: 2},
		"CRLFCR":          {msg: "\r\n\r", expected: 3},
		"CRLFLF":          {msg: "\r\n\n", expected: 4},
		"test-carry-CR":   {msg: "message\ra", expected: 9},
		"test-carry-LF":   {msg: "message\na", expected: 10},
		"test-carry-CRLF": {msg: "message\r\n", expected: 9},
		"test-carry-CRCR": {msg: "message\r\r", expected: 9},
		"test-carry-LFLF": {msg: "message\n\n", expected: 11},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			dynAttr := DovecotMessageRFC822Size()

			for c := range slices.Chunk([]byte(tc.msg), 8) {
				_, _ = dynAttr.Write(c)
			}
			_ = dynAttr.Close()

			count := dynAttr.Compute()
			if count != fmt.Sprint(tc.expected) {
				t.Fatalf("expected count: %d, actual count: %s", tc.expected, count)
			}

			if key := dynAttr.Key(); key != "W" {
				t.Fatalf("unexpected key: %s", key)
			}
		})
	}
}

func TestPanicFinalizingNotClosed(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		dynAttr := DovecotMessageSize()
		dynAttr = PanicFinalizingNotClosed(dynAttr)

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("did not panic")
			}
		}()

		_ = dynAttr.Compute()
	})

	t.Run("proper", func(t *testing.T) {
		dynAttr := DovecotMessageSize()
		dynAttr = PanicFinalizingNotClosed(dynAttr)

		defer func() {
			if r := recover(); r != nil {
				t.Fatal("did panic")
			}
		}()

		io.WriteString(dynAttr, "ciao")
		_ = dynAttr.Close()

		count := dynAttr.Compute()
		if count != "4" {
			t.Fatalf("unexpected count: %s", count)
		}

		key := dynAttr.Key()
		if key != "S" {
			t.Fatalf("unexpected key: %s", key)
		}
	})
}

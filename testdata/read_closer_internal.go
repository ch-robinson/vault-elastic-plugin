package testdata

import (
	"bytes"
	"io"
)

type closingBuffer struct {
	*bytes.Buffer
}

// NewClosingBuffer instantiates a new ClosingBuffer
func newClosingBuffer(buf *bytes.Buffer) *closingBuffer {
	return &closingBuffer{buf}
}

// Close returns nil. This is just needed when reading bytes
// and can dispose when called from
func (cb *closingBuffer) Close() error {
	return nil
}

// GetReadCloser returns ClosingBuffer after the buffer is set
func (cb *closingBuffer) GetReadCloser() io.ReadCloser {
	return cb
}

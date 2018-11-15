package httputil

import (
	"bytes"
	"io"
)

// ClosingBuffer contains bytes.Buffer with a close function
type ClosingBuffer struct {
	*bytes.Buffer
}

// NewClosingBuffer instantiates a new ClosingBuffer
func NewClosingBuffer(buf *bytes.Buffer) *ClosingBuffer {
	return &ClosingBuffer{buf}
}

// Close returns nil. This is just needed when reading bytes
// and can dispose when called from
func (cb *ClosingBuffer) Close() error {
	return nil
}

// GetReadCloser returns ClosingBuffer after the buffer is set
func (cb *ClosingBuffer) GetReadCloser() io.ReadCloser {
	return cb
}

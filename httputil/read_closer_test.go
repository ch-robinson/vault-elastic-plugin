package httputil

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClosingBuffer(t *testing.T) {
	buf := bytes.NewBufferString("test")

	cb := NewClosingBuffer(buf)

	assert.Equal(t, buf, cb.Buffer)
}

package testdata

import "time"

type MockVaultContext struct {
}

func NewMockVaultContext() *MockVaultContext {
	return &MockVaultContext{}
}

func (*MockVaultContext) Done() <-chan struct{} {
	var c = make(chan struct{})
	return c
}

func (*MockVaultContext) Err() error {
	return nil
}

func (*MockVaultContext) Value(i interface{}) interface{} {
	return nil
}

func (*MockVaultContext) Deadline() (time.Time, bool) {
	return time.Now(), true
}

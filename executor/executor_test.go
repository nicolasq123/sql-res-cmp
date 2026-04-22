package executor

import "testing"

func TestNewDB_InvalidDSN(t *testing.T) {
	_, err := NewDB("invalid://dsn")
	if err == nil {
		t.Error("expected error for invalid dsn")
	}
}

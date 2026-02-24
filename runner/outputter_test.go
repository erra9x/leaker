package runner

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestWritePlainResult_NotVerbose(t *testing.T) {
	var buf bytes.Buffer
	err := WritePlainResult(&buf, false, "leakcheck", "user@example.com:password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if got != "user@example.com:password123\n" {
		t.Errorf("expected plain value line, got: %q", got)
	}
	if strings.Contains(got, "leakcheck") {
		t.Errorf("source should not appear in non-verbose output, got: %q", got)
	}
}

func TestWritePlainResult_Verbose(t *testing.T) {
	var buf bytes.Buffer
	err := WritePlainResult(&buf, true, "proxynova", "user@example.com:secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if got != "[proxynova] user@example.com:secret\n" {
		t.Errorf("expected verbose line with source, got: %q", got)
	}
}

// errWriter always returns an error on Write.
type errWriter struct{}

func (e *errWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("write error")
}

func TestWritePlainResult_PropagatesWriteError(t *testing.T) {
	err := WritePlainResult(&errWriter{}, false, "src", "value")
	if err == nil {
		t.Error("expected error from failing writer, got nil")
	}
}

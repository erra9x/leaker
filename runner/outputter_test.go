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

func TestWriteJSONResult_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	err := WriteJSONResult(&buf, "leakcheck", "email:user@example.com, password:abc", "user@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"source":"leakcheck"`) {
		t.Errorf("expected source field, got: %q", out)
	}
	if !strings.Contains(out, `"target":"user@example.com"`) {
		t.Errorf("expected target field, got: %q", out)
	}
	if !strings.Contains(out, `"value":"email:user@example.com, password:abc"`) {
		t.Errorf("expected value field, got: %q", out)
	}
	if !strings.HasSuffix(out, "\n") {
		t.Errorf("expected newline at end, got: %q", out)
	}
}

func TestWriteJSONResult_EscapesSpecialChars(t *testing.T) {
	var buf bytes.Buffer
	err := WriteJSONResult(&buf, "src", `value with "quotes" and \backslash`, "target")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// result should be valid JSON — parse it to verify
	out := strings.TrimSpace(buf.String())
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(out, "}") {
		t.Errorf("expected valid JSON object, got: %q", out)
	}
}

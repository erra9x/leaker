package runner

import (
	"bytes"
	"errors"
	"regexp"

	"github.com/vflame6/leaker/logger"
	"github.com/vflame6/leaker/runner/sources"
	"strings"
	"testing"
)

var ansiStripper = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func TestWritePlainResult_NotVerbose(t *testing.T) {
	var buf bytes.Buffer
	r := &sources.Result{Source: "leakcheck", Email: "user@example.com", Password: "password123"}
	err := WritePlainResult(&buf, false, false, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "email:user@example.com") {
		t.Errorf("expected email field, got: %q", got)
	}
	if !strings.Contains(got, "password:password123") {
		t.Errorf("expected password field, got: %q", got)
	}
	if strings.Contains(got, "leakcheck") {
		t.Errorf("source should not appear in non-verbose output, got: %q", got)
	}
}

func TestWritePlainResult_Verbose(t *testing.T) {
	var buf bytes.Buffer
	logger.SetNoColor(true)
	defer logger.SetNoColor(false)
	r := &sources.Result{Source: "proxynova", Email: "user@example.com", Password: "secret"}
	err := WritePlainResult(&buf, true, false, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := ansiStripper.ReplaceAllString(buf.String(), "")
	if !strings.HasPrefix(got, "[proxynova] ") {
		t.Errorf("expected verbose line with source prefix, got: %q", buf.String())
	}
}

// errWriter always returns an error on Write.
type errWriter struct{}

func (e *errWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("write error")
}

func TestWritePlainResult_PropagatesWriteError(t *testing.T) {
	r := &sources.Result{Source: "src", Email: "test@test.com"}
	err := WritePlainResult(&errWriter{}, false, false, r)
	if err == nil {
		t.Error("expected error from failing writer, got nil")
	}
}

func TestWriteJSONResult_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	r := &sources.Result{Source: "leakcheck", Email: "user@example.com", Password: "abc"}
	err := WriteJSONResult(&buf, false, r, "user@example.com")
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
	if !strings.Contains(out, `"email":"user@example.com"`) {
		t.Errorf("expected email field, got: %q", out)
	}
	if !strings.Contains(out, `"password":"abc"`) {
		t.Errorf("expected password field, got: %q", out)
	}
	if !strings.HasSuffix(out, "\n") {
		t.Errorf("expected newline at end, got: %q", out)
	}
}

func TestWritePlainResult_IncludeMetadata(t *testing.T) {
	var buf bytes.Buffer
	r := &sources.Result{Source: "leakcheck", Email: "user@example.com", Password: "abc", Database: "Canva.com"}
	err := WritePlainResult(&buf, false, true, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "database:Canva.com") {
		t.Errorf("expected database field with metadata flag, got: %q", got)
	}
}

func TestWritePlainResult_NoMetadataByDefault(t *testing.T) {
	var buf bytes.Buffer
	r := &sources.Result{Source: "leakcheck", Email: "user@example.com", Password: "abc", Database: "Canva.com"}
	err := WritePlainResult(&buf, false, false, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if strings.Contains(got, "database") {
		t.Errorf("database should not appear without metadata flag, got: %q", got)
	}
}

func TestWriteJSONResult_IncludeMetadata(t *testing.T) {
	var buf bytes.Buffer
	r := &sources.Result{Source: "leakcheck", Email: "user@example.com", Database: "Canva.com"}
	err := WriteJSONResult(&buf, true, r, "user@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"database":"Canva.com"`) {
		t.Errorf("expected database in JSON with metadata flag, got: %q", buf.String())
	}
}

func TestWriteJSONResult_NoMetadataByDefault(t *testing.T) {
	var buf bytes.Buffer
	r := &sources.Result{Source: "leakcheck", Email: "user@example.com", Database: "Canva.com"}
	err := WriteJSONResult(&buf, false, r, "user@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), `"database"`) {
		t.Errorf("database should not appear in JSON without metadata flag, got: %q", buf.String())
	}
}

func TestWriteJSONResult_EscapesSpecialChars(t *testing.T) {
	var buf bytes.Buffer
	r := &sources.Result{Source: "src", Password: `value with "quotes" and \backslash`}
	err := WriteJSONResult(&buf, false, r, "target")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// result should be valid JSON — parse it to verify
	out := strings.TrimSpace(buf.String())
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(out, "}") {
		t.Errorf("expected valid JSON object, got: %q", out)
	}
}

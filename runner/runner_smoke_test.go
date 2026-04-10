package runner

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/vflame6/leaker/runner/sources"
)

type smokeSource struct{}

func (s *smokeSource) Run(_ context.Context, target string, scanType sources.ScanType, _ *sources.Session) <-chan sources.Result {
	ch := make(chan sources.Result, 1)
	if scanType == sources.TypeEmail {
		ch <- sources.Result{
			Source:   s.Name(),
			Email:    target,
			Password: "hunter2",
		}
	}
	close(ch)
	return ch
}

func (s *smokeSource) Name() string          { return "smoke" }
func (s *smokeSource) UsesKey() bool         { return false }
func (s *smokeSource) NeedsKey() bool        { return false }
func (s *smokeSource) AddApiKeys(_ []string) {}
func (s *smokeSource) RateLimit() int        { return 1 }

func TestRunEnumeration_SmokePassiveEmailFlow(t *testing.T) {
	tmpDir := t.TempDir()
	targetsPath := filepath.Join(tmpDir, "targets.txt")
	if err := os.WriteFile(targetsPath, []byte("user@example.com\n"), 0o600); err != nil {
		t.Fatalf("write targets: %v", err)
	}

	var out bytes.Buffer
	r := &Runner{
		options: &Options{
			Targets:     targetsPath,
			Type:        sources.TypeEmail,
			Timeout:     2 * time.Second,
			Output:      &out,
			UserAgent:   "leaker-test",
			Version:     "test",
			NoRateLimit: true,
		},
		scanSources: []sources.Source{&smokeSource{}},
	}

	if err := r.RunEnumeration(context.Background()); err != nil {
		t.Fatalf("RunEnumeration error: %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "email:user@example.com") {
		t.Fatalf("expected smoke output to contain enumerated email, got %q", got)
	}
	if !strings.Contains(got, "password:hunter2") {
		t.Fatalf("expected smoke output to contain enumerated password, got %q", got)
	}
}

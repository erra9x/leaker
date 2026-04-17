package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Dir(file)
}

func runLeaker(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	cmd.Dir = repoRoot(t)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run %v failed: %v\n%s", args, err, out)
	}
	return string(out)
}

func readmeText(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot(t), "README.md"))
	if err != nil {
		t.Fatalf("read README: %v", err)
	}
	return string(data)
}

func sorted(items []string) []string {
	out := append([]string(nil), items...)
	sort.Strings(out)
	return out
}

func TestReadmeMatchesLiveSearchTypesAndSources(t *testing.T) {
	readme := readmeText(t)
	help := runLeaker(t, "--help")
	sourcesOutput := runLeaker(t, "--list-sources")

	commandRe := regexp.MustCompile(`(?m)^\s+(domain|email|keyword|phone|username)\s+Search by .+$`)
	matches := commandRe.FindAllStringSubmatch(help, -1)
	if len(matches) == 0 {
		t.Fatal("failed to find search commands in --help output")
	}
	searchTypes := make([]string, 0, len(matches))
	for _, m := range matches {
		searchTypes = append(searchTypes, m[1])
	}
	if !strings.Contains(readme, "**5 search types** - email, username, domain, keyword, phone") {
		t.Fatal("README search types summary drifted from expected CLI surface")
	}
	if got, want := sorted(searchTypes), sorted([]string{"domain", "email", "keyword", "phone", "username"}); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("unexpected live search types\ngot=%v\nwant=%v", got, want)
	}

	sourceLineRe := regexp.MustCompile(`(?m)^([a-z0-9]+)(?: \*)?$`)
	sourceMatches := sourceLineRe.FindAllStringSubmatch(sourcesOutput, -1)
	liveSources := make([]string, 0, len(sourceMatches))
	for _, m := range sourceMatches {
		liveSources = append(liveSources, m[1])
	}
	if !strings.Contains(readme, "**12 sources** - aggregates results from multiple leak databases") {
		t.Fatal("README source count summary drifted from expected CLI surface")
	}

	readmeSourceLinks := regexp.MustCompile(`\| \[([^\]]+)\]\(`).FindAllStringSubmatch(readme, -1)
	aliases := map[string]string{
		"intelligencex": "intelx",
	}
	readmeSources := make([]string, 0, len(readmeSourceLinks))
	for _, m := range readmeSourceLinks {
		name := strings.ToLower(strings.ReplaceAll(m[1], " ", ""))
		name = strings.ReplaceAll(name, "-", "")
		if alias, ok := aliases[name]; ok {
			name = alias
		}
		readmeSources = append(readmeSources, name)
	}
	if len(readmeSources) != len(liveSources) {
		t.Fatalf("README available sources table count drifted from live CLI\nreadme=%d live=%d", len(readmeSources), len(liveSources))
	}
	if strings.Join(sorted(readmeSources), ",") != strings.Join(sorted(liveSources), ",") {
		t.Fatalf("README available sources drifted from live CLI\nreadme=%v\nlive=%v", sorted(readmeSources), sorted(liveSources))
	}
}

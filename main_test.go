package main

import (
	"strings"
	"testing"
)

// Helper to call processEnvContent with default keepQuotes=false
func processEnvContentDefault(content string, compact bool) (string, error) {
	return processEnvContent(content, compact, false)
}

// Helper to call processEnvContent with keepQuotes=true
func processEnvContentKeepQuotes(content string, compact bool) (string, error) {
	return processEnvContent(content, compact, true)
}

func TestProcessEnvContent_ExpandsParsedVarsBeforeOS(t *testing.T) {
	t.Setenv("APP_NAME", "FromOS")

	in := "APP_NAME=FromFile\nMAIL_FROM_NAME=${APP_NAME}\n"
	out, err := processEnvContentDefault(in, false)
	if err != nil {
		t.Fatalf("processEnvContent returned error: %v", err)
	}

	expected := "APP_NAME=FromFile\nMAIL_FROM_NAME=FromFile\n"
	if out != expected {
		t.Fatalf("unexpected output\nexpected:\n%q\nactual:\n%q", expected, out)
	}
}

func TestProcessEnvContent_ExpandsBracedAndUnbraced(t *testing.T) {
	t.Setenv("HOST", "db.local")
	t.Setenv("PORT", "5432")

	in := "DATABASE_URL=postgres://$HOST:${PORT}/app\n"
	out, err := processEnvContentDefault(in, false)
	if err != nil {
		t.Fatalf("processEnvContent returned error: %v", err)
	}

	expected := "DATABASE_URL=postgres://db.local:5432/app\n"
	if out != expected {
		t.Fatalf("unexpected output\nexpected:\n%q\nactual:\n%q", expected, out)
	}
}

func TestProcessEnvContent_PreservesUnresolvedVariables(t *testing.T) {
	in := "GREETING=Hello $UNKNOWN\n"
	out, err := processEnvContentDefault(in, false)
	if err != nil {
		t.Fatalf("processEnvContent returned error: %v", err)
	}

	expected := "GREETING=Hello $UNKNOWN\n"
	if out != expected {
		t.Fatalf("unexpected output\nexpected:\n%q\nactual:\n%q", expected, out)
	}
}

func TestProcessEnvContent_StripsQuotesByDefault(t *testing.T) {
	in := "APP_NAME=Order Online\nMAIL_FROM_NAME=\"${APP_NAME}\"\n"
	out, err := processEnvContentDefault(in, false)
	if err != nil {
		t.Fatalf("processEnvContent returned error: %v", err)
	}

	expected := "APP_NAME=Order Online\nMAIL_FROM_NAME=Order Online\n"
	if out != expected {
		t.Fatalf("unexpected output\nexpected:\n%q\nactual:\n%q", expected, out)
	}
}

func TestProcessEnvContent_PreservesQuotesWithKeepQuotes(t *testing.T) {
	in := "APP_NAME=\"Order Online\"\nMAIL_FROM_NAME=\"${APP_NAME}\"\n"
	out, err := processEnvContentKeepQuotes(in, false)
	if err != nil {
		t.Fatalf("processEnvContent returned error: %v", err)
	}

	expected := "APP_NAME=\"Order Online\"\nMAIL_FROM_NAME=\"Order Online\"\n"
	if out != expected {
		t.Fatalf("unexpected output\nexpected:\n%q\nactual:\n%q", expected, out)
	}
}

func TestProcessEnvContent_CompactModeStripsCommentsAndBlankLines(t *testing.T) {
	in := "# header\n\nA=1\n  # inline comment line\nB = ${A}\n\n"
	out, err := processEnvContentDefault(in, true)
	if err != nil {
		t.Fatalf("processEnvContent returned error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	expectedLines := []string{"A=1", "B=1"}
	if len(lines) != len(expectedLines) {
		t.Fatalf("unexpected compact line count: got %d, want %d; output=%q", len(lines), len(expectedLines), out)
	}

	for i := range expectedLines {
		if lines[i] != expectedLines[i] {
			t.Fatalf("unexpected line %d: got %q want %q", i, lines[i], expectedLines[i])
		}
	}
}

func TestProcessEnvContent_CircularReferencesDoNotLoop(t *testing.T) {
	in := "A=$B\nB=$A\nC=$A\n"
	out, err := processEnvContentDefault(in, false)
	if err != nil {
		t.Fatalf("processEnvContent returned error: %v", err)
	}

	expected := "A=$B\nB=$B\nC=$B\n"
	if out != expected {
		t.Fatalf("unexpected output\nexpected:\n%q\nactual:\n%q", expected, out)
	}
}

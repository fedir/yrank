package main

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/fedir/yrank/youtube"
)

// resetFlags replaces the global flag set so tests don't share state.
func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

// --- filterFrom ---

func TestFilterFrom(t *testing.T) {
	date := func(s string) time.Time {
		t, _ := time.Parse("2006-01-02", s)
		return t
	}

	videos := []youtube.VideoStatistics{
		{Title: "old", PublishedAt: date("2022-06-01")},
		{Title: "boundary", PublishedAt: date("2024-01-01")},
		{Title: "recent", PublishedAt: date("2025-03-15")},
		{Title: "newest", PublishedAt: date("2026-01-10")},
	}

	tests := []struct {
		from  string
		want  []string
	}{
		{"2024-01-01", []string{"boundary", "recent", "newest"}},
		{"2025-01-01", []string{"recent", "newest"}},
		{"2026-06-01", []string{}},
		{"2020-01-01", []string{"old", "boundary", "recent", "newest"}},
	}

	for _, tc := range tests {
		fromDate, _ := time.Parse("2006-01-02", tc.from)
		got := filterFrom(append([]youtube.VideoStatistics(nil), videos...), fromDate)
		if len(got) != len(tc.want) {
			t.Errorf("from=%s: got %d results, want %d", tc.from, len(got), len(tc.want))
			continue
		}
		for i, v := range got {
			if v.Title != tc.want[i] {
				t.Errorf("from=%s [%d]: got %q, want %q", tc.from, i, v.Title, tc.want[i])
			}
		}
	}
}

// --- -from flag parsing ---

func TestFromFlag_valid(t *testing.T) {
	resetFlags()
	os.Args = []string{"yrank", "-p", "PLAYLIST", "-from", "2025-06-01"}
	_, _, _, _, _, from, _ := cliParameters()
	if from != "2025-06-01" {
		t.Errorf("expected from=2025-06-01, got %q", from)
	}
}

func TestFromFlag_empty(t *testing.T) {
	resetFlags()
	os.Args = []string{"yrank", "-p", "PLAYLIST"}
	_, _, _, _, _, from, _ := cliParameters()
	if from != "" {
		t.Errorf("expected empty from, got %q", from)
	}
}

package main

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/fedir/yrank/youtube"
)

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
		from string
		want []string
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

// --- -from flag ---

func TestFromFlag_valid(t *testing.T) {
	resetFlags()
	os.Args = []string{"yrank", "-p", "PLAYLIST", "-from", "2025-06-01"}
	_, _, _, _, _, from, _, _, _ := cliParameters()
	if from != "2025-06-01" {
		t.Errorf("expected from=2025-06-01, got %q", from)
	}
}

func TestFromFlag_empty(t *testing.T) {
	resetFlags()
	os.Args = []string{"yrank", "-p", "PLAYLIST"}
	_, _, _, _, _, from, _, _, _ := cliParameters()
	if from != "" {
		t.Errorf("expected empty from, got %q", from)
	}
}

// --- -strategy flag ---

func TestStrategyFlag_valid(t *testing.T) {
	resetFlags()
	os.Args = []string{"yrank", "-p", "PLAYLIST", "-strategy", "viral"}
	_, _, _, _, strategy, _, _, _, _ := cliParameters()
	if strategy != "viral" {
		t.Errorf("expected strategy=viral, got %q", strategy)
	}
}

func TestStrategyFlag_empty(t *testing.T) {
	resetFlags()
	os.Args = []string{"yrank", "-p", "PLAYLIST"}
	_, _, _, sorting, strategy, _, _, _, _ := cliParameters()
	if strategy != "" {
		t.Errorf("expected empty strategy, got %q", strategy)
	}
	if sorting != "total-interest" {
		t.Errorf("expected default sort=total-interest, got %q", sorting)
	}
}

// --- -weights flag parsing ---

func TestParseWeightsFlag(t *testing.T) {
	w := parseWeightsFlag("engagement=0.7,reach=0.2,comments=0.1")
	if w["engagement"] != 0.7 {
		t.Errorf("expected engagement=0.7, got %f", w["engagement"])
	}
	if w["reach"] != 0.2 {
		t.Errorf("expected reach=0.2, got %f", w["reach"])
	}
}

func TestParseWeightsFlag_empty(t *testing.T) {
	w := parseWeightsFlag("")
	if len(w) != 0 {
		t.Errorf("expected empty weights, got %v", w)
	}
}

func TestWeightsFlag_roundtrip(t *testing.T) {
	resetFlags()
	os.Args = []string{"yrank", "-p", "PLAYLIST", "-strategy", "viral", "-weights", "engagement=0.9,reach=0.05,comments=0.05"}
	_, _, _, _, _, _, weightsRaw, _, _ := cliParameters()
	w := parseWeightsFlag(weightsRaw)
	if w["engagement"] != 0.9 {
		t.Errorf("expected engagement=0.9 from CLI, got %f", w["engagement"])
	}
}

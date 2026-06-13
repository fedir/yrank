package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"os"
	"strings"
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

// --- filterByLength ---

func TestFilterByLength(t *testing.T) {
	videos := []youtube.VideoStatistics{
		{Title: "short", Duration: 30},
		{Title: "mid", Duration: 300},
		{Title: "long", Duration: 1800},
		{Title: "unknown", Duration: 0},
	}
	tests := []struct {
		name     string
		min, max int
		want     []string
	}{
		{"no limits", 0, 0, []string{"short", "mid", "long", "unknown"}},
		{"min only drops short and unknown", 120, 0, []string{"mid", "long"}},
		{"max only", 0, 600, []string{"short", "mid", "unknown"}},
		{"window", 60, 600, []string{"mid"}},
		{"min keeps boundary", 300, 0, []string{"mid", "long"}},
	}
	for _, tc := range tests {
		got := filterByLength(append([]youtube.VideoStatistics(nil), videos...), tc.min, tc.max)
		if len(got) != len(tc.want) {
			t.Errorf("%s: got %d results, want %d", tc.name, len(got), len(tc.want))
			continue
		}
		for i, v := range got {
			if v.Title != tc.want[i] {
				t.Errorf("%s [%d]: got %q, want %q", tc.name, i, v.Title, tc.want[i])
			}
		}
	}
}

// --- filterByViews ---

func TestFilterByViews(t *testing.T) {
	videos := []youtube.VideoStatistics{
		{Title: "low", ViewCount: 50},
		{Title: "mid", ViewCount: 1000},
		{Title: "high", ViewCount: 100000},
	}
	tests := []struct {
		name string
		min  int
		want []string
	}{
		{"no min", 0, []string{"low", "mid", "high"}},
		{"drops low", 1000, []string{"mid", "high"}},
		{"keeps boundary", 100000, []string{"high"}},
		{"all dropped", 200000, []string{}},
	}
	for _, tc := range tests {
		got := filterByViews(append([]youtube.VideoStatistics(nil), videos...), tc.min)
		if len(got) != len(tc.want) {
			t.Errorf("%s: got %d results, want %d", tc.name, len(got), len(tc.want))
			continue
		}
		for i, v := range got {
			if v.Title != tc.want[i] {
				t.Errorf("%s [%d]: got %q, want %q", tc.name, i, v.Title, tc.want[i])
			}
		}
	}
}

// --- -from flag ---

func TestFromFlag_valid(t *testing.T) {
	resetFlags()
	os.Args = []string{"yrank", "-p", "PLAYLIST", "-from", "2025-06-01"}
	_, _, _, _, _, _, from, _, _, _, _, _, _, _, _ := cliParameters()
	if from != "2025-06-01" {
		t.Errorf("expected from=2025-06-01, got %q", from)
	}
}

func TestFromFlag_empty(t *testing.T) {
	resetFlags()
	os.Args = []string{"yrank", "-p", "PLAYLIST"}
	_, _, _, _, _, _, from, _, _, _, _, _, _, _, _ := cliParameters()
	if from != "" {
		t.Errorf("expected empty from, got %q", from)
	}
}

// --- -min-length / -max-length / -min-views flags ---

func TestLengthViewFlags(t *testing.T) {
	resetFlags()
	os.Args = []string{"yrank", "-p", "PLAYLIST", "-min-length", "120", "-max-length", "600", "-min-views", "5000"}
	_, _, _, _, _, _, _, _, _, _, minLen, maxLen, minViews, _, _ := cliParameters()
	if minLen != 120 {
		t.Errorf("expected min-length=120, got %d", minLen)
	}
	if maxLen != 600 {
		t.Errorf("expected max-length=600, got %d", maxLen)
	}
	if minViews != 5000 {
		t.Errorf("expected min-views=5000, got %d", minViews)
	}
}

func TestLengthViewFlags_defaults(t *testing.T) {
	resetFlags()
	os.Args = []string{"yrank", "-p", "PLAYLIST"}
	_, _, _, _, _, _, _, _, _, _, minLen, maxLen, minViews, _, _ := cliParameters()
	if minLen != 0 || maxLen != 0 || minViews != 0 {
		t.Errorf("expected all length/view flags to default to 0, got min-length=%d max-length=%d min-views=%d", minLen, maxLen, minViews)
	}
}

// --- -strategy flag ---

func TestStrategyFlag_valid(t *testing.T) {
	resetFlags()
	os.Args = []string{"yrank", "-p", "PLAYLIST", "-strategy", "viral"}
	_, _, _, _, _, strategy, _, _, _, _, _, _, _, _, _ := cliParameters()
	if strategy != "viral" {
		t.Errorf("expected strategy=viral, got %q", strategy)
	}
}

func TestStrategyFlag_empty(t *testing.T) {
	resetFlags()
	os.Args = []string{"yrank", "-p", "PLAYLIST"}
	_, _, _, _, sorting, strategy, _, _, _, _, _, _, _, _, _ := cliParameters()
	if strategy != "" {
		t.Errorf("expected empty strategy, got %q", strategy)
	}
	if sorting != "total-interest" {
		t.Errorf("expected default sort=total-interest, got %q", sorting)
	}
}

// --- -top-search flag ---

func TestTopSearchFlag_valid(t *testing.T) {
	resetFlags()
	os.Args = []string{"yrank", "-top-search", "kubernetes operator"}
	_, _, topSearch, _, sorting, _, _, _, _, _, _, _, _, _, _ := cliParameters()
	if topSearch != "kubernetes operator" {
		t.Errorf("expected top-search=%q, got %q", "kubernetes operator", topSearch)
	}
	// With no -s/-strategy, the default sort still applies.
	if sorting != "total-interest" {
		t.Errorf("expected default sort=total-interest, got %q", sorting)
	}
}

func TestTopSearchFlag_empty(t *testing.T) {
	resetFlags()
	os.Args = []string{"yrank", "-p", "PLAYLIST"}
	_, _, topSearch, _, _, _, _, _, _, _, _, _, _, _, _ := cliParameters()
	if topSearch != "" {
		t.Errorf("expected empty top-search, got %q", topSearch)
	}
}

// --- envWeights ---

func TestEnvWeights_parsed(t *testing.T) {
	os.Setenv("WEIGHT_VIRAL_ENGAGEMENT", "0.8")
	os.Setenv("WEIGHT_HYPE_VELOCITY", "2.0")
	defer os.Unsetenv("WEIGHT_VIRAL_ENGAGEMENT")
	defer os.Unsetenv("WEIGHT_HYPE_VELOCITY")

	w := envWeights()
	if w["viral_engagement"] != 0.8 {
		t.Errorf("expected viral_engagement=0.8, got %f", w["viral_engagement"])
	}
	if w["hype_velocity"] != 2.0 {
		t.Errorf("expected hype_velocity=2.0, got %f", w["hype_velocity"])
	}
}

func TestEnvWeights_ignoresUnrelated(t *testing.T) {
	os.Setenv("YOUTUBE_API_KEY", "dummy")
	defer os.Unsetenv("YOUTUBE_API_KEY")
	w := envWeights()
	if _, ok := w["youtube_api_key"]; ok {
		t.Error("envWeights should not include YOUTUBE_API_KEY")
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
	_, _, _, _, _, _, _, weightsRaw, _, _, _, _, _, _, _ := cliParameters()
	w := parseWeightsFlag(weightsRaw)
	if w["engagement"] != 0.9 {
		t.Errorf("expected engagement=0.9 from CLI, got %f", w["engagement"])
	}
}

// --- CSV output ---

func sampleVideos() []youtube.VideoStatistics {
	pub, _ := time.Parse("2006-01-02", "2025-01-15")
	return []youtube.VideoStatistics{
		{
			Title:                       "Video One",
			URL:                         "https://youtu.be/aaa",
			PublishedAt:                 pub,
			Duration:                    245,
			PositiveInterestingness:     0.0512,
			PositiveNegativeCoefficient: 1234.0,
			TotalInterestingness:        0.0530,
			ViewCount:                   10000,
			LikeCount:                   512,
			DislikeCount:                0,
			CommentCount:                100,
			TotalReaction:               612,
			GlobalBuzzIndex:             6120000,
		},
		{Title: ""}, // blank title must be skipped
	}
}

func TestPrintTo_CSV_headers(t *testing.T) {
	var buf bytes.Buffer
	printTo(&buf, sampleVideos(), "csv", false, false)

	r := csv.NewReader(strings.NewReader(buf.String()))
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("invalid CSV: %v", err)
	}
	if len(records) == 0 {
		t.Fatal("expected at least a header row")
	}
	want := []string{"Title", "URL", "Published at", "Duration", "Positive interestingness",
		"Positive negative coefficient", "Total interestingness",
		"Views", "Likes", "Dislikes", "Comments", "Total reaction", "Global buzz index"}
	for i, h := range want {
		if records[0][i] != h {
			t.Errorf("header[%d]: got %q want %q", i, records[0][i], h)
		}
	}
}

func TestPrintTo_CSV_row_count(t *testing.T) {
	var buf bytes.Buffer
	printTo(&buf, sampleVideos(), "csv", false, false)

	r := csv.NewReader(strings.NewReader(buf.String()))
	records, _ := r.ReadAll()
	// 1 header + 1 real video (blank-title entry must be skipped)
	if len(records) != 2 {
		t.Errorf("expected 2 rows (header+1 video), got %d", len(records))
	}
}

func TestPrintTo_CSV_values(t *testing.T) {
	var buf bytes.Buffer
	printTo(&buf, sampleVideos(), "csv", false, false)

	r := csv.NewReader(strings.NewReader(buf.String()))
	records, _ := r.ReadAll()
	row := records[1]
	if row[0] != "Video One" {
		t.Errorf("Title: got %q want %q", row[0], "Video One")
	}
	if row[3] != "245" {
		t.Errorf("Duration: got %q want %q", row[3], "245")
	}
	if row[7] != "10000" {
		t.Errorf("Views: got %q want %q", row[7], "10000")
	}
	if row[4] != "0.0512" {
		t.Errorf("PositiveInterestingness: got %q want %q", row[4], "0.0512")
	}
}

func TestPrintToFile_atomic(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/out.csv"
	tmp := path + ".tmp"

	if err := printToFile(path, sampleVideos(), "csv", false, false); err != nil {
		t.Fatalf("printToFile: %v", err)
	}
	// final file must exist
	if _, err := os.Stat(path); err != nil {
		t.Errorf("output file missing: %v", err)
	}
	// temp file must be gone
	if _, err := os.Stat(tmp); err == nil {
		t.Error("temp file should have been removed after rename")
	}
	// content must be valid CSV
	data, _ := os.ReadFile(path)
	r := csv.NewReader(strings.NewReader(string(data)))
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("invalid CSV in output file: %v", err)
	}
	if len(records) != 2 {
		t.Errorf("expected 2 rows, got %d", len(records))
	}
}

func TestPrintTo_CSV_specialChars(t *testing.T) {
	pub, _ := time.Parse("2006-01-02", "2025-01-15")
	videos := []youtube.VideoStatistics{
		{Title: `Title with "quotes"`, URL: "https://youtu.be/a", PublishedAt: pub},
		{Title: "Title with, comma", URL: "https://youtu.be/b", PublishedAt: pub},
		{Title: "Title with\nnewline", URL: "https://youtu.be/c", PublishedAt: pub},
		{Title: "Title with | pipe", URL: "https://youtu.be/d", PublishedAt: pub},
		{Title: "Emoji 😱🎯", URL: "https://youtu.be/e", PublishedAt: pub},
	}
	var buf bytes.Buffer
	printTo(&buf, videos, "csv", false, false)

	r := csv.NewReader(strings.NewReader(buf.String()))
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("CSV with special chars is invalid: %v", err)
	}
	// header + 5 data rows
	if len(records) != 6 {
		t.Fatalf("expected 6 rows, got %d", len(records))
	}
	titles := []string{
		`Title with "quotes"`,
		"Title with, comma",
		"Title with\nnewline",
		"Title with | pipe",
		"Emoji 😱🎯",
	}
	for i, want := range titles {
		if got := records[i+1][0]; got != want {
			t.Errorf("row %d title: got %q want %q", i+1, got, want)
		}
	}
}

func TestMdSafe(t *testing.T) {
	cases := []struct{ in, want string }{
		{"no pipes here", "no pipes here"},
		{"a | b | c", `a \| b \| c`},
		{"| leading", `\| leading`},
		{"trailing |", `trailing \|`},
		{"Istio Day: Zero-Downtime | Migration", `Istio Day: Zero-Downtime \| Migration`},
	}
	for _, c := range cases {
		if got := mdSafe(c.in); got != c.want {
			t.Errorf("mdSafe(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestPrintTo_Markdown_pipeInTitle(t *testing.T) {
	pub, _ := time.Parse("2006-01-02", "2025-01-15")
	videos := []youtube.VideoStatistics{
		{Title: "Istio Day: Zero-Downtime | Migration", URL: "https://youtu.be/x", PublishedAt: pub},
	}
	var buf bytes.Buffer
	printTo(&buf, videos, "markdown", false, false)

	out := buf.String()
	if strings.Contains(out, "Zero-Downtime | Migration") {
		t.Error("unescaped pipe found in markdown output — table will break")
	}
	if !strings.Contains(out, `Zero-Downtime \| Migration`) {
		t.Error("escaped pipe not found in markdown output")
	}
}

func TestPrintTo_CSV_withScore(t *testing.T) {
	vs := sampleVideos()
	vs[0].Score = 0.987654
	var buf bytes.Buffer
	printTo(&buf, vs, "csv", true, false)

	r := csv.NewReader(strings.NewReader(buf.String()))
	records, _ := r.ReadAll()
	if records[0][0] != "Score" {
		t.Errorf("first header should be Score, got %q", records[0][0])
	}
	if records[1][0] != "0.987654" {
		t.Errorf("Score value: got %q want %q", records[1][0], "0.987654")
	}
}

// --- -strategy all ---

func TestApplyAllStrategies_columns(t *testing.T) {
	pub, _ := time.Parse("2006-01-02", "2023-06-01")
	vs := []youtube.VideoStatistics{
		{Title: "A", ViewCount: 1000, LikeCount: 50, CommentCount: 10, TotalReaction: 60, TotalInterestingness: 0.06, PublishedAt: pub},
		{Title: "B", ViewCount: 500, LikeCount: 10, CommentCount: 2, TotalReaction: 12, TotalInterestingness: 0.02, PublishedAt: pub},
	}
	youtube.ApplyAllStrategies(vs)

	for _, v := range vs {
		if v.AllScores == nil {
			t.Fatalf("AllScores is nil for %q", v.Title)
		}
		for _, slug := range youtube.StrategyOrder {
			if _, ok := v.AllScores[slug]; !ok {
				t.Errorf("missing strategy score %q for video %q", slug, v.Title)
			}
		}
	}
}

func TestPrintTo_CSV_allScores_headers(t *testing.T) {
	pub, _ := time.Parse("2006-01-02", "2023-06-01")
	vs := []youtube.VideoStatistics{
		{Title: "A", ViewCount: 1000, LikeCount: 50, CommentCount: 10, TotalReaction: 60, TotalInterestingness: 0.06, PublishedAt: pub},
	}
	youtube.ApplyAllStrategies(vs)

	var buf bytes.Buffer
	printTo(&buf, vs, "csv", false, true)

	r := csv.NewReader(strings.NewReader(buf.String()))
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("invalid CSV: %v", err)
	}
	// first two columns must be Title and URL
	if records[0][0] != "Title" {
		t.Errorf("header[0]: got %q want %q", records[0][0], "Title")
	}
	if records[0][1] != "URL" {
		t.Errorf("header[1]: got %q want %q", records[0][1], "URL")
	}
	// next N headers should be Score:<slug> in StrategyOrder
	for i, slug := range youtube.StrategyOrder {
		want := "Score:" + slug
		if records[0][2+i] != want {
			t.Errorf("header[%d]: got %q want %q", 2+i, records[0][2+i], want)
		}
	}
	// total columns = 2 (title+url) + 6 strategy scores + 11 remaining metrics
	wantCols := 2 + len(youtube.StrategyOrder) + 11
	if len(records[0]) != wantCols {
		t.Errorf("expected %d columns, got %d", wantCols, len(records[0]))
	}
}

func TestStrategyFlag_all(t *testing.T) {
	resetFlags()
	os.Args = []string{"yrank", "-p", "PLAYLIST", "-strategy", "all"}
	_, _, _, _, _, strategy, _, _, _, _, _, _, _, _, _ := cliParameters()
	if strategy != "all" {
		t.Errorf("expected strategy=all, got %q", strategy)
	}
}

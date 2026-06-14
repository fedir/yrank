package youtube

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"testing"
)

func TestChunkRefs(t *testing.T) {
	mk := func(n int) []videoRef {
		r := make([]videoRef, n)
		for i := range r {
			r[i] = videoRef{ID: strconv.Itoa(i)}
		}
		return r
	}
	cases := []struct {
		n, wantChunks, wantLast int
	}{
		{0, 0, 0},
		{1, 1, 1},
		{50, 1, 50},
		{51, 2, 1},
		{120, 3, 20},
	}
	for _, c := range cases {
		chunks := chunkRefs(mk(c.n), maxIDsPerBatch)
		if len(chunks) != c.wantChunks {
			t.Errorf("n=%d: got %d chunks, want %d", c.n, len(chunks), c.wantChunks)
			continue
		}
		total := 0
		for _, ch := range chunks {
			if len(ch) > maxIDsPerBatch {
				t.Errorf("n=%d: chunk larger than %d", c.n, maxIDsPerBatch)
			}
			total += len(ch)
		}
		if total != c.n {
			t.Errorf("n=%d: chunks hold %d refs total, want %d", c.n, total, c.n)
		}
		if c.wantChunks > 0 && len(chunks[len(chunks)-1]) != c.wantLast {
			t.Errorf("n=%d: last chunk size %d, want %d", c.n, len(chunks[len(chunks)-1]), c.wantLast)
		}
	}
}

// collectStats must attach each batched videos.list result to the correct ref
// (title/publishedAt) by ID — the core correctness property of batching.
func TestCollectStats_mapsByID(t *testing.T) {
	defer withMockClient(t)()

	refs := []videoRef{
		{ID: "NCU_Sebq6Tw", Title: "first", PublishedAt: "2020-01-01T00:00:00Z"},
		{ID: "K4dEFsIxZoM", Title: "second", PublishedAt: "2020-01-02T00:00:00Z"},
	}
	vs := collectStats(refs, "dummy-key", false)
	if len(vs) != 2 {
		t.Fatalf("expected 2 videos, got %d", len(vs))
	}

	byTitle := make(map[string]VideoStatistics, len(vs))
	for _, v := range vs {
		byTitle[v.Title] = v
	}
	// Per-video values from testdata/video_stats.json — must NOT be uniform.
	if got := byTitle["first"].ViewCount; got != 619067 {
		t.Errorf("first ViewCount = %d, want 619067", got)
	}
	if got := byTitle["first"].Duration; got != 725 { // PT12M5S
		t.Errorf("first Duration = %d, want 725", got)
	}
	if got := byTitle["second"].ViewCount; got != 709073 {
		t.Errorf("second ViewCount = %d, want 709073", got)
	}
	if got := byTitle["second"].Duration; got != 45 { // PT45S
		t.Errorf("second Duration = %d, want 45", got)
	}
}

func Test_fetchVideo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(loadRespFromFile("./responses/video.json"))
	}))
	defer ts.Close()

	v := fetchVideo(ts.URL)
	if v.Items[0].Statistics.ViewCount != "222" {
		t.Errorf("unexpected view count: %s", v.Items[0].Statistics.ViewCount)
	}
	if v.Items[0].Statistics.LikeCount != "8" {
		t.Errorf("unexpected like count: %s", v.Items[0].Statistics.LikeCount)
	}
	if v.Items[0].ContentDetails.Duration != "PT4M13S" {
		t.Errorf("unexpected duration: %s", v.Items[0].ContentDetails.Duration)
	}
}

func TestParseISO8601Duration(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{"PT1H2M10S", 3730},
		{"PT45S", 45},
		{"PT12M", 720},
		{"PT8M30S", 510},
		{"P1DT2H", 93600},
		{"P0D", 0},
		{"", 0},
		{"garbage", 0},
		{"PT", 0},
	}
	for _, tc := range tests {
		if got := parseISO8601Duration(tc.in); got != tc.want {
			t.Errorf("parseISO8601Duration(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestIsAnomalousStats(t *testing.T) {
	tests := []struct {
		name         string
		views, likes int
		want         bool
	}{
		{"vsauce live stream (1 view, 112 likes)", 1, 112, true},
		{"zero views with likes", 0, 5, true},
		{"zero views, zero likes", 0, 0, true},
		{"normal video", 619067, 43990, false},
		{"boundary likes == views", 100, 100, false},
		{"likes one over views", 100, 101, true},
	}
	for _, tc := range tests {
		if got := isAnomalousStats(tc.views, tc.likes); got != tc.want {
			t.Errorf("%s: isAnomalousStats(%d, %d) = %v, want %v", tc.name, tc.views, tc.likes, got, tc.want)
		}
	}
}

// An unaired/live-stream placeholder (1 view, 112 likes) must be dropped end-to-end,
// not just by the helper. The anomalous fixture dir has a single such video.
func TestPlaylistStatistics_dropsAnomalous(t *testing.T) {
	orig := httpClient
	SetHTTPClient(NewMockClient(filepath.Join(testDataDir(), "anomalous")))
	defer func() { httpClient = orig }()

	vs := PlaylistStatistics("PLanomalous", "dummy-key", "", false)
	if len(vs) != 0 {
		t.Errorf("expected 0 videos (anomalous row dropped), got %d", len(vs))
	}
}

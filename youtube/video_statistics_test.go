package youtube

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

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

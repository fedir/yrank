package youtube

import (
	"net/http"
	"net/http/httptest"
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

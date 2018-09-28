package youtube

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_video(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f := loadRespFromFile("./responses/video.json")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(f)
	}))
	defer ts.Close()

	testURL := ts.URL
	t.Logf("%s\n", testURL)
	video := video(testURL)
	if video.Items[0].Statistics.ViewCount != "222" {
		t.Errorf("Wrong number of views count")
	}
	if video.Items[0].Statistics.LikeCount != "8" {
		t.Errorf("Wrong number of likes count")
	}

}

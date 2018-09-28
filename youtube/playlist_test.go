package youtube

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_playlist(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f := loadRespFromFile("./responses/playlistItems.json")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(f)
	}))
	defer ts.Close()

	testURL := ts.URL
	t.Logf("%s\n", testURL)
	playlist := playlist(testURL)
	if len(playlist.Items) != 5 {
		t.Errorf("Wrong number of playlists found")
	}
	if playlist.Items[0].Snippet.Title != "GopherCon UK 2018: Aditya Mukerjee - You Might Be a Go Contributor Already and Not Know It" {
		t.Errorf("Wrong first item title")
	}

}

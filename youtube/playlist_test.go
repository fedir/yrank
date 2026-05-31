package youtube

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_fetchPlaylist(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(loadRespFromFile("./responses/playlistItems.json"))
	}))
	defer ts.Close()

	pl := fetchPlaylist(ts.URL)
	if len(pl.Items) != 5 {
		t.Errorf("expected 5 items, got %d", len(pl.Items))
	}
	want := "GopherCon UK 2018: Aditya Mukerjee - You Might Be a Go Contributor Already and Not Know It"
	if pl.Items[0].Snippet.Title != want {
		t.Errorf("unexpected first item title: %s", pl.Items[0].Snippet.Title)
	}
}

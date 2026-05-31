package youtube

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	_, _, err := httpRequest(ts.URL)
	if err != nil {
		t.Errorf("httpRequest() returned an error: %s", err)
	}
}

func Test_readJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(loadRespFromFile("./responses/playlistItems.json"))
	}))
	defer ts.Close()

	body, _, err := httpRequest(ts.URL)
	if err != nil {
		t.Fatalf("httpRequest() error: %v", err)
	}
	var pl Playlist
	if err := json.Unmarshal(body, &pl); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}
	if len(pl.Items) != 5 {
		t.Errorf("expected 5 items, got %d", len(pl.Items))
	}
}

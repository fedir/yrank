package youtube

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func bytesToString(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, ",")
}

func TestOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	testURL := ts.URL
	t.Logf("%s\n", testURL)
	_, _, err := httpRequest(testURL)
	if err != nil {
		t.Errorf("httpRequest() returned an error: %s", err)
	}
}

func Test_readJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f := loadRespFromFile("./responses/playlistItems.json")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(f)
	}))
	defer ts.Close()
	testURL := ts.URL
	t.Logf("%s\n", testURL)
	resp, _, err := httpRequest(testURL)
	if err != nil {
		panic(err)
	}
	jsonResponse, err := readResp(resp)
	if err != nil {
		panic(err)
	}
	playlist := Playlist{}
	json.Unmarshal(jsonResponse, &playlist)
	if len(playlist.Items) != 5 {
		t.Errorf("Wrong number of playlists found")
	}
}

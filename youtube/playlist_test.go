package youtube

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// cyclicTokenClient serves playlistItems pages whose nextPageToken always points
// back to the same page (the YouTube quirk that froze @francetv). It hard-caps the
// number of calls so a pagination regression fails the test instead of hanging it.
type cyclicTokenClient struct {
	calls int
	cap   int
}

func (c *cyclicTokenClient) Do(req *http.Request) (*http.Response, error) {
	c.calls++
	next := `"LOOP"`
	if c.calls > c.cap {
		next = `""` // force termination so a broken guard fails on call count, not by hanging
	}
	body := `{"items":[{"contentDetails":{"videoId":"v` +
		string(rune('0'+c.calls%10)) + `","videoPublishedAt":"2020-01-01T00:00:00Z"},"snippet":{"title":"t"}}],"nextPageToken":` + next + `}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Header:     make(http.Header),
	}, nil
}

// Test_playlistRefs_cyclicToken guards against the infinite-pagination freeze:
// a playlist returning a self-referential nextPageToken must stop quickly.
func Test_playlistRefs_cyclicToken(t *testing.T) {
	prev := httpClient
	c := &cyclicTokenClient{cap: 100}
	SetHTTPClient(c)
	defer SetHTTPClient(prev)

	refs := playlistRefs("PLcyclic", "key", "", false)
	if c.calls > 3 {
		t.Fatalf("expected pagination to stop after the token repeats (≤3 calls), got %d", c.calls)
	}
	if len(refs) == 0 {
		t.Fatalf("expected at least one ref before stopping")
	}
}

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

package youtube

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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

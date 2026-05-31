package youtube

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_fetchChannel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(loadRespFromFile("./responses/channel.json"))
	}))
	defer ts.Close()

	ch := fetchChannel(ts.URL)
	if ch.Items[0].PlaylistID != "PLDWZ5uzn69ewsMyuGjVsAnpQIjyud1Cv9" {
		t.Errorf("unexpected first playlist ID: %s", ch.Items[0].PlaylistID)
	}
}

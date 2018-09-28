package youtube

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_channel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f := loadRespFromFile("./responses/channel.json")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(f)
	}))
	defer ts.Close()

	testURL := ts.URL
	t.Logf("%s\n", testURL)
	channel := channel(testURL)
	if channel.Items[0].PlaylistID != "PLDWZ5uzn69ewsMyuGjVsAnpQIjyud1Cv9" {
		t.Errorf("Wrong first item title")
	}

}

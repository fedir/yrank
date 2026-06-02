package youtube

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// MockTransport serves pre-recorded API responses from testdata/ files.
// URL routing rules:
//   - playlistItems → testdata/playlist_page<N>.json  (N from pageToken, default 1)
//   - videos        → testdata/video_stats.json
type MockTransport struct {
	// DataDir is the directory containing the JSON fixture files.
	// Defaults to the testdata/ dir relative to this source file.
	DataDir string
}

func (m *MockTransport) dataDir() string {
	if m.DataDir != "" {
		return m.DataDir
	}
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "testdata")
}

func (m *MockTransport) Do(req *http.Request) (*http.Response, error) {
	fixture := m.fixtureFor(req.URL)
	data, err := os.ReadFile(fixture)
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(data)),
		Header:     make(http.Header),
	}, nil
}

func (m *MockTransport) fixtureFor(u *url.URL) string {
	dir := m.dataDir()
	path := u.Path
	switch {
	case strings.Contains(path, "playlistItems"):
		page := u.Query().Get("pageToken")
		if page == "" {
			page = "1"
		}
		return filepath.Join(dir, "playlist_page"+page+".json")
	case strings.Contains(path, "/videos"):
		return filepath.Join(dir, "video_stats.json")
	default:
		return filepath.Join(dir, "unknown.json")
	}
}

// NewMockClient returns an HTTPClient backed by local fixture files.
func NewMockClient(dataDir string) HTTPClient {
	return &MockTransport{DataDir: dataDir}
}

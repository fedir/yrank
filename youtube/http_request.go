package youtube

import (
	"io"
	"log"
	"net/http"
	"time"
)

// HTTPClient is the interface used for all API calls. Override in tests or
// local-test mode by calling SetHTTPClient before any fetch.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var httpClient HTTPClient = &http.Client{Timeout: 10 * time.Second}

// SetHTTPClient replaces the package-level client (used by tests and -local-test mode).
func SetHTTPClient(c HTTPClient) {
	httpClient = c
}

func httpRequest(url string) ([]byte, int, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal("cannot prepare HTTP request", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal("cannot process HTTP request", err)
	}
	defer resp.Body.Close()

	statusProcessing(resp.StatusCode, url)

	body, err := io.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

func statusProcessing(statusCode int, url string) {
	switch statusCode {
	case http.StatusForbidden:
		log.Fatalf("rate limit exceeded, please try again later")
	case http.StatusAccepted:
		log.Printf("server needs time to prepare the request")
	case http.StatusOK:
		return
	default:
		log.Fatalf("unexpected status %d for URL %s", statusCode, url)
	}
}

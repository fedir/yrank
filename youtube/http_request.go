package youtube

import (
	"io"
	"log"
	"net/http"
	"time"
)

func httpRequest(url string) ([]byte, int, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal("cannot prepare HTTP request", err)
	}

	resp, err := client.Do(req)
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

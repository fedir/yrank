package youtube

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/url"
	"sync"
)

// SearchStatistics searches YouTube for a word or phrase and returns video
// statistics for the matching videos, ranked among themselves afterwards by the
// caller's chosen sort/strategy.
//
// It paginates the search.list endpoint until it has collected at least
// maxResults videos (or results run out). With maxResults <= 0 it fetches a
// single page (up to 50 results). Each page costs 100 YouTube quota units.
func SearchStatistics(query string, apiKey string, maxResults int, debug bool) []VideoStatistics {
	var stats []VideoStatistics
	token := ""

	for {
		reqURL := "https://www.googleapis.com/youtube/v3/search?part=snippet&type=video&maxResults=50&q=" +
			url.QueryEscape(query) + "&key=" + apiKey
		if token != "" {
			reqURL += "&pageToken=" + token
		}
		if debug {
			fmt.Printf("Search URL: %s\n", reqURL)
		}

		sr := fetchSearch(reqURL)

		dataChan := make(chan VideoStatistics, len(sr.Items))
		var wg sync.WaitGroup
		wg.Add(len(sr.Items))
		for _, item := range sr.Items {
			// search.list HTML-escapes titles (&amp;, &#39;), so unescape here.
			go videoStatistics(item.ID.VideoID, html.UnescapeString(item.Snippet.Title), item.Snippet.PublishedAt, apiKey, dataChan, &wg, debug)
		}
		wg.Wait()
		close(dataChan)

		for vs := range dataChan {
			if vs.Title != "" {
				stats = append(stats, vs)
			}
		}

		if sr.NextPageToken == "" {
			break
		}
		if maxResults <= 0 || len(stats) >= maxResults {
			break
		}
		token = sr.NextPageToken
	}

	return stats
}

func fetchSearch(url string) Search {
	body, _, err := httpRequest(url)
	if err != nil {
		log.Fatalf("search request failed: %v", err)
	}
	var sr Search
	if err := json.Unmarshal(body, &sr); err != nil {
		log.Fatalf("search JSON decode failed: %v", err)
	}
	return sr
}

package youtube

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// PlaylistStatistics returns video statistics for every video in a playlist,
// following pagination automatically.
func PlaylistStatistics(playlistKey string, apiKey string, pageToken string, debug bool) []VideoStatistics {
	url := "https://www.googleapis.com/youtube/v3/playlistItems?playlistId=" + playlistKey + "&maxResults=50&part=snippet%2CcontentDetails&key=" + apiKey
	if pageToken != "" {
		url += "&pageToken=" + pageToken
	}
	if debug {
		fmt.Printf("Playlist URL: %s\n", url)
	}

	pl := fetchPlaylist(url)

	dataChan := make(chan VideoStatistics, len(pl.Items))
	var wg sync.WaitGroup
	wg.Add(len(pl.Items))
	for _, item := range pl.Items {
		go videoStatistics(item.ContentDetails.VideoID, item.Snippet.Title, item.ContentDetails.VideoPublishedAt, apiKey, dataChan, &wg, debug)
	}
	wg.Wait()
	close(dataChan)

	var stats []VideoStatistics
	for vs := range dataChan {
		if vs.Title != "" {
			stats = append(stats, vs)
		}
	}

	if pl.NextPageToken != "" {
		stats = append(stats, PlaylistStatistics(playlistKey, apiKey, pl.NextPageToken, debug)...)
	}
	return stats
}

func fetchPlaylist(url string) Playlist {
	body, _, err := httpRequest(url)
	if err != nil {
		log.Fatalf("playlist request failed: %v", err)
	}
	var pl Playlist
	if err := json.Unmarshal(body, &pl); err != nil {
		log.Fatalf("playlist JSON decode failed: %v", err)
	}
	return pl
}

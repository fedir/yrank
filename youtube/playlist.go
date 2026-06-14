package youtube

import (
	"encoding/json"
	"fmt"
	"log"
)

// PlaylistStatistics returns video statistics for every video in a playlist.
func PlaylistStatistics(playlistKey string, apiKey string, pageToken string, debug bool) []VideoStatistics {
	return collectStats(playlistRefs(playlistKey, apiKey, pageToken, debug), apiKey, debug)
}

// playlistRefs paginates a playlist's items and returns the listing metadata
// (id, title, publishedAt) for every video, without fetching any stats.
func playlistRefs(playlistKey string, apiKey string, pageToken string, debug bool) []videoRef {
	var refs []videoRef
	token := pageToken

	for {
		url := "https://www.googleapis.com/youtube/v3/playlistItems?playlistId=" + playlistKey + "&maxResults=50&part=snippet%2CcontentDetails&key=" + apiKey
		if token != "" {
			url += "&pageToken=" + token
		}
		if debug {
			fmt.Printf("Playlist URL: %s\n", url)
		}

		pl := fetchPlaylist(url)
		for _, item := range pl.Items {
			refs = append(refs, videoRef{
				ID:          item.ContentDetails.VideoID,
				Title:       item.Snippet.Title,
				PublishedAt: item.ContentDetails.VideoPublishedAt,
			})
		}

		if pl.NextPageToken == "" {
			break
		}
		token = pl.NextPageToken
	}

	return refs
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

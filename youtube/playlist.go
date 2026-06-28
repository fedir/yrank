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

	// Some playlists return a nextPageToken that points back to a page we've
	// already fetched (a self-referential or cyclic token), which would loop
	// forever. Track the tokens we've requested and stop if one repeats.
	seenTokens := map[string]bool{}

	for {
		url := "https://www.googleapis.com/youtube/v3/playlistItems?playlistId=" + playlistKey + "&maxResults=50&part=snippet%2CcontentDetails&key=" + apiKey
		if token != "" {
			url += "&pageToken=" + token
		}
		if debug {
			fmt.Printf("Playlist URL: %s\n", url)
		}
		seenTokens[token] = true

		pl := fetchPlaylist(url)
		for _, item := range pl.Items {
			refs = append(refs, videoRef{
				ID:          item.ContentDetails.VideoID,
				Title:       item.Snippet.Title,
				PublishedAt: item.ContentDetails.VideoPublishedAt,
			})
		}

		if pl.NextPageToken == "" || seenTokens[pl.NextPageToken] {
			if debug && seenTokens[pl.NextPageToken] && pl.NextPageToken != "" {
				fmt.Printf("Stopping pagination for %s: nextPageToken %q repeats (non-advancing)\n", playlistKey, pl.NextPageToken)
			}
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

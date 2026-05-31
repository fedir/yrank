package youtube

import (
	"encoding/json"
	"fmt"
	"log"
)

// ChannelStatistics returns video statistics for all playlists in a channel.
func ChannelStatistics(cid string, apiKey string, debug bool) []VideoStatistics {
	url := "https://www.googleapis.com/youtube/v3/playlists?channelId=" + cid + "&part=id&maxResults=50&key=" + apiKey
	if debug {
		fmt.Printf("Channel URL: %s\n", url)
	}

	ch := fetchChannel(url)
	var all []VideoStatistics
	for _, pl := range ch.Items {
		if debug {
			fmt.Printf("Getting videos from playlist: %s\n", pl.PlaylistID)
		}
		all = append(all, PlaylistStatistics(pl.PlaylistID, apiKey, "", debug)...)
	}

	// Deduplicate — a video can appear in multiple playlists.
	seen := make(map[string]bool, len(all))
	unique := make([]VideoStatistics, 0, len(all))
	for _, v := range all {
		if !seen[v.Key] {
			seen[v.Key] = true
			unique = append(unique, v)
		}
	}
	return unique
}

// ResolveHandle resolves a YouTube handle (e.g. "@Squeezie") to a channel ID.
func ResolveHandle(handle string, apiKey string) string {
	url := "https://www.googleapis.com/youtube/v3/channels?part=id&forHandle=" + handle + "&key=" + apiKey
	body, _, err := httpRequest(url)
	if err != nil {
		log.Fatalf("handle resolution request failed: %v", err)
	}
	var result ChannelByHandle
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("handle resolution JSON decode failed: %v", err)
	}
	if len(result.Items) == 0 {
		log.Fatalf("no channel found for handle %q", handle)
	}
	return result.Items[0].ID
}

func fetchChannel(url string) Channel {
	body, _, err := httpRequest(url)
	if err != nil {
		log.Fatalf("channel request failed: %v", err)
	}
	var ch Channel
	if err := json.Unmarshal(body, &ch); err != nil {
		log.Fatalf("channel JSON decode failed: %v", err)
	}
	return ch
}

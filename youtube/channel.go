package youtube

import (
	"encoding/json"
	"fmt"
	"log"
)

// uploadsPlaylistID returns the auto-generated uploads playlist ID for a channel.
// YouTube derives it by replacing the "UC" prefix with "UU".
func uploadsPlaylistID(channelID string) string {
	if len(channelID) >= 2 && channelID[:2] == "UC" {
		return "UU" + channelID[2:]
	}
	return channelID
}

// ChannelStatistics returns video statistics for all videos in a channel,
// including videos not assigned to any manual playlist (via the uploads playlist).
func ChannelStatistics(cid string, apiKey string, debug bool) []VideoStatistics {
	// Always start with the uploads playlist — it contains every uploaded video.
	uploadsID := uploadsPlaylistID(cid)
	if debug {
		fmt.Printf("Uploads playlist: %s\n", uploadsID)
	}
	all := PlaylistStatistics(uploadsID, apiKey, "", debug)

	// Also fetch manual playlists to pick up any videos that may differ
	// (rare, but keeps behaviour consistent with prior versions).
	url := "https://www.googleapis.com/youtube/v3/playlists?channelId=" + cid + "&part=id&maxResults=50&key=" + apiKey
	if debug {
		fmt.Printf("Channel URL: %s\n", url)
	}
	ch := fetchChannel(url)
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

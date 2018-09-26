package youtube

import (
	"encoding/json"
	"fmt"
)

// ChannelStatistics returns videos statistics of a Youtube channel
func ChannelStatistics(cid string, apiKey string, debug bool) []VideoStatistics {

	var channelVideos = []VideoStatistics{}

	url := "https://www.googleapis.com/youtube/v3/playlists?channelId=" + cid + "&part=id&maxResults=50&key=" + apiKey

	resp, _, err := httpRequest(url)
	if err != nil {
		panic(err)
	}

	jsonResponse, _ := readResp(resp)
	channel := Channel{}
	json.Unmarshal(jsonResponse, &channel)

	for _, pl := range channel.Items {
		if debug {
			fmt.Printf("Getting videos from playlist: %s\n", pl.PlaylistID)
		}
		playlistVideos := PlaylistStatistics(pl.PlaylistID, apiKey, "", debug)
		channelVideos = append(channelVideos, playlistVideos...)
	}

	// Single video could be in multiple playlists, so we should make the array unique
	// (the best will be to do it even before to collect the videos details, to avoid additional usage of API quota)
	var uniqueVideos []VideoStatistics
	uniqueVideosKeys := make(map[string]bool)
	for _, video := range channelVideos {
		if _, ok := uniqueVideosKeys[video.Key]; !ok {
			uniqueVideosKeys[video.Key] = true
			uniqueVideos = append(uniqueVideos, video)
		}
	}

	return uniqueVideos
}

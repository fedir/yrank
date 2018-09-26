package youtube

import (
	"encoding/json"
	"fmt"
	"sync"
)

// PlaylistStatistics returns videos statistics of a Youtube playlist
func PlaylistStatistics(playlistKey string, apiKey string, pageToken string, debug bool) []VideoStatistics {
	url := "https://www.googleapis.com/youtube/v3/playlistItems?playlistId=" + playlistKey + "&part=snippet%2CcontentDetails&key=" + apiKey
	if pageToken != "" {
		url = url + "&pageToken=" + pageToken
	}
	if debug {
		fmt.Printf("Playlist URL: %s\n", url)
	}

	resp, _, err := httpRequest(url)
	if err != nil {
		panic(err)
	}
	jsonResponse, _, _ := readResp(resp)
	playlist := Playlist{}
	json.Unmarshal(jsonResponse, &playlist)

	var playlistStatistic = []VideoStatistics{}
	var wg sync.WaitGroup
	wg.Add(len(playlist.Items))

	dataChan := make(chan VideoStatistics, len(playlist.Items))
	for _, video := range playlist.Items {
		go videoStatistics(video.ContentDetails.VideoID, video.Snippet.Title, apiKey, dataChan, &wg, debug)
	}
	for range playlist.Items {
		vs := <-dataChan
		if vs.Title != "" {
			playlistStatistic = append(playlistStatistic, vs)
		}
	}
	wg.Wait()

	if playlist.NextPageToken != "" {
		nextPagePlaylistStatistics := PlaylistStatistics(playlistKey, apiKey, playlist.NextPageToken, debug)
		playlistStatistic = append(playlistStatistic, nextPagePlaylistStatistics...)
	}

	return playlistStatistic
}

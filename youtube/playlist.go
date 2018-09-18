package youtube

import (
	"encoding/json"
	"fmt"

	"github.com/fedir/ghstat/httpcache"
)

// RankedPlaylist returns ranked playlist
func RankedPlaylist(playlistKey string, apiKey string) string {
	url := "https://www.googleapis.com/youtube/v3/playlistItems?playlistId=" + playlistKey + "&maxResults=50&part=snippet%2CcontentDetails&key=" + apiKey
	fmt.Printf("Playlist URL: %s\n", url)

	resp, _, err := httpRequest(url)
	if err != nil {
		panic(err)
	}
	jsonResponse, _, _ := httpcache.ReadResp(resp)
	playlist := Playlist{}
	json.Unmarshal(jsonResponse, &playlist)
	for _, item := range playlist.Items {
		fmt.Println(item.ContentDetails.VideoID)
	}
	return "RankedPlaylist"
}

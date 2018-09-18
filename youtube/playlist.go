package youtube

import (
	"encoding/json"
	"fmt"
	"sync"

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

	var playlistStatistic = []VideoStatistics{}
	var wg sync.WaitGroup
	wg.Add(len(playlist.Items))

	dataChan := make(chan VideoStatistics, len(playlist.Items))
	for _, video := range playlist.Items {
		go videoStatistics(video.ContentDetails.VideoID, video.Snippet.Title, apiKey, dataChan, &wg)
	}
	for range playlist.Items {
		playlistStatistic = append(playlistStatistic, <-dataChan)
	}
	wg.Wait()

	fmt.Printf("%v#", playlistStatistic)

	//rateAndPrintGreetings(playlistStatistic)
	//writeCSVStatistics(playlistStatistic, csvFilePath)

	return "RankedPlaylist"
}

func videoStatistics(vid string, title string, apiKey string, dataChan chan VideoStatistics, wg *sync.WaitGroup) {

	defer wg.Done()

	vs := new(VideoStatistics)

	url := "https://www.googleapis.com/youtube/v3/videos?part=statistics&id=" + vid + "&key=" + apiKey

	resp, _, err := httpRequest(url)
	if err != nil {
		panic(err)
	}
	jsonResponse, _, _ := httpcache.ReadResp(resp)
	video := Video{}
	json.Unmarshal(jsonResponse, &video)

	for _, item := range video.Items {
		vs.Key = item.ID
		vs.URL = "https://www.youtube.com/watch?v=" + item.ID
		vs.ViewCount = item.Statistics.ViewCount
		vs.LikeCount = item.Statistics.LikeCount
		vs.DislikeCount = item.Statistics.DislikeCount
		vs.CommentCount = item.Statistics.CommentCount
		vs.Title = title

		// Coefficient of likes / dislikes
		// Proportion of views / likes-dislikes
		// Total reaction coefficient
		// Total likes
	}

	dataChan <- *vs
}

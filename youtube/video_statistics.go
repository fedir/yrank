package youtube

import (
	"encoding/json"
	"log"
	"sync"
)

func videoStatistics(vid string, title string, apiKey string, dataChan chan VideoStatistics, wg *sync.WaitGroup) {

	defer wg.Done()

	vs := new(VideoStatistics)

	url := "https://www.googleapis.com/youtube/v3/videos?part=statistics&id=" + vid + "&key=" + apiKey

	resp, _, err := httpRequest(url)
	if err != nil {
		panic(err)
	}
	jsonResponse, _, err := readResp(resp)
	if err != nil {
		log.Println(err)
	}
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

		// Total views
		// Total likes
		// Total dislikes
		// Total comments
		// Global rating coefficient of likes / dislikes
		// Relative Most rated (Proportion of views / likes-dislikes)
		// Total reaction coefficient
	}

	dataChan <- *vs
}

package youtube

import (
	"encoding/json"
	"log"
	"strconv"
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
		vs.ViewCount, _ = strconv.Atoi(item.Statistics.ViewCount)
		vs.LikeCount, _ = strconv.Atoi(item.Statistics.LikeCount)
		vs.DislikeCount, _ = strconv.Atoi(item.Statistics.DislikeCount)
		vs.CommentCount, _ = strconv.Atoi(item.Statistics.CommentCount)
		vs.Title = title
		vs.PositiveInterestingness = float64(vs.LikeCount-vs.DislikeCount) / float64(vs.ViewCount)
		vs.TotalReaction = float64(vs.LikeCount+vs.DislikeCount+vs.CommentCount) / float64(vs.ViewCount)
		vs.GlobalBuzz = vs.ViewCount + vs.CommentCount
	}

	dataChan <- *vs
}

package youtube

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

func videoStatistics(vid string, title string, publishedAt string, apiKey string, dataChan chan VideoStatistics, wg *sync.WaitGroup, debug bool) {

	defer wg.Done()

	vs := new(VideoStatistics)

	url := "https://www.googleapis.com/youtube/v3/videos?part=statistics&id=" + vid + "&key=" + apiKey

	if debug {
		fmt.Printf("Video URL: %s\n", url)
	}

	video := video(url)

	for _, item := range video.Items {
		vs.Key = item.ID
		vs.URL = "https://www.youtube.com/watch?v=" + item.ID
		var err error
		vs.PublishedAt, err = time.Parse(time.RFC3339, publishedAt)
		if err != nil {
			panic(err)
		}
		vs.ViewCount, _ = strconv.Atoi(item.Statistics.ViewCount)
		vs.LikeCount, _ = strconv.Atoi(item.Statistics.LikeCount)
		vs.DislikeCount, _ = strconv.Atoi(item.Statistics.DislikeCount)
		vs.CommentCount, _ = strconv.Atoi(item.Statistics.CommentCount)
		vs.TotalReaction = vs.LikeCount + vs.DislikeCount + vs.CommentCount
		vs.Title = title
		vs.PositiveInterestingness = float64(vs.LikeCount-vs.DislikeCount) / float64(vs.ViewCount)
		vs.PositiveNegativeCoefficient = float64(vs.LikeCount) / float64(1+vs.DislikeCount)
		vs.TotalInterestingness = float64(vs.LikeCount+vs.DislikeCount+vs.CommentCount) / float64(vs.ViewCount)
		vs.GlobalBuzzIndex = vs.ViewCount * (vs.LikeCount + vs.DislikeCount + vs.CommentCount)
	}

	dataChan <- *vs
}

func video(url string) Video {
	resp, _, err := httpRequest(url)
	if err != nil {
		panic(err)
	}
	jsonResponse, err := readResp(resp)
	if err != nil {
		log.Println(err)
	}
	video := Video{}
	json.Unmarshal(jsonResponse, &video)
	return video
}

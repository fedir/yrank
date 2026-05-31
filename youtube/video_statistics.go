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

	url := "https://www.googleapis.com/youtube/v3/videos?part=statistics&id=" + vid + "&key=" + apiKey
	if debug {
		fmt.Printf("Video URL: %s\n", url)
	}

	v := fetchVideo(url)
	if len(v.Items) == 0 {
		dataChan <- VideoStatistics{}
		return
	}

	item := v.Items[0]
	pub, err := time.Parse(time.RFC3339, publishedAt)
	if err != nil {
		log.Printf("invalid publishedAt %q for video %s: %v", publishedAt, vid, err)
	}

	views, _ := strconv.Atoi(item.Statistics.ViewCount)
	likes, _ := strconv.Atoi(item.Statistics.LikeCount)
	dislikes, _ := strconv.Atoi(item.Statistics.DislikeCount)
	comments, _ := strconv.Atoi(item.Statistics.CommentCount)

	dataChan <- VideoStatistics{
		Key:                         item.ID,
		URL:                         "https://www.youtube.com/watch?v=" + item.ID,
		Title:                       title,
		PublishedAt:                 pub,
		ViewCount:                   views,
		LikeCount:                   likes,
		DislikeCount:                dislikes,
		CommentCount:                comments,
		TotalReaction:               likes + dislikes + comments,
		PositiveInterestingness:     float64(likes-dislikes) / float64(views),
		PositiveNegativeCoefficient: float64(likes) / float64(1+dislikes),
		TotalInterestingness:        float64(likes+dislikes+comments) / float64(views),
		GlobalBuzzIndex:             views * (likes + dislikes + comments),
	}
}

func fetchVideo(url string) Video {
	body, _, err := httpRequest(url)
	if err != nil {
		log.Fatalf("video request failed: %v", err)
	}
	var v Video
	if err := json.Unmarshal(body, &v); err != nil {
		log.Fatalf("video JSON decode failed: %v", err)
	}
	return v
}

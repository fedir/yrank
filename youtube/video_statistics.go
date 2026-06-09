package youtube

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"sync"
	"time"
)

func videoStatistics(vid string, title string, publishedAt string, apiKey string, dataChan chan VideoStatistics, wg *sync.WaitGroup, debug bool) {
	defer wg.Done()

	url := "https://www.googleapis.com/youtube/v3/videos?part=statistics,contentDetails&id=" + vid + "&key=" + apiKey
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

	if isAnomalousStats(views, likes) {
		if debug {
			fmt.Printf("Skipping anomalous video %s: views=%d likes=%d (likely unaired/live stream)\n", vid, views, likes)
		}
		dataChan <- VideoStatistics{}
		return
	}

	dataChan <- VideoStatistics{
		Key:                         item.ID,
		URL:                         "https://www.youtube.com/watch?v=" + item.ID,
		Title:                       title,
		PublishedAt:                 pub,
		Duration:                    parseISO8601Duration(item.ContentDetails.Duration),
		ViewCount:                   views,
		LikeCount:                   likes,
		DislikeCount:                dislikes,
		CommentCount:                comments,
		TotalReaction:               likes + dislikes + comments,
		PositiveInterestingness:     safeDiv(float64(likes-dislikes), float64(views)),
		PositiveNegativeCoefficient: float64(likes) / float64(1+dislikes),
		TotalInterestingness:        safeDiv(float64(likes+dislikes+comments), float64(views)),
		GlobalBuzzIndex:             views * (likes + dislikes + comments),
	}
}

// isAnomalousStats reports whether a video's stats are physically impossible
// (e.g. unaired/live-stream placeholders that report ~0 views but carry likes).
// A like requires a view, so likes > views — or no views at all — means the data
// is bogus and the video must be dropped to avoid distorting view-normalised metrics.
func isAnomalousStats(views, likes int) bool {
	return views <= 0 || likes > views
}

var iso8601Duration = regexp.MustCompile(`^P(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?)?$`)

// parseISO8601Duration converts a YouTube ISO-8601 duration ("PT1H2M10S") to
// seconds. Returns 0 for empty or unparseable values (e.g. live placeholders).
func parseISO8601Duration(s string) int {
	m := iso8601Duration.FindStringSubmatch(s)
	if m == nil {
		return 0
	}
	atoi := func(v string) int {
		n, _ := strconv.Atoi(v)
		return n
	}
	return atoi(m[1])*86400 + atoi(m[2])*3600 + atoi(m[3])*60 + atoi(m[4])
}

func safeDiv(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
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

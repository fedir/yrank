package youtube

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// videoRef is the listing metadata (from playlistItems/search) for one video,
// carried alongside the ID so it can be re-attached after the batched stats fetch.
type videoRef struct {
	ID          string
	Title       string
	PublishedAt string
}

// maxIDsPerBatch is the videos.list cap: up to 50 IDs per call, all for 1 quota unit.
const maxIDsPerBatch = 50

// statsWorkers bounds how many batch requests run concurrently, keeping throughput
// up while capping burst so we don't trip the per-100s rate limit on large channels.
const statsWorkers = 6

// collectStats fetches statistics for every ref using batched videos.list calls
// (≤50 IDs each), then maps each result back to its listing title/publishedAt.
// Anomalous rows (impossible engagement) are dropped. Order is not preserved; the
// caller sorts afterwards.
func collectStats(refs []videoRef, apiKey string, debug bool) []VideoStatistics {
	chunks := chunkRefs(refs, maxIDsPerBatch)
	if len(chunks) == 0 {
		return nil
	}

	results := make([][]VideoStatistics, len(chunks))
	jobs := make(chan int)
	var wg sync.WaitGroup

	workers := statsWorkers
	if len(chunks) < workers {
		workers = len(chunks)
	}
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				results[i] = fetchStatsBatch(chunks[i], apiKey, debug)
			}
		}()
	}
	for i := range chunks {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	var stats []VideoStatistics
	for _, chunk := range results {
		stats = append(stats, chunk...)
	}
	return stats
}

// fetchStatsBatch fetches one chunk of refs (≤50) in a single videos.list call.
func fetchStatsBatch(refs []videoRef, apiKey string, debug bool) []VideoStatistics {
	byID := make(map[string]videoRef, len(refs))
	ids := make([]string, 0, len(refs))
	for _, r := range refs {
		byID[r.ID] = r
		ids = append(ids, r.ID)
	}

	url := "https://www.googleapis.com/youtube/v3/videos?part=statistics,contentDetails&id=" +
		strings.Join(ids, ",") + "&key=" + apiKey
	if debug {
		fmt.Printf("Video URL (%d ids): %s\n", len(ids), url)
	}

	v := fetchVideo(url)
	stats := make([]VideoStatistics, 0, len(v.Items))
	for _, item := range v.Items {
		ref, ok := byID[item.ID]
		if !ok {
			continue
		}
		vs, keep := buildVideoStatistics(item, ref.Title, ref.PublishedAt, debug)
		if keep {
			stats = append(stats, vs)
		}
	}
	return stats
}

// chunkRefs splits refs into groups of at most size.
func chunkRefs(refs []videoRef, size int) [][]videoRef {
	if size <= 0 || len(refs) == 0 {
		return nil
	}
	chunks := make([][]videoRef, 0, (len(refs)+size-1)/size)
	for i := 0; i < len(refs); i += size {
		end := i + size
		if end > len(refs) {
			end = len(refs)
		}
		chunks = append(chunks, refs[i:end])
	}
	return chunks
}

// buildVideoStatistics computes the derived metrics for one video item, attaching
// the listing title/publishedAt. The bool is false when the row must be dropped
// (anomalous/impossible engagement).
func buildVideoStatistics(item VideoItem, title, publishedAt string, debug bool) (VideoStatistics, bool) {
	pub, err := time.Parse(time.RFC3339, publishedAt)
	if err != nil {
		log.Printf("invalid publishedAt %q for video %s: %v", publishedAt, item.ID, err)
	}

	views, _ := strconv.Atoi(item.Statistics.ViewCount)
	likes, _ := strconv.Atoi(item.Statistics.LikeCount)
	dislikes, _ := strconv.Atoi(item.Statistics.DislikeCount)
	comments, _ := strconv.Atoi(item.Statistics.CommentCount)

	if isAnomalousStats(views, likes) {
		if debug {
			fmt.Printf("Skipping anomalous video %s: views=%d likes=%d (likely unaired/live stream)\n", item.ID, views, likes)
		}
		return VideoStatistics{}, false
	}

	return VideoStatistics{
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
	}, true
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

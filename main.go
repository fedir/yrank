package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fedir/yrank/youtube"
)

// version is overridden at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	cid, pid, topSearch, of, sorting, strategy, from, weightsRaw, outFile, m, minLength, maxLength, minViews, d, localTest, inFile := cliParameters()

	// -in mode: filter an existing CSV export locally, no API key required.
	if inFile != "" {
		if err := filterCSVFile(inFile, outFile, minViews, minLength, maxLength); err != nil {
			log.Fatalf("filter file: %v", err)
		}
		return
	}

	c := configuration()
	if localTest {
		youtube.SetHTTPClient(youtube.NewMockClient("testdata"))
	}
	if d {
		fmt.Printf("API key: %s\n", c.apikey)
	}

	rankedVideos := fetchVideos(c.apikey, cid, pid, topSearch, m, d)

	if from != "" {
		fromDate, _ := time.Parse("2006-01-02", from)
		rankedVideos = filterFrom(rankedVideos, fromDate)
	}
	rankedVideos = filterByLength(rankedVideos, minLength, maxLength)
	rankedVideos = filterByViews(rankedVideos, minViews)

	allStrategies := strategy == "all"
	if allStrategies {
		youtube.ApplyAllStrategies(rankedVideos)
	} else if strategy != "" {
		env := envWeights()
		cli := parseWeightsFlag(weightsRaw)
		weights := youtube.ResolveWeights(strategy, env, cli)
		youtube.ApplyStrategy(rankedVideos, strategy, weights)
	} else {
		youtube.SortBy(rankedVideos, sorting)
	}

	if m > 0 && m <= len(rankedVideos) {
		rankedVideos = rankedVideos[:m]
	}
	if outFile != "" {
		if err := printToFile(outFile, rankedVideos, of, strategy != "", allStrategies); err != nil {
			log.Fatalf("write output: %v", err)
		}
	} else {
		print(rankedVideos, of, strategy != "", allStrategies)
	}
}

// fetchVideos resolves the chosen input source (channel, playlist or search)
// to its ranked video statistics.
func fetchVideos(apikey, cid, pid, topSearch string, m int, d bool) []youtube.VideoStatistics {
	switch {
	case cid != "":
		if strings.HasPrefix(cid, "@") {
			cid = youtube.ResolveHandle(cid, apikey)
		}
		if d {
			fmt.Printf("Channel ID: %s\n", cid)
		}
		return youtube.ChannelStatistics(cid, apikey, d)
	case pid != "":
		if d {
			fmt.Printf("Playlist ID: %s\n", pid)
		}
		return youtube.PlaylistStatistics(pid, apikey, "", d)
	case topSearch != "":
		if d {
			fmt.Printf("Search query: %s\n", topSearch)
		}
		return youtube.SearchStatistics(topSearch, apikey, m, d)
	}
	return nil
}

// filterByViews keeps videos with at least min views. A min of 0 is a no-op.
func filterByViews(videos []youtube.VideoStatistics, min int) []youtube.VideoStatistics {
	if min <= 0 {
		return videos
	}
	out := videos[:0]
	for _, v := range videos {
		if v.ViewCount >= min {
			out = append(out, v)
		}
	}
	return out
}

func filterFrom(videos []youtube.VideoStatistics, from time.Time) []youtube.VideoStatistics {
	out := videos[:0]
	for _, v := range videos {
		if !v.PublishedAt.Before(from) {
			out = append(out, v)
		}
	}
	return out
}

// filterByLength keeps videos whose duration (seconds) falls within
// [min, max]. A bound of 0 means "no limit" for that side. Videos with an
// unknown duration (0, e.g. live placeholders) are dropped only when min > 0.
func filterByLength(videos []youtube.VideoStatistics, min, max int) []youtube.VideoStatistics {
	if min <= 0 && max <= 0 {
		return videos
	}
	out := videos[:0]
	for _, v := range videos {
		if min > 0 && v.Duration < min {
			continue
		}
		if max > 0 && v.Duration > max {
			continue
		}
		out = append(out, v)
	}
	return out
}

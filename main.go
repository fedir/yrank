package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fedir/yrank/youtube"
)

func main() {
	c := configuration()
	cid, pid, topSearch, of, sorting, strategy, from, weightsRaw, outFile, m, d, localTest := cliParameters()
	if localTest {
		youtube.SetHTTPClient(youtube.NewMockClient("testdata"))
	}
	if d {
		fmt.Printf("API key: %s\n", c.apikey)
	}

	var rankedVideos []youtube.VideoStatistics
	if cid != "" {
		if strings.HasPrefix(cid, "@") {
			cid = youtube.ResolveHandle(cid, c.apikey)
		}
		if d {
			fmt.Printf("Channel ID: %s\n", cid)
		}
		rankedVideos = youtube.ChannelStatistics(cid, c.apikey, d)
	} else if pid != "" {
		if d {
			fmt.Printf("Playlist ID: %s\n", pid)
		}
		rankedVideos = youtube.PlaylistStatistics(pid, c.apikey, "", d)
	} else if topSearch != "" {
		if d {
			fmt.Printf("Search query: %s\n", topSearch)
		}
		rankedVideos = youtube.SearchStatistics(topSearch, c.apikey, m, d)
	}

	if from != "" {
		fromDate, _ := time.Parse("2006-01-02", from)
		rankedVideos = filterFrom(rankedVideos, fromDate)
	}

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

func filterFrom(videos []youtube.VideoStatistics, from time.Time) []youtube.VideoStatistics {
	out := videos[:0]
	for _, v := range videos {
		if !v.PublishedAt.Before(from) {
			out = append(out, v)
		}
	}
	return out
}

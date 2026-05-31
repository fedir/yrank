package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/fedir/yrank/youtube"
)

func main() {
	c := configuration()
	cid, pid, of, sorting, strategy, from, weightsRaw, m, d := cliParameters()
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
	}

	if from != "" {
		fromDate, _ := time.Parse("2006-01-02", from)
		rankedVideos = filterFrom(rankedVideos, fromDate)
	}

	if strategy != "" {
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
	print(rankedVideos, of, strategy != "")
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

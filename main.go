package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/fedir/yrank/youtube"
)

func main() {

	// Getting configuration from configuration files and from CLI parameters
	c := configuration()
	cid, pid, of, s, m, from, d := cliParameters()
	if d {
		fmt.Printf("API key: %s\n", c.apikey)
	}

	// Statistics retrieve from Youtube
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

	// Filter by date
	if from != "" {
		fromDate, _ := time.Parse("2006-01-02", from)
		filtered := rankedVideos[:0]
		for _, v := range rankedVideos {
			if !v.PublishedAt.Before(fromDate) {
				filtered = append(filtered, v)
			}
		}
		rankedVideos = filtered
	}

	// Sorting
	youtube.SortBy(rankedVideos, s)

	// Limiting number of results
	if m > 0 && m <= len(rankedVideos) {
		rankedVideos = rankedVideos[:m]
	}
	print(rankedVideos, of)
}

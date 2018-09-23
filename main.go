package main

import (
	"fmt"

	"github.com/fedir/yrank/youtube"
)

func main() {

	// Getting configuration from configuration files and from CLI parameters
	c := configuration()
	cid, pid, of, s, d := cliParameters()
	if d {
		fmt.Printf("API key: %s\n", c.apikey)
	}

	// Statistics retrieve from Youtube
	var rankedVideos []youtube.VideoStatistics
	if cid != "" {
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

	// Sorting
	if s == "likes" {
		youtube.SortByLikes(rankedVideos)
	} else if s == "total-interest" {
		youtube.SortByTotalInterestingness(rankedVideos)
	} else if s == "positive-interest" {
		youtube.SortByPositiveInterestingness(rankedVideos)
	}

	print(rankedVideos, of)
}

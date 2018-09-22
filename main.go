package main

import (
	"fmt"

	"github.com/fedir/yrank/youtube"
)

func main() {

	// Getting configuration from configuration files and from CLI parameters
	c := configuration()
	fmt.Printf("API key: %s\n", c.apikey)
	cid, pid, of, s := cliParameters()

	// Statistics retrieve from Youtube
	var rankedVideos []youtube.VideoStatistics
	if cid != "" {
		fmt.Printf("Channel ID: %s\n", cid)
		rankedVideos = youtube.ChannelStatistics(cid, c.apikey)
	} else if pid != "" {
		fmt.Printf("Playlist ID: %s\n", pid)
		rankedVideos = youtube.PlaylistStatistics(pid, c.apikey, "")
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

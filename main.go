package main

import (
	"fmt"

	"github.com/fedir/yrank/youtube"
)

func main() {
	c := configuration()
	fmt.Printf("API key: %s\n", c.apikey)
	pk, of, s := cliParameters()
	fmt.Printf("Playlist key: %s\n", pk)
	ps := youtube.PlaylistStatistics(pk, c.apikey)
	if s == "likes" {
		youtube.SortByLikes(ps)
	} else if s == "total-interest" {
		youtube.SortByTotalInterestingness(ps)
	} else if s == "positive-interest" {
		youtube.SortByPositiveInterestingness(ps)
	}
	print(ps, of)
}

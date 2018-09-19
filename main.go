package main

import (
	"fmt"

	"github.com/fedir/yrank/youtube"
)

func main() {
	c := configuration()
	fmt.Printf("API key: %s\n", c.apikey)
	pk := playlistKey()
	fmt.Printf("Playlist key: %s\n", pk)
	ps := youtube.PlaylistStatistics(pk, c.apikey)
	print(ps)

	//rateStatistics(playlistStatistic)
	//writeCSVStatistics(playlistStatistic, csvFilePath)
}

package main

import (
	"fmt"

	"github.com/fedir/yrank/youtube"
)

func main() {
	c := configuration()
	fmt.Printf("API key: %s\n", c.apikey)
	pk, of := cliParameters()
	fmt.Printf("Playlist key: %s\n", pk)
	ps := youtube.PlaylistStatistics(pk, c.apikey)
	print(ps, of)

	//rateStatistics(playlistStatistic)
	//writeCSVStatistics(playlistStatistic, csvFilePath)
}

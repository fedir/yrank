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
	rankedPlaylist := youtube.RankedPlaylist(pk, c.apikey)
	fmt.Printf("%v#", rankedPlaylist)
}

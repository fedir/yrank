package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func configuration() Configuration {
	godotenv.Load()
	apikey := os.Getenv("YOUTUBE_API_KEY")
	if apikey == "" {
		log.Fatalln("YOUTUBE_API_KEY environment variable is not set")
	}
	return Configuration{apikey: apikey}
}

func cliParameters() (string, string, string, string, int, string, bool) {
	var (
		playlistID = flag.String("p", "", "Youtube playlist ID")
		channelID  = flag.String("c", "", "Youtube channel ID")
		output     = flag.String("o", "table", "Output format {table|markdown}")
		sorting    = flag.String("s", "total-interest", "Sorting {total-interest|positive-interest|global-buzz-index|total-reaction|positive-negative-coefficient|pnc}")
		maxResults = flag.Int("m", 0, "The maximum number of items that should be returned")
		from       = flag.String("from", "", "Only include videos published on or after this date (YYYY-MM-DD)")
		debug      = flag.Bool("d", false, "Debug mode")
	)
	flag.Parse()
	if *playlistID == "" && *channelID == "" {
		log.Fatalln("Playlist ID or channel ID must be defined")
	} else if *playlistID != "" && *channelID != "" {
		log.Fatalln("Playlist ID & channel ID cannot be used together")
	}
	if *output != "table" && *output != "markdown" {
		log.Fatalln("Unknown output format")
	}
	if *sorting != "likes" && *sorting != "total-interest" && *sorting != "positive-interest" && *sorting != "global-buzz-index" && *sorting != "total-reaction" && *sorting != "pnc" && *sorting != "positive-negative-coefficient" {
		log.Fatalln("Unknown sorting column")
	}
	if *from != "" {
		if _, err := time.Parse("2006-01-02", *from); err != nil {
			log.Fatalln("Invalid -from date, expected format YYYY-MM-DD")
		}
	}
	return *channelID, *playlistID, *output, *sorting, *maxResults, *from, *debug
}

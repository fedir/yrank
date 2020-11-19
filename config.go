package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

// Getting configuration from configuration files and from CLI parameters
func configuration() Configuration {
	var c Configuration

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("couldn't load config: %s", err)
		os.Exit(1)
	}

	c.apikey = v.GetString("apikey")

	return c
}

// Getting parameters from CLI
func cliParameters() (string, string, string, string, int, bool) {
	var (
		playlistID = flag.String("p", "", "Youtube playlist ID")
		channelID  = flag.String("c", "", "Youtube channel ID")
		output     = flag.String("o", "table", "Output format {table|markdown}")
		sorting    = flag.String("s", "total-interest", "Sorting {total-interest|positive-interest|global-buzz-index|total-reaction|positive-negative-coefficient|pnc}")
		maxResults = flag.Int("m", 0, "The maximum number of items that should be returned")
		debug      = flag.Bool("d", false, "Debug mode")
	)
	flag.Parse()
	if *playlistID == "" && *channelID == "" {
		log.Fatalln("Playlist ID or channel ID must be defined")
	} else if *playlistID != "" && *channelID != "" {
		log.Fatalln("Playlist ID & channel ID could not be used together")
	}
	if *output != "table" && *output != "markdown" {
		log.Fatalln("Output format unknown")
	}
	if *sorting != "likes" && *sorting != "total-interest" && *sorting != "positive-interest" && *sorting != "global-buzz-index" && *sorting != "total-reaction" && *sorting != "pnc" && *sorting != "positive-negative-coefficient" {
		log.Fatalln("Unknown sorting column")
	}

	return *channelID, *playlistID, *output, *sorting, *maxResults, *debug
}

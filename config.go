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
	c.maxResults = v.GetString("maxresults")

	return c
}

// Getting parameters from CLI
func cliParameters() (string, string, string, string, int, bool) {
	var (
		playlistID = flag.String("p", "", "Youtube playlist ID")
		channelID  = flag.String("c", "", "Youtube channel ID")
		output     = flag.String("o", "table", "Output format")
		sorting    = flag.String("s", "likes", "Sorting")
		onPage     = flag.Int("n", 0, "Count of items to show")
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
	if *sorting != "likes" && *sorting != "total-interest" && *sorting != "positive-interest" {
		log.Fatalln("Unknown sorting column")
	}
	return *channelID, *playlistID, *output, *sorting, *onPage, *debug
}

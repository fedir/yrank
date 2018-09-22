package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

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
func cliParameters() (string, string, string) {
	var (
		playlistKey = flag.String("p", "", "Youtube playlist key")
		output      = flag.String("o", "table", "Output format")
		sorting     = flag.String("s", "likes", "Sorting")
	)
	flag.Parse()
	if *playlistKey == "" {
		log.Fatalln("Playlist key must be defined")
	}
	if *output != "table" && *output != "markdown" {
		log.Fatalln("Output format unknown")
	}
	if *sorting != "likes" && *sorting != "total-interest" && *sorting != "positive-interest" {
		log.Fatalln("Unknown sorting column")
	}
	return *playlistKey, *output, *sorting
}

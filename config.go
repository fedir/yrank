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

// Getting playlist key from CLI
func playlistKey() string {
	var (
		playlistKey = flag.String("p", "", "Youtube playlist key")
	)
	flag.Parse()
	if *playlistKey == "" {
		log.Fatalln("Playlist key must be defined")
	}
	return *playlistKey
}

package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Configuration global container
type Configuration struct {
	apikey string
}

func main() {
	c := configuration()
	fmt.Printf("API key: %s\n", c.apikey)
}

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

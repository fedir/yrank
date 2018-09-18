package main

import (
	"fmt"
)

func main() {
	c := configuration()
	fmt.Printf("API key: %s\n", c.apikey)
}

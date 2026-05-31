package youtube

import (
	"os"
)

func loadRespFromFile(file string) []byte {
	resp, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return resp
}

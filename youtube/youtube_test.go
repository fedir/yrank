package youtube

import (
	"io/ioutil"
)

func loadRespFromFile(file string) []byte {
	resp, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return resp
}

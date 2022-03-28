package shared

import (
	"encoding/json"
	"log"
)

//
func Manifest(data string) string {
	var manifestJson string

	json.Unmarshal([]byte(data), &manifestJson)
	log.Print(manifestJson)

	return "NotImplemented"
}

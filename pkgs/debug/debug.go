package pkg_debug

import (
	"encoding/json"
	"log"
)

func Debug(data any) {
	json, _ := json.Marshal(data)
	log.Printf("[package][debug]: %v \n", string(json))
}

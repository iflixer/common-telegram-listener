package helper

import (
	"encoding/json"
	"log"
)

func P(s interface{}) {
	log.Println("debug:")
	enc := json.NewEncoder(log.Writer())
	enc.SetIndent("", "  ")
	_ = enc.Encode(s)
}

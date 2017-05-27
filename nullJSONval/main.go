package main

import (
	"encoding/json"
	"fmt"
)

// Payload activity to be recorded in Keen
type Payload struct {
	Device string `json:"device" binding:"required"`
	Value  uint64 `json:"value" binding:"required"`
}

func main() {
	value := uint64(0)
	key := "lkjslk"

	p := new(Payload)

	p.Value = value
	p.Device = key

	jsonMsg, _ := json.Marshal(p)

	fmt.Println(string(jsonMsg))

}

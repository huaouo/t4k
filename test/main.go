package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	j := "{\"a\":12, \"b\": 10086}"
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(j), &m)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(uint64(m["a"].(float64)))
}

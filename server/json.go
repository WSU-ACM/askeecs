package main

import (
	"encoding/json"
)

type JM map[string]string

func Stringify(i interface{}) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(b)
}

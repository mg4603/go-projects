package main

import (
	"encoding/json"
	"fmt"
)

func compareJson(json1, json2 string) (bool, error) {
	var obj1, obj2 Movie
	if err := json.Unmarshal([]byte(json1), &obj1); err != nil {
		return false, fmt.Errorf("error parsing first JSON: %v\n", err)
	}
	if err := json.Unmarshal([]byte(json2), &obj2); err != nil {
		return false, fmt.Errorf("error parsing second JSON: %v\n", err)
	}
	return obj1 == obj2, nil
}

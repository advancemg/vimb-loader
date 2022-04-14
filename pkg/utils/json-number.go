package utils

import (
	"encoding/json"
	"strings"
)

// Convert json.Number to int64, float64
func JsonNumber(value interface{}) interface{} {
	if number, ok := value.(json.Number); ok {
		dot := strings.Contains(number.String(), ".")
		if dot {
			i, err := number.Float64()
			if err != nil {
				panic(err)
			}
			value = i
		} else {
			i, err := number.Int64()
			if err != nil {
				panic(err)
			}
			value = int(i)
		}
	}
	return value
}

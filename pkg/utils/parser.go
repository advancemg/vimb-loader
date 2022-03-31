package utils

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strconv"
	"time"
)

func Int64(input string) int64 {
	result, _ := strconv.ParseInt(input, 10, 64)
	return result
}

func Int64I(input interface{}) *int64 {
	if input == nil {
		return nil
	}
	result, _ := strconv.ParseInt(fmt.Sprintf("%v", input), 10, 64)
	return &result
}

func Int(input string) int {
	result, _ := strconv.Atoi(input)
	return result
}

func IntI(input interface{}) *int {
	if input == nil {
		return nil
	}
	result, _ := strconv.Atoi(fmt.Sprintf("%v", input))
	return &result
}

func Float(input string) float64 {
	result, _ := strconv.ParseFloat(input, 64)
	return result
}

func FloatI(input interface{}) *float64 {
	if input == nil {
		return nil
	}
	result, _ := strconv.ParseFloat(fmt.Sprintf("%v", input), 64)
	return &result
}

func BoolI(input interface{}) *bool {
	if input == nil {
		return nil
	}
	parseBool, err := strconv.ParseBool(fmt.Sprintf("%v", input))
	if err != nil {
		return nil
	}
	return &parseBool
}

func StringI(input interface{}) *string {
	if input == nil {
		return nil
	}
	sprint := fmt.Sprintf("%v", input)
	return &sprint
}

func FloatInterface(input interface{}) float64 {
	parse, _ := input.(float64)
	return parse
}
func FloatFromJson(input interface{}) float64 {
	parse, _ := input.(json.Number).Float64()
	return parse
}

func TimeI(input interface{}, layout string) *time.Time {
	if input == nil {
		return nil
	}
	parse, _ := time.Parse(layout, fmt.Sprintf("%s", input))
	return &parse
}

func Time(input string) time.Time {
	parse, _ := time.Parse(`2006-01-02 15:04:05`, input)
	return parse
}

func TimeByTimeZone(input, zone string) time.Time {
	timeZone, _ := time.LoadLocation(zone)
	parse, _ := time.ParseInLocation(`2006-01-02 15:04:05`, input, timeZone)
	return parse
}
func ToJsonBytes(input interface{}) []byte {
	marshal, err := json.Marshal(input)
	if err != nil {
		return []byte(``)
	}
	return marshal
}
func UniRequestId(inputToken string) uint64 {
	h := fnv.New64()
	h.Write([]byte(inputToken))
	return h.Sum64()
}

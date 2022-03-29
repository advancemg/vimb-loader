package utils

import (
	"encoding/json"
	"hash/fnv"
	"strconv"
	"time"
)

func Int64(input string) int64 {
	result, _ := strconv.ParseInt(input, 10, 64)
	return result
}
func Int(input string) int {
	result, _ := strconv.Atoi(input)
	return result
}
func Float(input string) float64 {
	result, _ := strconv.ParseFloat(input, 64)
	return result
}
func FloatInterface(input interface{}) float64 {
	parse, _ := input.(float64)
	return parse
}
func FloatFromJson(input interface{}) float64 {
	parse, _ := input.(json.Number).Float64()
	return parse
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

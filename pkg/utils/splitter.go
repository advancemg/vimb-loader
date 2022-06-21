package utils

import "strings"

func SplitDbAndTable(s string) (string, string) {
	split := strings.Split(s, "/")
	return split[0], split[1]
}

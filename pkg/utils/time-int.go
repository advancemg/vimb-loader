package utils

import (
	"fmt"
	"time"
)

func DateTimeNowInt() string {
	year, month, day := time.Now().Date()
	return fmt.Sprintf("%d%d%d", year, int(month), day)
}

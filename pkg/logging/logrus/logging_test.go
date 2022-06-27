package logrus

import (
	"os"
	"testing"
)

func TestPrintLog(t *testing.T) {
	_, err := os.Open("fff")
	if err != nil {
		PrintLog("vimb", "new-client", "error", err)
	}
}

func BenchmarkPrintLog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := os.Open("fff")
		if err != nil {
			PrintLog("vimb", "new-client", "info", err)
		}
	}
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = os.RemoveAll(pwd + "/" + "logs")
	if err != nil {
		panic(err)
	}
}

package zap

import (
	"os"
	"testing"
)

func TestPrintLog(t *testing.T) {
	Init()
	_, err := os.Open("fff")
	if err != nil {
		PrintLog("vimb", "new-client", "info", err)
	}
}

func BenchmarkPrintLog(b *testing.B) {
	Init()
	for i := 0; i < b.N; i++ {
		PrintLog("vimb", "new-client", "info", "err")
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

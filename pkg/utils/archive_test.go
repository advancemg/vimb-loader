package utils

import (
	"os"
	"testing"
)

func TestRecursiveZip(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = RecursiveZip(pwd, "outPath.zip")
	if err != nil {
		panic(err)
	}
	os.RemoveAll("logs")
	os.RemoveAll("outPath.zip")
}

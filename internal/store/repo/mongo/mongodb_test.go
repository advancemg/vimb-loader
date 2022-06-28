package mongo

import (
	"fmt"
	"github.com/advancemg/vimb-loader/internal/config"
	"github.com/advancemg/vimb-loader/pkg/logging/zap"
	"os"
	"testing"
	"time"
)

func TestDbRepo_AddWithTTL(t *testing.T) {
	config.Load()
	zap.Init()
	db := New("timeout", "db")
	timeout := Timeout{
		IsTimeout: true,
	}
	for i := 0; i < 10000; i++ {
		err := db.AddWithTTL("_id", timeout, time.Second*30)
		if err != nil {
			panic(err)
		}
	}
	os.RemoveAll("logs")
}

func TestGetTTL(t *testing.T) {
	config.Load()
	zap.Init()
	db := New("timeout", "db")
	var timeout Timeout
	err := db.Get("_id", &timeout)
	if err != nil {
		panic(err)
	}
	os.RemoveAll("logs")
	fmt.Println(timeout)
}

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
	if testing.Short() {
		t.SkipNow()
	}
	err := config.Load()
	if err != nil {
		panic(err)
	}
	zap.Init()
	db, _ := New("timeout", "db")
	timeout := Timeout{
		IsTimeout: true,
	}
	err = db.AddWithTTL("_id", timeout, time.Second*30)
	if err != nil {
		panic(err)
	}

	os.RemoveAll("logs")
}

func TestGetTTL(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	config.Load()
	zap.Init()
	db, _ := New("timeout", "db")
	var timeout Timeout
	err := db.Get("_id", &timeout)
	if err != nil {
		panic(err)
	}
	os.RemoveAll("logs")
	fmt.Println(timeout)
}

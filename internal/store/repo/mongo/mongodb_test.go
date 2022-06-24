package mongo

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestDbRepo_AddWithTTL(t *testing.T) {
	var db = DbRepo{
		Client:   nil,
		table:    "timeout",
		database: "db",
	}
	timeout := Timeout{
		IsTimeout: true,
	}
	err := db.AddWithTTL("_id", timeout, time.Second*30)
	if err != nil {
		panic(err)
	}
	os.RemoveAll("logs")
}

func TestGetTTL(t *testing.T) {
	var db = DbRepo{
		Client:   nil,
		table:    "timeout",
		database: "db",
	}
	var timeout Timeout
	err := db.Get("_id", &timeout)
	if err != nil {
		panic(err)
	}
	os.RemoveAll("logs")
	fmt.Println(timeout)
}

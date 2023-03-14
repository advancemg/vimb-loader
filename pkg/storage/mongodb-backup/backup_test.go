package mongodb_backup

import (
	"fmt"
	"github.com/advancemg/vimb-loader/internal/config"
	"github.com/advancemg/vimb-loader/pkg/logging/zap"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"os"
	"strings"
	"testing"
	"time"
)

func TestDump(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	err := config.Load()
	if err != nil {
		panic(err)
	}
	err = zap.Init()
	if err != nil {
		panic(err)
	}
	dm := InitConfig()
	start := time.Now()
	path, err := dm.Backup()
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Since(start))
	fmt.Println(path)
	//os.RemoveAll(path)
}

func TestRestore(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	err := config.Load()
	if err != nil {
		panic(err)
	}
	err = zap.Init()
	if err != nil {
		panic(err)
	}
	dm := InitConfig()
	start := time.Now()
	path := "/var/folders/dz/38x3jsnd74q6fg_97qrxnhzc0000gn/T/mongo-dump/db-2022-06-30T13:50:25Z.gz"
	//path := "/Users/eminshakh/tmp/db"
	err = dm.Restore(path)
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Since(start))
}

func TestDumpS3(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	err := config.Load()
	if err != nil {
		panic(err)
	}
	err = zap.Init()
	if err != nil {
		panic(err)
	}
	dm := InitConfig()
	start := time.Now()
	path, err := dm.Backup()
	if err != nil {
		panic(err)
	}
	fmt.Println(path)
	s3.InitConfig()
	index := strings.LastIndex(path, "/")
	s3Key := fmt.Sprintf("%s/%s", "mongo-backup", path[index:])
	_, err = s3.UploadFileWithBucket(path, s3Key)
	if err != nil {
		panic(err)
	}
	os.RemoveAll(path)
	fmt.Println(time.Since(start))
}

func TestListBackup(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	config.Load()
	err := zap.Init()
	if err != nil {
		panic(err)
	}
	s3.InitConfig()
	backups, _ := ListBackups()
	for _, k := range backups {
		fmt.Println(k)
	}

}

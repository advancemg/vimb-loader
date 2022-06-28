package mongodb_backup

import (
	"fmt"
	"github.com/advancemg/vimb-loader/internal/config"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"os"
	"strings"
	"testing"
	"time"
)

func TestDump(t *testing.T) {
	err := config.Load()
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

func TestRecursiveZip(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = RecursiveZip(pwd, "outPath.zip")
	if err != nil {
		panic(err)
	}
}

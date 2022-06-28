package mongodb_backup

import (
	"archive/zip"
	"fmt"
	cfg "github.com/advancemg/vimb-loader/internal/config"
	log "github.com/advancemg/vimb-loader/pkg/logging/zap"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/mongodb/mongo-tools/common/options"
	"github.com/mongodb/mongo-tools/mongodump"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Host       string `json:"Host"`
	Port       string `json:"Port"`
	DB         string `json:"Db"`
	Username   string `json:"Username"`
	Password   string `json:"Password"`
	CronBackup string `json:"CronBackup"`
}

func InitConfig() *Config {
	return &Config{
		Host:       cfg.Config.Mongo.Host,
		Port:       cfg.Config.Mongo.Port,
		DB:         cfg.Config.Mongo.DB,
		Username:   cfg.Config.Mongo.Username,
		Password:   cfg.Config.Mongo.Password,
		CronBackup: cfg.Config.Mongo.CronBackup,
	}
}

func (cfg *Config) Backup() (string, error) {
	enabledOptions := options.EnabledOptions{
		Auth:       true,
		Connection: true,
		Namespace:  true,
		URI:        true,
	}
	opts := options.New("mongodump", "", "", "", false, enabledOptions)
	outputOptions := &mongodump.OutputOptions{
		NumParallelCollections: 1,
	}
	inputOptions := &mongodump.InputOptions{}
	url := options.URI{
		ConnectionString: fmt.Sprintf(`mongodb://%s:%s/%s`, cfg.Host, cfg.Port, cfg.DB),
	}
	auth := options.Auth{
		Username: cfg.Username,
		Password: cfg.Password,
	}
	opts.URI = &url
	opts.Auth = &auth
	dump := mongodump.MongoDump{
		ToolOptions:   opts,
		InputOptions:  inputOptions,
		OutputOptions: outputOptions,
	}
	//dump.ToolOptions.Namespace.DB = cfg.DB
	path := fmt.Sprintf("%s%s/mongo-backup-%s", os.TempDir(), "mongo-dump", time.Now().UTC().Format(time.RFC3339))
	dump.OutputOptions.Out = path
	dump.OutputOptions.Gzip = true
	dump.OutputOptions.NumParallelCollections = 1
	err := dump.Init()
	if err != nil {
		return "", err
	}
	err = dump.Dump()
	if err != nil {
		return "", err
	}
	index := strings.LastIndex(path, "/")
	zipFile := fmt.Sprintf("%s%s.zip", os.TempDir(), path[index:])
	err = RecursiveZip(path, zipFile)
	if err != nil {
		return "", fmt.Errorf("RecursiveZip: %w", err)
	}
	os.RemoveAll(path)
	return zipFile, nil
}

func (cfg *Config) StartBackup() func() {
	return func() {
		path, err := cfg.Backup()
		if err != nil {
			log.PrintLog("vimb-loader", "StartBackup", "error", "err:", err.Error())
			return
		}
		fmt.Println(path)
		s3.InitConfig()
		index := strings.LastIndex(path, "/")
		s3Key := fmt.Sprintf("%s/%s", "mongo-backup", path[index:])
		_, err = s3.UploadFileWithBucket(path, s3Key)
		if err != nil {
			log.PrintLog("vimb-loader", "StartBackup", "error", "err:", err.Error())
			return
		}
		os.RemoveAll(path)
	}
}

func RecursiveZip(outputZip, inputPath string) error {
	destinationFile, err := os.Create(inputPath)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	err = filepath.Walk(outputZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(outputZip))
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}

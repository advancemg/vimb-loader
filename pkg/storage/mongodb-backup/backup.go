package mongodb_backup

import (
	"archive/zip"
	"fmt"
	cfg "github.com/advancemg/vimb-loader/internal/config"
	log "github.com/advancemg/vimb-loader/pkg/logging/zap"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/mongodb/mongo-tools/common/db"
	mlog "github.com/mongodb/mongo-tools/common/log"
	"github.com/mongodb/mongo-tools/common/options"
	"github.com/mongodb/mongo-tools/common/util"
	"github.com/mongodb/mongo-tools/mongodump"
	"github.com/mongodb/mongo-tools/mongorestore"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SwaggerBackupRequest struct {
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	DB       string `json:"Db"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type JsonResponse struct {
	Request string `json:"request"`
}

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

func (cfg *Config) Restore(path string) {
	enabledOptions := options.EnabledOptions{
		Auth:       true,
		Connection: true,
		Namespace:  true,
		URI:        true,
	}
	opts := options.New("mongodump", "", "", "", false, enabledOptions)
	nsOpts := &mongorestore.NSOptions{
		NSInclude: []string{"db.timeout"},
	}
	inputOpts := &mongorestore.InputOptions{}
	outputOpts := &mongorestore.OutputOptions{}
	opts.AddOptions(inputOpts)
	opts.AddOptions(nsOpts)
	opts.AddOptions(outputOpts)
	//opts.Namespace = &options.Namespace{DB: nsOpts.NSInclude[0]}
	//url := options.URI{
	//	ConnectionString: fmt.Sprintf(`mongodb://%s:%s/%s`, cfg.Host, cfg.Port, cfg.DB),
	//}
	//auth := options.Auth{
	//	Username: cfg.Username,
	//	Password: cfg.Password,
	//}
	//opts.URI = &url
	//opts.Auth = &auth
	targetDir := util.ToUniversalPath(path)
	provider, err := db.NewSessionProvider(*opts)
	if err != nil {
		panic(err)
	}
	defer provider.Close()
	restore := mongorestore.MongoRestore{
		ToolOptions:     opts,
		InputOptions:    inputOpts,
		OutputOptions:   outputOpts,
		NSOptions:       nsOpts,
		SessionProvider: provider,
		TargetDirectory: targetDir,
	}
	defer restore.Close()
	result := restore.Restore()
	err = result.Err
	if err != nil {
		panic(err)
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
	filePath := fmt.Sprintf("logs/mongo-backup-%v", time.Now().Format(time.RFC3339))
	logFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll("logs", 0777)
			if err != nil {
				return "", fmt.Errorf("backup: %w", err)
			}
			logFile, err = os.Create(filePath)
			if err != nil {
				return "", fmt.Errorf("backup: %w", err)
			}
		} else {
			return "", fmt.Errorf("backup: %w", err)
		}
	}
	defer logFile.Close()
	mlog.SetWriter(logFile)
	dump := mongodump.MongoDump{
		ToolOptions:   opts,
		InputOptions:  inputOptions,
		OutputOptions: outputOptions,
	}
	//dump.ToolOptions.Namespace.DB = cfg.DB
	path := fmt.Sprintf("%s%s/mongo-backup-%s", os.TempDir(), "mongo-dump", time.Now().UTC().Format(time.RFC3339))
	dump.OutputOptions.Out = path
	//dump.OutputOptions.Gzip = true
	dump.OutputOptions.NumParallelCollections = 1
	err = dump.Init()
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

func (cfg *Config) RunBackup() (*JsonResponse, error) {
	path, err := cfg.Backup()
	if err != nil {
		log.PrintLog("vimb-loader", "RunBackup", "error", "err:", err.Error())
		return nil, err
	}
	s3.InitConfig()
	index := strings.LastIndex(path, "/")
	s3Key := fmt.Sprintf("%s/%s", "mongo-backup", path[index:])
	_, err = s3.UploadFileWithBucket(path, s3Key)
	if err != nil {
		log.PrintLog("vimb-loader", "RunBackup", "error", "err:", err.Error())
		return nil, err
	}
	os.RemoveAll(path)
	return &JsonResponse{"s3Key: " + s3Key}, nil
}

func (cfg *Config) StartBackup() func() {
	return func() {
		_, err := cfg.RunBackup()
		if err != nil {
			log.PrintLog("vimb-loader", "StartBackup", "error", "err:", err.Error())
			return
		}
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

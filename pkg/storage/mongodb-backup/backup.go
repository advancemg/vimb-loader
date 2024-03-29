package mongodb_backup

import (
	"fmt"
	cfg "github.com/advancemg/vimb-loader/internal/config"
	log "github.com/advancemg/vimb-loader/pkg/logging/zap"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/mongodb/mongo-tools/common/db"
	mlog "github.com/mongodb/mongo-tools/common/log"
	"github.com/mongodb/mongo-tools/common/options"
	"github.com/mongodb/mongo-tools/mongodump"
	"github.com/mongodb/mongo-tools/mongorestore"
	"os"
	"strings"
	"time"
)

const backupDir = "mongo-backup"

type SwaggerBackupRequest struct {
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	DB       string `json:"Db"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type SwaggerListBackupsRequest struct {
}

type SwaggerRestoreRequest struct {
	S3Key string `json:"S3Key"`
}

type JsonResponse struct {
	Request string `json:"request"`
}

type Config struct {
	Host       string `json:"Host"`
	Port       string `json:"Port"`
	DB         string `json:"Db"`
	AuthDB     string `json:"AuthDb"`
	Username   string `json:"Username"`
	Password   string `json:"Password"`
	CronBackup string `json:"CronBackup"`
}

func InitConfig() *Config {
	return &Config{
		Host:       cfg.Config.Mongo.Host,
		Port:       cfg.Config.Mongo.Port,
		AuthDB:     cfg.Config.Mongo.AuthDB,
		DB:         cfg.Config.Mongo.DB,
		Username:   cfg.Config.Mongo.Username,
		Password:   cfg.Config.Mongo.Password,
		CronBackup: cfg.Config.Mongo.CronBackup,
	}
}

func (cfg *Config) options() *options.ToolOptions {
	enabledOptions := options.EnabledOptions{
		Auth:       true,
		Connection: true,
		Namespace:  true,
		URI:        true,
	}
	opts := options.New("mongo-tools", "", "", "", false, enabledOptions)
	url := options.URI{
		ConnectionString: fmt.Sprintf(`mongodb://%s:%s/%s`, cfg.Host, cfg.Port, cfg.AuthDB),
	}
	auth := options.Auth{
		Username: cfg.Username,
		Password: cfg.Password,
		Source:   cfg.AuthDB,
	}
	opts.URI = &url
	opts.Auth = &auth
	return opts
}

func (cfg *Config) Restore(s3Key string) error {
	path, err := s3.Download(s3Key)
	if err != nil {
		return err
	}
	mlog.SetWriter(log.LogWriter)
	opts := cfg.options()
	nsOpts := &mongorestore.NSOptions{}
	inputOpts := &mongorestore.InputOptions{
		Archive: path,
		Gzip:    true,
	}
	outputOpts := &mongorestore.OutputOptions{
		Drop:                   true,
		NoIndexRestore:         true,
		NumParallelCollections: 4,
		NumInsertionWorkers:    1,
	}
	opt := mongorestore.Options{
		ToolOptions:   opts,
		InputOptions:  inputOpts,
		NSOptions:     nsOpts,
		OutputOptions: outputOpts,
	}
	restore, err := mongorestore.New(opt)
	if err != nil {
		return err
	}
	defer restore.Close()
	result := restore.Restore()
	if result.Err != nil {
		return result.Err
	}
	return nil
}

func ListBackups() ([]string, error) {
	var result []string
	bucket := cfg.Config.S3.S3Bucket
	keys, err := s3.ListKeys(bucket, backupDir)
	if err != nil {
		return nil, err
	}
	for key, _ := range keys {
		result = append(result, key)
	}
	return result, nil
}

func (cfg *Config) Backup() (string, error) {
	opts := cfg.options()
	provider, err := db.NewSessionProvider(*opts)
	if err != nil {
		return "", err
	}
	defer provider.Close()
	outputOptions := &mongodump.OutputOptions{}
	inputOptions := &mongodump.InputOptions{}
	mlog.SetWriter(log.LogWriter)
	dump := mongodump.MongoDump{
		ToolOptions:     opts,
		InputOptions:    inputOptions,
		OutputOptions:   outputOptions,
		SessionProvider: provider,
	}
	dump.ToolOptions.Namespace.DB = cfg.DB
	path := fmt.Sprintf("%s%s/%s-%s.gz", os.TempDir(), "mongo-dump", cfg.DB, time.Now().UTC().Format(time.RFC3339))
	//dump.OutputOptions.Out = path
	dump.OutputOptions.Archive = path
	dump.OutputOptions.Gzip = true
	dump.OutputOptions.NumParallelCollections = 4
	err = dump.Init()
	if err != nil {
		return "", err
	}
	err = dump.Dump()
	if err != nil {
		return "", err
	}
	return path, nil
}

func (cfg *Config) RunBackup() (*JsonResponse, error) {
	path, err := cfg.Backup()
	if err != nil {
		log.PrintLog("vimb-loader", "RunBackup", "error", "err:", err.Error())
		return nil, err
	}
	s3.InitConfig()
	index := strings.LastIndex(path, "/")
	s3Key := fmt.Sprintf("%s/%s", backupDir, path[index:])
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

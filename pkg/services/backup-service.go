package services

import (
	"errors"
	backup "github.com/advancemg/vimb-loader/pkg/storage/mongodb-backup"
	"github.com/robfig/cron/v3"
	"os"
	"os/signal"
	"time"
)

type BackupService struct {
}

func (svc *BackupService) Start() error {
	cfg := backup.InitConfig()
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}
	scheduler := cron.New(cron.WithLocation(location))
	_, err = scheduler.AddFunc(cfg.CronBackup, cfg.StartBackup())
	if err != nil {
		return err
	}
	defer scheduler.Stop()
	go scheduler.Start()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	return errors.New("scheduler stop")
}

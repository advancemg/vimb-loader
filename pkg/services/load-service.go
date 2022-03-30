package services

import (
	"github.com/advancemg/vimb-loader/pkg/models"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"os"
	"os/signal"
	"time"
)

type LoadService struct {
}

func (svc *LoadService) Start() error {
	config := models.Config
	scheduler := cron.New()
	err := scheduler.AddFunc(config.Budget.Cron, config.Budget.InitJob())
	if err != nil {
		return err
	}
	err = scheduler.AddFunc(config.ProgramBreaks.Cron, config.ProgramBreaks.InitJob())
	if err != nil {
		return err
	}
	err = scheduler.AddFunc(config.ProgramBreaksLight.Cron, config.ProgramBreaksLight.InitJob())
	if err != nil {
		return err
	}
	err = scheduler.AddFunc(config.Channel.Cron, config.Channel.InitJob())
	if err != nil {
		return err
	}
	err = scheduler.AddFunc(config.Rank.Cron, config.Rank.InitJob())
	if err != nil {
		return err
	}
	err = scheduler.AddFunc(config.CustomersWithAdvertisers.Cron, config.CustomersWithAdvertisers.InitJob())
	if err != nil {
		return err
	}
	err = scheduler.AddFunc(config.Spots.Cron, config.Spots.InitJob())
	if err != nil {
		return err
	}
	err = scheduler.AddFunc(config.Mediaplan.Cron, config.Mediaplan.InitJob())
	if err != nil {
		return err
	}
	err = scheduler.AddFunc(config.AdvMessages.Cron, config.AdvMessages.InitJob())
	if err != nil {
		return err
	}
	err = scheduler.AddFunc(config.DeletedSpotInfo.Cron, config.DeletedSpotInfo.InitJob())
	if err != nil {
		return err
	}
	defer scheduler.Stop()
	scheduler.Start()
	/*jobs*/
	budgetErrChan := config.Budget.StartJob()
	programBreaksErrChan := config.ProgramBreaks.StartJob()
	programBreaksLightErrChan := config.ProgramBreaksLight.StartJob()
	mediaplanErrChan := config.Mediaplan.StartJob()
	channelErrChan := config.Channel.StartJob()
	rankErrChan := config.Rank.StartJob()
	spotsChan := config.Spots.StartJob()
	advMessagesChan := config.AdvMessages.StartJob()
	customersWithAdvertisersErrChan := config.CustomersWithAdvertisers.StartJob()
	deletedSpotInfoErrChan := config.DeletedSpotInfo.StartJob()
	allStartJobs := true
	for allStartJobs {
		select {
		case <-budgetErrChan:
			allStartJobs = false
		case <-programBreaksErrChan:
			allStartJobs = false
		case <-programBreaksLightErrChan:
			allStartJobs = false
		case <-mediaplanErrChan:
			allStartJobs = false
		case <-channelErrChan:
			allStartJobs = false
		case <-rankErrChan:
			allStartJobs = false
		case <-customersWithAdvertisersErrChan:
			allStartJobs = false
		case <-spotsChan:
			allStartJobs = false
		case <-advMessagesChan:
			allStartJobs = false
		case <-deletedSpotInfoErrChan:
			allStartJobs = false
		case <-time.After(time.Second * 3):
			continue
		}
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	return errors.New("scheduler stop")
}

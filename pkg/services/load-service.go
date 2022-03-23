package services

import (
	"github.com/advancemg/vimb-loader/pkg/models"
	"github.com/robfig/cron"
)

type LoadService struct {
}

func (svc *LoadService) Start() error {
	config := models.Config
	scheduler := cron.New()
	err := scheduler.AddFunc(config.Budget.Cron, config.Budget.GetJob())
	if err != nil {
		return err
	}
	err = scheduler.AddFunc(config.Channel.Cron, config.Channel.GetJob())
	if err != nil {
		return err
	}
	err = scheduler.AddFunc(config.CustomersWithAdvertisers.Cron, config.CustomersWithAdvertisers.GetJob())
	if err != nil {
		return err
	}
	err = scheduler.AddFunc(config.Mediaplan.Cron, config.Mediaplan.GetJob())
	if err != nil {
		return err
	}
	err = scheduler.AddFunc(config.AdvMessages.Cron, config.AdvMessages.GetJob())
	if err != nil {
		return err
	}
	err = scheduler.AddFunc(config.DeletedSpotInfo.Cron, config.DeletedSpotInfo.GetJob())
	if err != nil {
		return err
	}
	err = scheduler.AddFunc(config.Rank.Cron, config.Rank.GetJob())
	if err != nil {
		return err
	}
	err = scheduler.AddFunc(config.ProgramBreaks.Cron, config.ProgramBreaks.GetJob())
	if err != nil {
		return err
	}
	err = scheduler.AddFunc(config.ProgramBreaksLight.Cron, config.ProgramBreaksLight.GetJob())
	if err != nil {
		return err
	}
	err = scheduler.AddFunc(config.Spots.Cron, config.Spots.GetJob())
	if err != nil {
		return err
	}
	defer scheduler.Stop()
	scheduler.Start()
	return nil
}

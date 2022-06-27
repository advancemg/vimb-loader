package services

import (
	cfg "github.com/advancemg/vimb-loader/internal/config"
	"github.com/advancemg/vimb-loader/internal/models"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"os"
	"os/signal"
	"time"
)

type LoadService struct {
}

type Configuration struct {
	Mediaplan                models.MediaplanConfiguration                `json:"mediaplan"`
	Budget                   models.BudgetConfiguration                   `json:"budget"`
	Channel                  models.ChannelConfiguration                  `json:"channel"`
	AdvMessages              models.AdvMessagesConfiguration              `json:"advMessages"`
	CustomersWithAdvertisers models.CustomersWithAdvertisersConfiguration `json:"customersWithAdvertisers"`
	DeletedSpotInfo          models.DeletedSpotInfoConfiguration          `json:"deletedSpotInfo"`
	Rank                     models.RanksConfiguration                    `json:"rank"`
	ProgramBreaks            models.ProgramBreaksConfiguration            `json:"programBreaks"`
	ProgramBreaksLight       models.ProgramBreaksLightConfiguration       `json:"programBreaksLight"`
	Spots                    models.SpotsConfiguration                    `json:"spots"`
}

func InitConfig() *Configuration {
	return &Configuration{
		Mediaplan:                models.MediaplanConfiguration(cfg.Config.Mediaplan),
		Budget:                   models.BudgetConfiguration(cfg.Config.Budget),
		Channel:                  models.ChannelConfiguration(cfg.Config.Channel),
		AdvMessages:              models.AdvMessagesConfiguration(cfg.Config.AdvMessages),
		CustomersWithAdvertisers: models.CustomersWithAdvertisersConfiguration(cfg.Config.CustomersWithAdvertisers),
		DeletedSpotInfo:          models.DeletedSpotInfoConfiguration(cfg.Config.DeletedSpotInfo),
		Rank:                     models.RanksConfiguration(cfg.Config.Rank),
		ProgramBreaks:            models.ProgramBreaksConfiguration(cfg.Config.ProgramBreaks),
		ProgramBreaksLight:       models.ProgramBreaksLightConfiguration(cfg.Config.ProgramBreaksLight),
		Spots:                    models.SpotsConfiguration(cfg.Config.Spots),
	}
}

func (svc *LoadService) Start() error {
	config := InitConfig()
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

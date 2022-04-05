package services

import (
	"errors"
	"github.com/advancemg/vimb-loader/pkg/models"
	"os"
	"os/signal"
	"time"
)

type UpdateService struct {
}

func (svc *UpdateService) Start() error {
	budgetErrChan := models.BudgetStartJob()
	channelErrChan := models.ChannelsStartJob()
	advertiserErrChan := models.AdvertiserStartJob()
	customersWithAdvertisersErrChan := models.CustomersWithAdvertisersStartJob()
	spotErrChan := models.SpotStartJob()
	mediaplanErrChan := models.MediaplanStartJob()
	programBreaksLightErrChan := models.ProgramBreaksLightStartJob()
	programBreaksStartJobErrChan := models.ProgramBreaksStartJob()
	ranksErrChan := models.RanksStartJob()
	programsErrChan := models.ProgramStartJob()
	mediaplanAggErrChan := models.MediaplanAggStartJob()
	allStartJobs := true
	for allStartJobs {
		select {
		case <-budgetErrChan:
			allStartJobs = false
		case <-channelErrChan:
			allStartJobs = false
		case <-advertiserErrChan:
			allStartJobs = false
		case <-customersWithAdvertisersErrChan:
			allStartJobs = false
		case <-spotErrChan:
			allStartJobs = false
		case <-spotErrChan:
			allStartJobs = false
		case <-mediaplanErrChan:
			allStartJobs = false
		case <-programBreaksLightErrChan:
			allStartJobs = false
		case <-programBreaksStartJobErrChan:
			allStartJobs = false
		case <-ranksErrChan:
			allStartJobs = false
		case <-programsErrChan:
			allStartJobs = false
		case <-mediaplanAggErrChan:
			allStartJobs = false
		case <-time.After(time.Second * 3):
			continue
		}
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	return errors.New("update service stop")
}

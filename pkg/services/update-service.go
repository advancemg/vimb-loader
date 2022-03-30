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
	allStartJobs := true
	for allStartJobs {
		select {
		case <-budgetErrChan:
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

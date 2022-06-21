package models

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/internal/usecase"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

type ProgramUpdateRequest struct {
}

type Program struct {
	CnlID        *int64    `json:"CnlID"`
	ProgID       *int64    `json:"ProgID"`
	ProID        *int64    `json:"ProID"`
	Pro2         *int64    `json:"Pro2"`
	RPID         *int64    `json:"RPID"`
	PrgName      *string   `json:"PrgName"`
	CnlName      *string   `json:"CnlName"`
	PrgNameShort *string   `json:"PrgNameShort"`
	Timestamp    time.Time `json:"Timestamp"`
}

func (program *Program) Key() string {
	return fmt.Sprintf("%d", *program.ProgID)
}

func ProgramStartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := ProgramsUpdateQueue
		amqpConfig := mq_broker.InitConfig()
		err := amqpConfig.DeclareSimpleQueue(qName)
		if err != nil {
			errorCh <- err
		}
		ch, err := amqpConfig.Channel()
		if err != nil {
			errorCh <- err
		}
		defer ch.Close()
		err = ch.Qos(1, 0, false)
		messages, err := ch.Consume(qName, "",
			false,
			false,
			false,
			false,
			nil)
		for msg := range messages {
			var bodyJson ProgramUpdateRequest
			err := json.Unmarshal(msg.Body, &bodyJson)
			if err != nil {
				errorCh <- err
			}
			err = bodyJson.Update()
			if err != nil {
				errorCh <- err
			}
			msg.Ack(false)
		}
		defer close(errorCh)
	}()
	return errorCh
}

func (request *ProgramUpdateRequest) Update() error {
	timestamp := time.Now()
	db, table := utils.SplitDbAndTable(DbPrograms)
	dbPrograms := usecase.OpenDb(db, table)
	var programBreaks []ProgramBreaks
	err := dbPrograms.FindWhereGe(&programBreaks, "ProgID", int64(-1))
	if err != nil {
		return err
	}
	for _, programBreak := range programBreaks {
		program := &Program{
			CnlID:        programBreak.CnlID,
			ProgID:       programBreak.ProgID,
			ProID:        programBreak.ProID,
			Pro2:         programBreak.Pro2,
			RPID:         programBreak.RPID,
			PrgName:      programBreak.PrgName,
			CnlName:      programBreak.CnlName,
			PrgNameShort: programBreak.PrgNameShort,
			Timestamp:    timestamp,
		}
		err = dbPrograms.AddOrUpdate(program.Key(), program)
		if err != nil {
			return err
		}
	}
	return nil
}

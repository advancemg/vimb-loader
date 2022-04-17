package models

import (
	"encoding/json"
	"fmt"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

type ProgramBreaksLightUpdateRequest struct {
	S3Key string
}

type ProgramBreaksLight struct {
	Booked              *int64    `json:"Booked"`
	BlockID             *int64    `json:"BlockID"`
	RankID              *int64    `json:"RankID"`
	VM                  *int64    `json:"VM"`
	VR                  *int64    `json:"VR"`
	SimpleFixVolume     *int64    `json:"SimpleFixVolume"`
	ReserveVol          *int64    `json:"ReserveVol"`
	SimpleFixReserveVol *int64    `json:"SimpleFixReserveVol"`
	BlockDur            *int64    `json:"BlockDur"`
	BlkOpenReserv       *int64    `json:"BlkOpenReserv"`
	HasAucSpots         *bool     `json:"HasAucSpots"`
	Timestamp           time.Time `json:"Timestamp"`
}

func (program *ProgramBreaksLight) Key() string {
	return fmt.Sprintf("%d", *program.BlockID)
}

func (b *internalB) Convert() (*ProgramBreaksLight, error) {
	timestamp := time.Now()
	items := &ProgramBreaksLight{
		Booked:              utils.Int64I(b.B["Booked"]),
		BlockID:             utils.Int64I(b.B["BlockID"]),
		RankID:              utils.Int64I(b.B["RankID"]),
		VM:                  utils.Int64I(b.B["VM"]),
		VR:                  utils.Int64I(b.B["VR"]),
		SimpleFixVolume:     utils.Int64I(b.B["SimpleFixVolume"]),
		ReserveVol:          utils.Int64I(b.B["ReserveVol"]),
		SimpleFixReserveVol: utils.Int64I(b.B["SimpleFixReserveVol"]),
		BlockDur:            utils.Int64I(b.B["BlockDur"]),
		BlkOpenReserv:       utils.Int64I(b.B["BlkOpenReserv"]),
		HasAucSpots:         utils.BoolI(b.B["HasAucSpots"]),
		Timestamp:           timestamp,
	}
	return items, nil
}

func ProgramBreaksLightStartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := ProgramBreaksLightModeUpdateQueue
		amqpConfig := mq_broker.InitConfig()
		err := amqpConfig.DeclareSimpleQueue(qName)
		if err != nil {
			errorCh <- err
		}
		ch, err := amqpConfig.Channel()
		if err != nil {
			errorCh <- err
		}
		err = ch.Qos(1, 0, false)
		messages, err := ch.Consume(qName, "",
			false,
			false,
			false,
			false,
			nil)
		for msg := range messages {
			var bodyJson MqUpdateMessage
			err := json.Unmarshal(msg.Body, &bodyJson)
			if err != nil {
				errorCh <- err
			}
			/*read from s3 by s3Key*/
			req := ProgramBreaksLightUpdateRequest{
				S3Key: bodyJson.Key,
			}
			err = req.Update()
			if err != nil {
				errorCh <- err
			}
			msg.Ack(false)
		}
		defer close(errorCh)
	}()
	return errorCh
}

func (request *ProgramBreaksLightUpdateRequest) Update() error {
	var err error
	request.S3Key, err = s3.Download(request.S3Key)
	if err != nil {
		return err
	}
	err = request.loadFromFile()
	if err != nil {
		return err
	}
	return nil
}

func (request *ProgramBreaksLightUpdateRequest) loadFromFile() error {
	resp := utils.VimbResponse{FilePath: request.S3Key}
	convertData, err := resp.Convert("BreakList")
	if err != nil {
		return err
	}
	marshalData, err := json.Marshal(convertData)
	if err != nil {
		return err
	}
	var internalData []internalB
	err = json.Unmarshal(marshalData, &internalData)
	if err != nil {
		return err
	}
	badgerProgramBreaksLight := storage.Open(DbProgramBreaksLightMode)
	for _, dataB := range internalData {
		programBreaksLight, err := dataB.Convert()
		if err != nil {
			return err
		}
		err = badgerProgramBreaksLight.Upsert(programBreaksLight.Key(), programBreaksLight)
		if err != nil {
			return err
		}
	}
	return nil
}

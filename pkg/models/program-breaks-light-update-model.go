package models

import (
	"encoding/json"
	"fmt"
	log "github.com/advancemg/vimb-loader/pkg/logging"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"os"
	"reflect"
	"time"
)

type ProgramBreaksLightUpdateRequest struct {
	S3Key string
}

type ProgramBreaksLight struct {
	Booked    []BookedAttributes `json:"Booked"`
	BlockID   *int64             `json:"BlockID"`
	VM        *int64             `json:"VM"`
	VR        *int64             `json:"VR"`
	Timestamp time.Time          `json:"Timestamp"`
}

func (program *ProgramBreaksLight) Key() string {
	return fmt.Sprintf("%d", *program.BlockID)
}

func (b *internalB) Convert() (*ProgramBreaksLight, error) {
	timestamp := time.Now()
	var blockID *int64
	var vm *int64
	var vr *int64
	var attribute internalAttr
	var attributes []BookedAttributes
	if _, ok := b.B["Booked"]; ok {
		marshalData, err := json.Marshal(b.B["Booked"])
		if err != nil {
			return nil, err
		}
		switch reflect.TypeOf(b.B["Booked"]).Kind() {
		case reflect.Array, reflect.Slice:
			var internalAttributesData []internalAttributes
			err = json.Unmarshal(marshalData, &internalAttributesData)
			if err != nil {
				return nil, err
			}
			for _, qualityItem := range internalAttributesData {
				attr, err := qualityItem.Convert()
				if err != nil {
					return nil, err
				}
				attributes = append(attributes, *attr)
			}
		case reflect.Map, reflect.Struct:
			var internalAttributesData internalAttributes
			err = json.Unmarshal(marshalData, &internalAttributesData)
			if err != nil {
				return nil, err
			}
			attr, err := internalAttributesData.Convert()
			if err != nil {
				return nil, err
			}
			attributes = append(attributes, *attr)
		}
	}
	if _, ok := b.B["attributes"]; ok {
		marshalData, err := json.Marshal(b.B["attributes"])
		if err != nil {
			return nil, err
		}
		switch reflect.TypeOf(b.B["attributes"]).Kind() {
		case reflect.Map, reflect.Struct:
			err = json.Unmarshal(marshalData, &attribute)
			if err != nil {
				return nil, err
			}
			blockID = utils.Int64I(attribute.BlockID)
			vm = utils.Int64I(attribute.VM)
			vr = utils.Int64I(attribute.VR)
		}
	} else {
		blockID = utils.Int64I(b.B["BlockID"])
		vm = utils.Int64I(b.B["VM"])
		vr = utils.Int64I(b.B["VR"])
	}
	items := &ProgramBreaksLight{
		Booked:    attributes,
		BlockID:   blockID,
		VM:        vm,
		VR:        vr,
		Timestamp: timestamp,
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
	for {
		var err error
		request.S3Key, err = s3.Download(request.S3Key)
		if err != nil {
			return err
		}
		open, err := os.Open(request.S3Key)
		if err != nil {
			if os.IsNotExist(err) {
				log.PrintLog("ProgramBreaksLightUpdateRequest", "Update()", "error", "Empty s3Key ", err.Error())
				time.Sleep(time.Minute * 2)
				continue
			} else {
				return err
			}
		}
		open.Close()
		err = request.loadFromFile()
		if err != nil {
			return err
		}
		return nil
	}
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

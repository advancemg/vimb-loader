package models

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/internal/store"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"reflect"
	"time"
)

type RanksUpdateRequest struct {
	S3Key string
}

type internalRank struct {
	Rank map[string]interface{} `json:"Rank"`
}
type internalDetail struct {
	Detail map[string]interface{} `json:"detail"`
}

type Ranks struct {
	ID        *int64    `json:"ID" bson:"ID"`
	Name      *string   `json:"Name" bson:"Name"`
	Timestamp time.Time `json:"Timestamp" bson:"Timestamp"`
	Details   []Detail  `json:"Details" bson:"Details"`
}

type Detail struct {
	SellingDirectionID *int64     `json:"SellingDirectionID" bson:"SellingDirectionID"`
	OrderNo            *int64     `json:"OrderNo" bson:"OrderNo"`
	UsesAuction        *bool      `json:"UsesAuction" bson:"UsesAuction"`
	StartDate          *time.Time `json:"StartDate" bson:"StartDate"`
	EndDate            *time.Time `json:"EndDate" bson:"EndDate"`
}

func (r *Ranks) Key() string {
	return fmt.Sprintf("%d", r.ID)
}

func (d *internalDetail) Convert() (*Detail, error) {
	detail := &Detail{
		SellingDirectionID: utils.Int64I(d.Detail["SellingDirectionID"]),
		OrderNo:            utils.Int64I(d.Detail["OrderNo"]),
		UsesAuction:        utils.BoolI(d.Detail["UsesAuction"]),
		StartDate:          utils.TimeI(d.Detail["StartDate"], `2006-01-02T15:04:05`),
		EndDate:            utils.TimeI(d.Detail["EndDate"], `2006-01-02T15:04:05`),
	}
	return detail, nil
}

func (r *internalRank) Convert() (*Ranks, error) {
	timestamp := time.Now()
	var details []Detail
	if _, ok := r.Rank["Details"]; ok {
		marshalData, err := json.Marshal(r.Rank["Details"])
		if err != nil {
			return nil, err
		}
		switch reflect.TypeOf(r.Rank["Details"]).Kind() {
		case reflect.Array, reflect.Slice:
			var internalDetailData []internalDetail
			err = json.Unmarshal(marshalData, &internalDetailData)
			if err != nil {
				return nil, err
			}
			for _, qualityItem := range internalDetailData {
				detail, err := qualityItem.Convert()
				if err != nil {
					return nil, err
				}
				details = append(details, *detail)
			}
		case reflect.Map, reflect.Struct:
			var internalDetailData internalDetail
			err = json.Unmarshal(marshalData, &internalDetailData)
			if err != nil {
				return nil, err
			}
			detail, err := internalDetailData.Convert()
			if err != nil {
				return nil, err
			}
			details = append(details, *detail)
		}
	}
	budget := &Ranks{
		ID:        utils.Int64I(r.Rank["ID"]),
		Name:      utils.StringI(r.Rank["Name"]),
		Timestamp: timestamp,
		Details:   details,
	}
	return budget, nil
}

func RanksStartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := RanksUpdateQueue
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
			var bodyJson MqUpdateMessage
			err := json.Unmarshal(msg.Body, &bodyJson)
			if err != nil {
				errorCh <- err
			}
			/*read from s3 by s3Key*/
			req := RanksUpdateRequest{
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

func (request *RanksUpdateRequest) Update() error {
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

func (request *RanksUpdateRequest) loadFromFile() error {
	resp := utils.VimbResponse{FilePath: request.S3Key}
	convertData, err := resp.Convert("Ranks")
	if err != nil {
		return err
	}
	marshalData, err := json.Marshal(convertData)
	if err != nil {
		return err
	}
	var internalData []internalRank
	err = json.Unmarshal(marshalData, &internalData)
	if err != nil {
		return err
	}
	db, table := utils.SplitDbAndTable(DbRanks)
	dbRanks := store.OpenDb(db, table)
	for _, dataM := range internalData {
		rank, err := dataM.Convert()
		if err != nil {
			return err
		}
		err = dbRanks.AddOrUpdate(rank.Key(), rank)
		if err != nil {
			return err
		}
	}
	return nil
}

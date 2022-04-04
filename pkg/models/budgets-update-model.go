package models

import (
	"encoding/json"
	"fmt"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"reflect"
	"time"
)

const BudgetTable = "budgets"

type Budget struct {
	Month                 *int         `json:"Month"`
	CnlID                 *int         `json:"CnlID"`
	AdtID                 *int         `json:"AdtID"`
	AgrID                 *int         `json:"AgrID"`
	InventoryUnitDuration *int         `json:"InventoryUnitDuration"`
	DealChannelStatus     *int         `json:"DealChannelStatus"`
	FixPercent            *int         `json:"FixPercent"`
	GRPFix                *int         `json:"GRPFix"`
	AdtName               *string      `json:"AdtName"`
	AgrName               *string      `json:"AgrName"`
	CmpName               *string      `json:"CmpName"`
	CnlName               *string      `json:"CnlName"`
	TP                    *string      `json:"TP"`
	Budget                *float64     `json:"Budget"`
	CoordCost             *float64     `json:"CoordCost"`
	Cost                  *float64     `json:"Cost"`
	FixPercentPrime       *float64     `json:"FixPercentPrime"`
	FloatPercent          *float64     `json:"FloatPercent"`
	FloatPercentPrime     *float64     `json:"FloatPercentPrime"`
	GRP                   *float64     `json:"GRP"`
	GRPWithoutKF          *float64     `json:"GRPWithoutKF"`
	Timestamp             time.Time    `json:"Timestamp"`
	Quality               []BudgetItem `json:"Quality"`
}

type BudgetItem struct {
	RankID            *int     `json:"RankID"`
	Percent           *float64 `json:"Percent"`
	BudgetOffprime    *float64 `json:"BudgetOffprime"`
	BudgetPrime       *float64 `json:"BudgetPrime"`
	InventoryOffprime *float64 `json:"InventoryOffprime"`
	InventoryPrime    *float64 `json:"InventoryPrime"`
	PercentPrime      *float64 `json:"PercentPrime"`
}

type BudgetsUpdateRequest struct {
	S3Key string
}

func (budget *Budget) Key() string {
	return fmt.Sprintf("%d-%d-%d-%d", *budget.Month, *budget.CnlID, *budget.AdtID, *budget.AgrID)
}

func (item *internalItem) Convert() (*BudgetItem, error) {
	items := &BudgetItem{
		RankID:            utils.IntI(item.Item["RankID"]),
		Percent:           utils.FloatI(item.Item["Percent"]),
		BudgetOffprime:    utils.FloatI(item.Item["BudgetOffprime"]),
		BudgetPrime:       utils.FloatI(item.Item["BudgetPrime"]),
		InventoryOffprime: utils.FloatI(item.Item["InventoryOffprime"]),
		InventoryPrime:    utils.FloatI(item.Item["InventoryPrime"]),
		PercentPrime:      utils.FloatI(item.Item["PercentPrime"]),
	}
	return items, nil
}

func (m *internalM) ConvertBudget() (*Budget, error) {
	timestamp := time.Now()
	var qualities []BudgetItem
	if _, ok := m.M["Quality"]; ok {
		marshalData, err := json.Marshal(m.M["Quality"])
		if err != nil {
			return nil, err
		}
		switch reflect.TypeOf(m.M["Quality"]).Kind() {
		case reflect.Array, reflect.Slice:
			var internalQualityData []internalItem
			err = json.Unmarshal(marshalData, &internalQualityData)
			if err != nil {
				return nil, err
			}
			for _, qualityItem := range internalQualityData {
				quality, err := qualityItem.Convert()
				if err != nil {
					return nil, err
				}
				qualities = append(qualities, *quality)
			}
		case reflect.Map, reflect.Struct:
			var internalQualityData internalItem
			err = json.Unmarshal(marshalData, &internalQualityData)
			if err != nil {
				return nil, err
			}
			quality, err := internalQualityData.Convert()
			if err != nil {
				return nil, err
			}
			qualities = append(qualities, *quality)
		}
	}
	budget := &Budget{
		Month:                 utils.IntI(m.M["Month"]),
		CnlID:                 utils.IntI(m.M["CnlID"]),
		AdtID:                 utils.IntI(m.M["AdtID"]),
		AgrID:                 utils.IntI(m.M["AgrID"]),
		InventoryUnitDuration: utils.IntI(m.M["InventoryUnitDuration"]),
		DealChannelStatus:     utils.IntI(m.M["DealChannelStatus"]),
		FixPercent:            utils.IntI(m.M["FixPercent"]),
		GRPFix:                utils.IntI(m.M["GRPFix"]),
		AdtName:               utils.StringI(m.M["AdtName"]),
		AgrName:               utils.StringI(m.M["AgrName"]),
		CmpName:               utils.StringI(m.M["CmpName"]),
		CnlName:               utils.StringI(m.M["CnlName"]),
		TP:                    utils.StringI(m.M["TP"]),
		Budget:                utils.FloatI(m.M["Budget"]),
		CoordCost:             utils.FloatI(m.M["CoordCost"]),
		Cost:                  utils.FloatI(m.M["Cost"]),
		FixPercentPrime:       utils.FloatI(m.M["FixPercentPrime"]),
		FloatPercent:          utils.FloatI(m.M["FloatPercent"]),
		FloatPercentPrime:     utils.FloatI(m.M["FloatPercentPrime"]),
		GRP:                   utils.FloatI(m.M["GRP"]),
		GRPWithoutKF:          utils.FloatI(m.M["GRPWithoutKF"]),
		Timestamp:             timestamp,
		Quality:               qualities,
	}
	return budget, nil
}

func BudgetStartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := BudgetsUpdateQueue
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
			req := BudgetsUpdateRequest{
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

func (request *BudgetsUpdateRequest) Update() error {
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

func (request *BudgetsUpdateRequest) loadFromFile() error {
	resp := utils.VimbResponse{FilePath: request.S3Key}
	convertData, err := resp.Convert("BudgetList")
	if err != nil {
		return err
	}
	marshalData, err := json.Marshal(convertData)
	if err != nil {
		return err
	}
	var internalData []internalM
	err = json.Unmarshal(marshalData, &internalData)
	if err != nil {
		return err
	}
	badgerBudgets := storage.Open(DbBudgets)
	for _, dataM := range internalData {
		budget, err := dataM.ConvertBudget()
		if err != nil {
			return err
		}
		err = badgerBudgets.Upsert(budget.Key(), budget)
		if err != nil {
			return err
		}
	}
	return nil
}

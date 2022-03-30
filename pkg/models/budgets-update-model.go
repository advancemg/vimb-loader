package models

import (
	"encoding/json"
	"fmt"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

const BudgetTable = "budgets"

type Budget struct {
	Month                 int          `json:"month"`
	CnlID                 int          `json:"cnlID"`
	AdtID                 int          `json:"adtID"`
	AgrID                 int          `json:"agrID"`
	InventoryUnitDuration int          `json:"inventoryUnitDuration"`
	DealChannelStatus     int          `json:"dealChannelStatus"`
	FixPercent            int          `json:"fixPercent"`
	GRPFix                int          `json:"GRPFix"`
	AdtName               string       `json:"adtName"`
	AgrName               string       `json:"agrName"`
	CmpName               string       `json:"cmpName"`
	CnlName               string       `json:"cnlName"`
	TP                    string       `json:"TP"`
	Budget                float64      `json:"budget"`
	CoordCost             float64      `json:"coordCost"`
	Cost                  float64      `json:"cost"`
	FixPercentPrime       float64      `json:"fixPercentPrime"`
	FloatPercent          float64      `json:"floatPercent"`
	FloatPercentPrime     float64      `json:"floatPercentPrime"`
	GRP                   float64      `json:"GRP"`
	GRPWithoutKF          float64      `json:"GRPWithoutKF"`
	Quality               []BudgetItem `json:"quality"`
}

type BudgetItem struct {
	Percent           int     `json:"percent"`
	RankID            int     `json:"rankID"`
	BudgetOffprime    float64 `json:"budgetOffprime"`
	BudgetPrime       float64 `json:"budgetPrime"`
	InventoryOffprime float64 `json:"inventoryOffprime"`
	InventoryPrime    float64 `json:"inventoryPrime"`
	PercentPrime      float64 `json:"percentPrime"`
}

type BudgetsUpdateRequest struct {
	S3Key string
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
	filePath, err := s3.Download(request.S3Key)
	if err != nil {
		return err
	}
	resp := utils.VimbResponse{FilePath: filePath}
	catchField, err := resp.Convert("BudgetList")
	if err != nil {
		return err
	}
	var budget = Budget{}
	var quality = []BudgetItem{}
	badgerBudgets := storage.NewBadger(DbBudgets)
	badgerMonth := storage.NewBadger(DbCustomConfigMonth)
	badgerAdvertisers := storage.NewBadger(DbCustomConfigAdvertisers)
	badgerChannels := storage.NewBadger(DbCustomConfigChannels)
	defer badgerBudgets.Close()
	defer badgerMonth.Close()
	defer badgerAdvertisers.Close()
	defer badgerChannels.Close()
	if array1, ok := catchField.([]interface{}); ok {
		for _, arrVal1 := range array1 {
			if mapLevel1, ok := arrVal1.(map[string]interface{}); ok {
				for _, v1 := range mapLevel1 {
					if mapLevel2, ok := v1.(map[string]interface{}); ok {
						for _, v2 := range mapLevel2 {
							if array2, ok := v2.([]interface{}); ok {
								quality = []BudgetItem{}
								for _, arrVal2 := range array2 {
									if mapLevel3, ok := arrVal2.(map[string]interface{}); ok {
										for _, v3 := range mapLevel3 {
											if mapLevel4, ok := v3.(map[string]interface{}); ok {
												item := BudgetItem{
													Percent:           utils.Int(mapLevel4["Percent"].(string)),
													RankID:            utils.Int(mapLevel4["RankID"].(string)),
													BudgetOffprime:    utils.Float(mapLevel4["BudgetOffprime"].(string)),
													BudgetPrime:       utils.Float(mapLevel4["BudgetPrime"].(string)),
													InventoryOffprime: utils.Float(mapLevel4["InventoryOffprime"].(string)),
													InventoryPrime:    utils.Float(mapLevel4["InventoryPrime"].(string)),
													PercentPrime:      utils.Float(mapLevel4["PercentPrime"].(string)),
												}
												quality = append(quality, item)
											}
										}
									}
								}
							}
						}
						adtId := mapLevel2["AdtID"].(string)
						cnlId := mapLevel2["CnlID"].(string)
						month := mapLevel2["Month"].(string)
						badgerAdvertisers.Set(adtId, []byte(adtId))
						badgerChannels.Set(cnlId, []byte(cnlId))
						badgerMonth.Set(month, []byte(month))
						budget = Budget{
							Month:                 utils.Int(mapLevel2["Month"].(string)),
							CnlID:                 utils.Int(mapLevel2["CnlID"].(string)),
							AdtID:                 utils.Int(mapLevel2["AdtID"].(string)),
							AgrID:                 utils.Int(mapLevel2["AgrID"].(string)),
							InventoryUnitDuration: utils.Int(mapLevel2["InventoryUnitDuration"].(string)),
							DealChannelStatus:     utils.Int(mapLevel2["DealChannelStatus"].(string)),
							FixPercent:            utils.Int(mapLevel2["FixPercent"].(string)),
							GRPFix:                utils.Int(mapLevel2["GRPFix"].(string)),
							AdtName:               mapLevel2["AdtName"].(string),
							AgrName:               mapLevel2["AgrName"].(string),
							CmpName:               mapLevel2["CmpName"].(string),
							CnlName:               mapLevel2["CnlName"].(string),
							TP:                    mapLevel2["TP"].(string),
							Budget:                utils.Float(mapLevel2["Budget"].(string)),
							CoordCost:             utils.Float(mapLevel2["CoordCost"].(string)),
							Cost:                  utils.Float(mapLevel2["Cost"].(string)),
							FixPercentPrime:       utils.Float(mapLevel2["FixPercentPrime"].(string)),
							FloatPercent:          utils.Float(mapLevel2["FloatPercent"].(string)),
							FloatPercentPrime:     utils.Float(mapLevel2["FloatPercentPrime"].(string)),
							GRP:                   utils.Float(mapLevel2["GRP"].(string)),
							GRPWithoutKF:          utils.Float(mapLevel2["GRPWithoutKF"].(string)),
							Quality:               quality,
						}
					}
				}
			}
			key := fmt.Sprintf("%d%d%d%d", budget.Month, budget.CnlID, budget.AdtID, budget.AgrID)
			body, err := json.Marshal(budget)
			if err != nil {
				return err
			}
			badgerBudgets.Set(key, body)
		}
	}
	return nil
}

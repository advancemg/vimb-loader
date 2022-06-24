package models

import (
	"encoding/json"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/internal/store"
	log "github.com/advancemg/vimb-loader/pkg/logging"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type ChannelLoadRequest struct {
	SellingDirectionID string `json:"SellingDirectionID" example:"23"`
}

type ChannelQuery struct {
	ID struct {
		Eq int64 `json:"eq" example:"1018583"`
	} `json:"ID"`
	MainChnl struct {
		Eq int64 `json:"eq" example:"1018568"`
	} `json:"MainChnl"`
	SellingDirectionID struct {
		Ee int64 `json:"eq" example:"23"`
	} `json:"SellingDirectionID"`
}

func (request *ChannelLoadRequest) InitTasks() (CommonResponse, error) {
	qName := GetChannelsType
	amqpConfig := mq_broker.InitConfig()
	err := amqpConfig.DeclareSimpleQueue(qName)
	if err != nil {
		return nil, err
	}
	result := CommonResponse{}
	req := goConvert.New()
	req.Set("SellingDirectionID", request.SellingDirectionID)
	err = amqpConfig.PublishJson(qName, req)
	if err != nil {
		log.PrintLog("vimb-loader", "Channels InitTasks", "error", "Q:", qName, "err:", err.Error())
		return nil, err
	}
	result["status"] = "ok"
	return result, nil
}

func (request *Any) QueryChannels() ([]Channel, error) {
	var result []Channel
	db, table := utils.SplitDbAndTable(DbChannels)
	repo := store.OpenDb(db, table)
	marshal, err := json.Marshal(request.Body)
	if err != nil {
		return nil, err
	}
	err = repo.FindJson(&result, marshal)
	if err != nil {
		return nil, err
	}
	return result, nil
}

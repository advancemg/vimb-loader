package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/advancemg/badgerhold"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/internal/usecase"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

type SwaggerChangeMPlanFilmPlannedInventoryRequest struct {
	Data []struct {
		CommInMpl []struct {
			ID        string `json:"ID"`
			Inventory string `json:"Inventory"`
		} `json:"CommInMpl"`
	} `json:"data"`
}

type ChangeMPlanFilmPlannedInventory struct {
	goConvert.UnsortedMap
}

func (request *ChangeMPlanFilmPlannedInventory) GetDataJson() (*JsonResponse, error) {
	req, err := request.getXml()
	if err != nil {
		return nil, err
	}
	resp, err := utils.Actions.RequestJson(req)
	if err != nil {
		return nil, err
	}
	var body = map[string]interface{}{}
	err = json.Unmarshal(resp, &body)
	if err != nil {
		return nil, err
	}
	return &JsonResponse{
		Body:    body,
		Request: string(req),
	}, nil
}

func (request *ChangeMPlanFilmPlannedInventory) GetDataXmlZip() (*StreamResponse, error) {
	for {
		var isTimeout utils.Timeout
		db, table := utils.SplitDbAndTable(DbTimeout)
		repo := usecase.OpenDb(db, table)
		err := repo.Get("vimb-timeout", &isTimeout)
		if err != nil {
			if errors.Is(err, badgerhold.ErrNotFound) {
				isTimeout.IsTimeout = false
			} else {
				return nil, err
			}
		}
		if isTimeout.IsTimeout {
			time.Sleep(1 * time.Second)
			continue
		}
		if !isTimeout.IsTimeout {
			break
		}
	}
	req, err := request.getXml()
	if err != nil {
		return nil, err
	}
	resp, err := utils.Actions.Request(req)
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    resp,
		Request: string(req),
	}, nil
}

func (request *ChangeMPlanFilmPlannedInventory) UploadToS3() error {
	for {
		typeName := ChangeMPlanFilmPlannedInventoryType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				vimbError.CheckTimeout("ChangeMPlanFilmPlannedInventory")
				continue
			}
			return err
		}
		var newS3Key = fmt.Sprintf("vimb/%s/%s/%s-%s.gz", utils.Actions.Client, typeName, utils.DateTimeNowInt(), typeName)
		_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
		if err != nil {
			return err
		}
		return nil
	}
}

func (request *ChangeMPlanFilmPlannedInventory) getXml() ([]byte, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	Data, exist := request.Get("Data")
	if exist {
		body.Set("Data", Data)
	}
	xmlRequestHeader.Set("ChangeMPlanFilmPlannedInventory", body)
	return xmlRequestHeader.ToXml()
}

package models

import (
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
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

func (request *ChangeMPlanFilmPlannedInventory) GetData() (*StreamResponse, error) {
	req, err := request.getXml()
	if err != nil {
		return nil, err
	}
	resp, err := utils.RequestJson(req)
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    resp,
		Request: string(req),
	}, nil
}

func (request *ChangeMPlanFilmPlannedInventory) GetRawData() (*StreamResponse, error) {
	req, err := request.getXml()
	if err != nil {
		return nil, err
	}
	resp, err := utils.Request(req)
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    resp,
		Request: string(req),
	}, nil
}

func (request *ChangeMPlanFilmPlannedInventory) GetDataToS3() error {
	data, err := request.GetRawData()
	if err != nil {
		return err
	}
	var newS3Key = fmt.Sprintf("%s.zip", "GetProgramBreaks")
	_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
	if err != nil {
		return err
	}
	return nil
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

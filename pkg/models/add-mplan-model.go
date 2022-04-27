package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/advancemg/badgerhold"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type SwaggerAddMPlanRequest struct {
	OrdID             string `json:"OrdID"`
	MplCnlID          string `json:"MplCnlID"`
	DateFrom          string `json:"DateFrom"`
	DateTo            string `json:"DateTo"`
	MplName           string `json:"MplName"`
	BrandID           string `json:"BrandID"`
	MultiSpotsInBlock string `json:"MultiSpotsInBlock"`
}

type AddMPlan struct {
	goConvert.UnsortedMap
}

func (request *AddMPlan) GetDataJson() (*JsonResponse, error) {
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

func (request *AddMPlan) GetDataXmlZip() (*StreamResponse, error) {
	for {
		var isTimeout utils.Timeout
		err := storage.Open(DbTimeout).Get("vimb-timeout", &isTimeout)
		if err != nil {
			if errors.Is(err, badgerhold.ErrNotFound) {
				isTimeout.IsTimeout = false
			} else {
				return nil, err
			}
		}
		if isTimeout.IsTimeout {
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

func (request *AddMPlan) UploadToS3() error {
	for {
		typeName := AddMPlanType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				vimbError.CheckTimeout("AddMPlan")
				continue
			}
			return err
		}
		month, _ := request.Get("DateFrom")
		var newS3Key = fmt.Sprintf("vimb/%s/%s/%v/%s-%s.gz", utils.Actions.Client, typeName, month, utils.DateTimeNowInt(), typeName)
		_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
		if err != nil {
			return err
		}
		return nil
	}
}

func (request *AddMPlan) getXml() ([]byte, error) {
	attributes := goConvert.New()
	attributes.Set("xmlns:xsi", "\"http://www.w3.org/2001/XMLSchema-instance\"")
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	OrdID, exist := request.Get("OrdID")
	if exist {
		body.Set("OrdID", OrdID)
	}
	MplCnlID, exist := request.Get("MplCnlID")
	if exist {
		body.Set("MplCnlID", MplCnlID)
	}
	DateFrom, exist := request.Get("DateFrom")
	if exist {
		body.Set("DateFrom", DateFrom)
	}
	DateTo, exist := request.Get("DateTo")
	if exist {
		body.Set("DateTo", DateTo)
	}
	MplName, exist := request.Get("MplName")
	if exist {
		body.Set("MplName", MplName)
	}
	BrandID, exist := request.Get("BrandID")
	if exist {
		body.Set("BrandID", BrandID)
	}
	MultiSpotsInBlock, exist := request.Get("MultiSpotsInBlock")
	if exist {
		body.Set("MultiSpotsInBlock", MultiSpotsInBlock)
	}
	xmlRequestHeader.Set("AddMPlan", body)
	xmlRequestHeader.Set("attributes", attributes)
	return xmlRequestHeader.ToXml()
}

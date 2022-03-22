package models

import (
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/pkg/s3"
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

func (request *AddMPlan) GetDataJson() (*StreamResponse, error) {
	req, err := request.getXml()
	if err != nil {
		return nil, err
	}
	resp, err := utils.Actions.RequestJson(req)
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    resp,
		Request: string(req),
	}, nil
}

func (request *AddMPlan) GetDataXmlZip() (*StreamResponse, error) {
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
	typeName := AddMPlanType
	data, err := request.GetDataXmlZip()
	if err != nil {
		return err
	}
	var newS3Key = fmt.Sprintf("vimb/%s/%s/%s-%s.gz", utils.Actions.Client, typeName, utils.DateTimeNowInt(), typeName)
	_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
	if err != nil {
		return err
	}
	return nil
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

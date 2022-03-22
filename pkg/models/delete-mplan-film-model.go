package models

import (
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type SwaggerDeleteMPlanFilmRequest struct {
	CommInMplID string `json:"CommInMplID"`
}

type DeleteMPlanFilm struct {
	goConvert.UnsortedMap
}

func (request *DeleteMPlanFilm) GetDataJson() (*StreamResponse, error) {
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

func (request *DeleteMPlanFilm) GetDataXmlZip() (*StreamResponse, error) {
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

func (request *DeleteMPlanFilm) UploadToS3() error {
	typeName := DeleteMPlanFilmType
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

func (request *DeleteMPlanFilm) getXml() ([]byte, error) {
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	CommInMplID, exist := request.Get("CommInMplID")
	if exist {
		body.Set("CommInMplID", CommInMplID)
	}
	xmlRequestHeader.Set("DeleteMPlanFilm", body)
	return xmlRequestHeader.ToXml()
}

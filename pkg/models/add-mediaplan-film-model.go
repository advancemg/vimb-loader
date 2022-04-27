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
	"time"
)

type SwaggerAddMPlanFilmRequest struct {
	MplID     string `json:"MplID"`
	FilmID    string `json:"FilmID"`
	StartDate string `json:"StartDate"`
	EndDate   string `json:"EndDate"`
}

type AddMPlanFilm struct {
	goConvert.UnsortedMap
}

func (request *AddMPlanFilm) GetDataJson() (*JsonResponse, error) {
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

func (request *AddMPlanFilm) GetDataXmlZip() (*StreamResponse, error) {
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

func (request *AddMPlanFilm) UploadToS3() error {
	for {
		typeName := AddMPlanFilmType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				vimbError.CheckTimeout("AddMPlanFilm")
				continue
			}
			return err
		}
		month, _ := request.Get("StartDate")
		var newS3Key = fmt.Sprintf("vimb/%s/%s/%s/%s-%s.gz", utils.Actions.Client, typeName, month, utils.DateTimeNowInt(), typeName)
		_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
		if err != nil {
			return err
		}
		return nil
	}
}

func (request *AddMPlanFilm) getXml() ([]byte, error) {
	attributes := goConvert.New()
	attributes.Set("xmlns:xsi", "\"http://www.w3.org/2001/XMLSchema-instance\"")
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	MplID, exist := request.Get("MplID")
	if exist {
		body.Set("MplID", MplID)
	}
	FilmID, exist := request.Get("FilmID")
	if exist {
		body.Set("FilmID", FilmID)
	}
	StartDate, exist := request.Get("StartDate")
	if exist {
		body.Set("StartDate", StartDate)
	}
	EndDate, exist := request.Get("EndDate")
	if exist {
		body.Set("EndDate", EndDate)
	}
	xmlRequestHeader.Set("AddMPlanFilm", body)
	xmlRequestHeader.Set("attributes", attributes)
	return xmlRequestHeader.ToXml()
}

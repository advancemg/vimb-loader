package models

import (
	"encoding/json"
	"errors"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/internal/usecase"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

type SwaggerAddSpotRequest struct {
	BlockID         string `json:"BlockID"`
	FilmID          string `json:"FilmID"`
	Position        string `json:"Position"`
	FixedPosition   string `json:"FixedPosition"`
	AuctionBidValue string `json:"AuctionBidValue"`
}

type AddSpot struct {
	goConvert.UnsortedMap
}

func (request *AddSpot) GetDataJson() (*JsonResponse, error) {
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

func (request *AddSpot) GetDataXmlZip() (*StreamResponse, error) {
	for {
		var isTimeout utils.Timeout
		db, table := utils.SplitDbAndTable(DbTimeout)
		repo := usecase.OpenDb(db, table)
		err := repo.Get("vimb-timeout", &isTimeout)
		if err != nil {
			if errors.Is(err, usecase.ErrNotFound) {
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

func (request *AddSpot) UploadToS3() error {
	for {
		typeName := AddSpotType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				vimbError.CheckTimeout("AddSpot")
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

func (request *AddSpot) getXml() ([]byte, error) {
	attributes := goConvert.New()
	attributes.Set("xmlns:xsi", "\"http://www.w3.org/2001/XMLSchema-instance\"")
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	BlockID, exist := request.Get("BlockID")
	if exist {
		body.Set("BlockID", BlockID)
	}
	FilmID, exist := request.Get("FilmID")
	if exist {
		body.Set("FilmID", FilmID)
	}
	Position, exist := request.Get("Position")
	if exist {
		body.Set("Position", Position)
	}
	FixedPosition, exist := request.Get("FixedPosition")
	if exist {
		body.Set("FixedPosition", FixedPosition)
	}
	AuctionBidValue, exist := request.Get("AuctionBidValue")
	if exist {
		body.Set("AuctionBidValue", AuctionBidValue)
	}
	xmlRequestHeader.Set("AddSpot", body)
	xmlRequestHeader.Set("attributes", attributes)
	return xmlRequestHeader.ToXml()
}

package models

import (
	"encoding/json"
	"errors"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/internal/store"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
)

type SwaggerDeleteMPlanFilmRequest struct {
	CommInMplID string `json:"CommInMplID"`
}

type DeleteMPlanFilm struct {
	goConvert.UnsortedMap
}

func (request *DeleteMPlanFilm) GetDataJson() (*JsonResponse, error) {
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

func (request *DeleteMPlanFilm) GetDataXmlZip() (*StreamResponse, error) {
	for {
		var isTimeout utils.Timeout
		db, table := utils.SplitDbAndTable(DbTimeout)
		repo := store.OpenDb(db, table)
		err := repo.Get("_id", &isTimeout)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
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

func (request *DeleteMPlanFilm) UploadToS3() error {
	for {
		typeName := DeleteMPlanFilmType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				vimbError.CheckTimeout("DeleteMPlanFilm")
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

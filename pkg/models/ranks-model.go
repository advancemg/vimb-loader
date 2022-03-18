package models

import (
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type SwaggerGetRanksRequest struct {
}

type GetRanks struct {
	goConvert.UnsortedMap
}

func (request *GetRanks) GetData() (*StreamResponse, error) {
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

func (request *GetRanks) getXml() ([]byte, error) {
	xmlRequestHeader := goConvert.New()
	xmlRequestHeader.Set("GetRanks", "")
	return xmlRequestHeader.ToXml()
}

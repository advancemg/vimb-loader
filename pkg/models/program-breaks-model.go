package models

import (
	goConvert "github.com/advancemg/go-convert"
)

type SwaggerGetProgramBreaksRequest struct {
	SellingDirectionID string `json:"sellingDirectionID"`
	InclProgAttr       string `json:"inclProgAttr"`
	InclForecast       string `json:"inclForecast"`
	AudRatDec          string `json:"audRatDec"`
	StartDate          string `json:"startDate"`
	EndDate            string `json:"endDate"`
	LightMode          string `json:"lightMode"`
	CnlList            []struct {
		Cnl string `json:"Cnl"`
	} `json:"CnlList"`
	ProtocolVersion string `json:"protocolVersion"`
}

type GetProgramBreaks struct {
	goConvert.UnsortedMap
}

func (request GetProgramBreaks) GetData() (*StreamResponse, error) {
	attributes := goConvert.New()
	attributes.Set("xmlns:xsi", "\"http://www.w3.org/2001/XMLSchema-instance\"")
	xmlRequestHeader := goConvert.New()
	body := goConvert.New()
	sellingDirectionID, exist := request.Get("SellingDirectionID")
	if exist {
		body.Set("SellingDirectionID", sellingDirectionID)
	}
	inclProgAttr, exist := request.Get("InclProgAttr")
	if exist {
		body.Set("InclProgAttr", inclProgAttr)
	}
	inclForecast, exist := request.Get("InclForecast")
	if exist {
		body.Set("InclForecast", inclForecast)
	}
	audRatDec, exist := request.Get("AudRatDec")
	if exist {
		body.Set("AudRatDec", audRatDec)
	}
	startDate, exist := request.Get("StartDate")
	if exist {
		body.Set("StartDate", startDate)
	}
	endDate, exist := request.Get("EndDate")
	if exist {
		body.Set("EndDate", endDate)
	}
	lightMode, exist := request.Get("LightMode")
	if exist {
		body.Set("LightMode", lightMode)
	}
	cnlList, exist := request.Get("CnlList")
	if exist {
		body.Set("CnlList", cnlList)
	}
	protocolVersion, exist := request.Get("ProtocolVersion")
	if exist {
		body.Set("ProtocolVersion", protocolVersion)
	}
	xmlRequestHeader.Set("GetProgramBreaks", body)
	xmlRequestHeader.Set("attributes", attributes)
	req, err := xmlRequestHeader.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(req),
	}, nil
}

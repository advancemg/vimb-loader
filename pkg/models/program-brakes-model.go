package models

import goConvert "github.com/advancemg/go-convert"

type SwaggerGetProgramBrakesRequest struct {
	SellingDirectionID string `json:"SellingDirectionID"`
	InclProgAttr       string `json:"InclProgAttr"`
	InclForecast       string `json:"InclForecast"`
	AudRatDec          string `json:"AudRatDec"`
	StartDate          string `json:"StartDate"`
	EndDate            string `json:"EndDate"`
	LightMode          string `json:"LightMode"`
	CnlList            []struct {
		Cnl string `json:"Cnl"`
	} `json:"CnlList"`
	ProtocolVersion string `json:"ProtocolVersion"`
}

type GetProgramBrakes struct {
	goConvert.UnsortedMap
}

func (request *GetProgramBrakes) GetData() (*StreamResponse, error) {
	xml, err := request.ToXml()
	if err != nil {
		return nil, err
	}
	return &StreamResponse{
		Body:    nil,
		Request: string(xml),
	}, nil
}

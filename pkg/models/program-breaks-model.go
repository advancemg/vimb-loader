package models

import (
	goConvert "github.com/advancemg/go-convert"
	"github.com/advancemg/vimb-loader/pkg/utils"
)

type SwaggerGetProgramBreaksRequest struct {
	SellingDirectionID string `json:"SellingDirectionID" example:"21"` //ID направления продаж
	InclProgAttr       string `json:"InclProgAttr" example:"1"`        //Флаг "Заполнять секцию ProMaster". 1 - да, 0 - нет. (int, not nillable)
	InclForecast       string `json:"InclForecast" example:"1"`        //Признак "Как заполнять секцию прогнозных рейтингов". 0 - Не заполнять,  1 - Заполнять только ЦА программатика, 2 - Заполнять всеми возможными ЦА
	AudRatDec          string `json:"AudRatDec" example:"2"`           //Внимание! Этот элемент является устаревшим. Его значение игнорируется (при расчетах всегда используется значение 9) и скоро он будет удален. Кол-во десятичных знаков округления. Допустимы значения 1..9. По умолчанию - 9. (int, nillable)
	StartDate          string `json:"StartDate" example:"20170701"`    //Дата начала периода в формате YYYYMMDD
	EndDate            string `json:"EndDate" example:"20170702"`      //Дата окончания периода (включительно) в формате YYYYMMDD
	LightMode          string `json:"LightMode" example:"0"`           //Флаг "облегченного режима". Выходной xml имеет другой формат
	CnlList            []struct {
		Cnl string `json:"Cnl" example:"1018566"` //ID канала (int, not nillable)
	} `json:"CnlList"`
	ProtocolVersion string `json:"ProtocolVersion" example:"2"` //Флаг "ожидаемый формат ответа". (int, nillable). Допускается отсутствие этого элемента, в этом случае используется формат по умолчанию (1.0). 1 - старый формат (1.0), 2 - Новый формат (2.0), null - формат по-умолчанию (1.0). Внимание! В ближайших версиях этот элемент станет обязательным и не нулабельным.
}

type GetProgramBreaks struct {
	goConvert.UnsortedMap
}

func (request *GetProgramBreaks) GetData() (*StreamResponse, error) {
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

func (request *GetProgramBreaks) getXml() ([]byte, error) {
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
	return xmlRequestHeader.ToXml()
}

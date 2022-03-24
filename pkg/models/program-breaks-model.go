package models

import (
	"encoding/json"
	"fmt"
	goConvert "github.com/advancemg/go-convert"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"time"
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
	Path int
	goConvert.UnsortedMap
}

type ProgramBreaksConfiguration struct {
	Cron             string `json:"cron"`
	SellingDirection string `json:"sellingDirection"`
	Loading          bool   `json:"loading"`
}

func (cfg *ProgramBreaksConfiguration) StartJob() error {
	if !cfg.Loading {
		return nil
	}
	qName := GetProgramBreaksType
	amqpConfig := mq_broker.InitConfig()
	err := amqpConfig.DeclareSimpleQueue(qName)
	if err != nil {
		return err
	}
	ch, err := amqpConfig.Channel()
	if err != nil {
		return err
	}
	err = ch.Qos(1, 0, false)
	messages, err := ch.Consume(qName, "",
		false,
		false,
		false,
		false,
		nil)
	for msg := range messages {
		var bodyJson GetProgramBreaks
		err := json.Unmarshal(msg.Body, &bodyJson)
		if err != nil {
			return err
		}
		err = bodyJson.UploadToS3()
		if err != nil {
			return err
		}
		msg.Ack(false)
	}
	return nil
}

func (cfg *ProgramBreaksConfiguration) InitJob() func() {
	return func() {
		if !cfg.Loading {
			return
		}
		qName := GetProgramBreaksType
		amqpConfig := mq_broker.InitConfig()
		err := amqpConfig.DeclareSimpleQueue(qName)
		if err != nil {
			fmt.Printf("Q:%s - err:%s", qName, err.Error())
			return
		}
		qInfo, err := amqpConfig.GetQueueInfo(qName)
		if err != nil {
			fmt.Printf("Q:%s - err:%s", qName, err.Error())
			return
		}
		if qInfo.Messages > 0 {
			return
		}
		months, err := utils.GetActualMonths()
		if err != nil {
			fmt.Printf("Q:%s - err:%s", qName, err.Error())
			return
		}
		for _, month := range months {
			request := goConvert.New()
			request.Set("SellingDirectionID", cfg.SellingDirection)
			request.Set("StartDate", month.ValueString)
			request.Set("EndDate", month.ValueString)
			err := amqpConfig.PublishJson(qName, request)
			if err != nil {
				fmt.Printf("Q:%s - err:%s", qName, err.Error())
				return
			}
		}
	}
}

type ProgramBreaksLightConfiguration struct {
	Cron             string `json:"cron"`
	SellingDirection string `json:"sellingDirection"`
	Loading          bool   `json:"loading"`
}

func (cfg *ProgramBreaksLightConfiguration) StartJob() error {
	if !cfg.Loading {
		return nil
	}
	qName := GetProgramBreaksLightModeType
	amqpConfig := mq_broker.InitConfig()
	err := amqpConfig.DeclareSimpleQueue(qName)
	if err != nil {
		return err
	}
	ch, err := amqpConfig.Channel()
	if err != nil {
		return err
	}
	err = ch.Qos(1, 0, false)
	messages, err := ch.Consume(qName, "",
		false,
		false,
		false,
		false,
		nil)
	for msg := range messages {
		var bodyJson GetProgramBreaks
		err := json.Unmarshal(msg.Body, &bodyJson)
		if err != nil {
			return err
		}
		err = bodyJson.UploadToS3()
		if err != nil {
			return err
		}
		msg.Ack(false)
	}
	return nil
}

func (cfg *ProgramBreaksLightConfiguration) InitJob() func() {
	return func() {
		if !cfg.Loading {
			return
		}
		qName := GetProgramBreaksLightModeType
		amqpConfig := mq_broker.InitConfig()
		err := amqpConfig.DeclareSimpleQueue(qName)
		if err != nil {
			fmt.Printf("Q:%s - err:%s", qName, err.Error())
			return
		}
		qInfo, err := amqpConfig.GetQueueInfo(qName)
		if err != nil {
			fmt.Printf("Q:%s - err:%s", qName, err.Error())
			return
		}
		if qInfo.Messages > 0 {
			return
		}
		months, err := utils.GetActualMonths()
		if err != nil {
			fmt.Printf("Q:%s - err:%s", qName, err.Error())
			return
		}
		for _, month := range months {
			request := goConvert.New()
			request.Set("SellingDirectionID", cfg.SellingDirection)
			request.Set("StartDate", month.ValueString)
			request.Set("EndDate", month.ValueString)
			err := amqpConfig.PublishJson(qName, request)
			if err != nil {
				fmt.Printf("Q:%s - err:%s", qName, err.Error())
				return
			}
		}
	}
}

func (request *GetProgramBreaks) GetDataJson() (*JsonResponse, error) {
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

func (request *GetProgramBreaks) GetDataXmlZip() (*StreamResponse, error) {
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

func (request *GetProgramBreaks) UploadToS3() error {
	for {
		typeName := GetProgramBreaksType
		data, err := request.GetDataXmlZip()
		if err != nil {
			if vimbError, ok := err.(*utils.VimbError); ok {
				code := vimbError.Code
				switch code {
				case 1001:
					fmt.Printf("Vimb code %v timeout...", code)
					time.Sleep(time.Minute * 1)
					continue
				case 1003:
					fmt.Printf("Vimb code %v timeout...", code)
					time.Sleep(time.Minute * 2)
					continue
				default:
					fmt.Printf("Vimb code %v - not implemented timeout...", code)
					time.Sleep(time.Minute * 1)
					continue
				}
			}
			return err
		}
		sellingDirectionID, _ := request.Get("SellingDirectionID")
		startDate, _ := request.Get("StartDate")
		month, _ := request.Get("StartMonth")
		var newS3Key = fmt.Sprintf("vimb/%s/%s/%s/%s/%s/%d/%s-%s.gz", sellingDirectionID, utils.Actions.Client, typeName, startDate, month, request.Path, utils.DateTimeNowInt(), typeName)
		_, err = s3.UploadBytesWithBucket(newS3Key, data.Body)
		if err != nil {
			return err
		}
		return nil
	}
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

package models

import (
	"encoding/json"
	"fmt"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"reflect"
	"strconv"
	"time"
)

type ProgramBreaksUpdateRequest struct {
	S3Key string
	Month string
}

type ProgramBreaks struct {
	Month                  *int               `json:"Moth"`
	CnlID                  *int               `json:"CnlID"`
	ProgID                 *int               `json:"ProgID"`
	RCID                   *int               `json:"RCID"`
	RPID                   *int               `json:"RPID"`
	IssTime                *int               `json:"IssTime"`
	IssDuration            *int               `json:"IssDuration"`
	BlockTime              *int               `json:"BlockTime"`
	BlockNumber            *int               `json:"BlockNumber"`
	IsPrime                *int               `json:"IsPrime"`
	ProID                  *int               `json:"ProID"`
	ProOriginalPTR         *int               `json:"ProOriginalPTR"`
	BlockDistr             *int               `json:"BlockDistr"`
	SptOptions             *int               `json:"SptOptions"`
	TNSBlockFactTime       *int               `json:"TNSBlockFactTime"`
	TNSBlockFactDur        *int               `json:"TNSBlockFactDur"`
	Pro2                   *int               `json:"Pro2"`
	BlkAdvertTypePTR       *int               `json:"BlkAdvertTypePTR"`
	WeekDay                *int               `json:"WeekDay"`
	PrgBegTimeL            *int               `json:"PrgBegTimeL"`
	PrgEndMonthL           *int               `json:"PrgEndMonthL"`
	PrgBegMonthL           *int               `json:"PrgBegMonthL"`
	BlkAuc                 *int               `json:"BlkAuc"`
	AucRate                *int               `json:"AucRate"`
	NoRating               *int               `json:"NoRating"`
	AvailableAuctionVolume *int               `json:"AvailableAuctionVolume"`
	IsLocal                *int               `json:"IsLocal"`
	RootRCID               *int               `json:"RootRCID"`
	RootRPID               *int               `json:"RootRPID"`
	IsSpecialProject       *int               `json:"IsSpecialProject"`
	RankID                 *int64             `json:"RankID"`
	IssID                  *int64             `json:"IssID"`
	TNSBlockFactID         *int64             `json:"TNSBlockFactID"`
	ForecastRateBase       *float64           `json:"ForecastRateBase"`
	FactRateBase           *float64           `json:"FactRateBase"`
	RateAll                *float64           `json:"RateAll"`
	AuctionStepPrice       *float64           `json:"AuctionStepPrice"`
	PrgName                *string            `json:"PrgName"`
	CnlName                *string            `json:"CnlName"`
	PrgNameShort           *string            `json:"PrgNameShort"`
	TgrID                  *string            `json:"TgrID"`
	TgrName                *string            `json:"TgrName"`
	BlockDate              *int               `json:"BlockDate"`
	Booked                 []BookedAttributes `json:"Booked"`
	BlockID                *int64             `json:"BlockID"`
	VM                     *int               `json:"VM"`
	VR                     *int               `json:"VR"`
	DLDate                 *time.Time         `json:"DLDate"`
	DLTrDate               *time.Time         `json:"DLTrDate"`
	Timestamp              time.Time          `json:"Timestamp"`
}

type internalAttr struct {
	BlockID string `json:"BlockID"`
	VM      string `json:"VM"`
	VR      string `json:"VR"`
}

type Attributes struct {
	BlockID *int64 `json:"BlockID"`
	VM      *int   `json:"VM"`
	VR      *int   `json:"VR"`
}

func (a *internalAttr) Convert() (*Attributes, error) {
	attributes := &Attributes{
		BlockID: utils.Int64I(a.BlockID),
		VM:      utils.IntI(a.VM),
		VR:      utils.IntI(a.VR),
	}
	return attributes, nil
}

type BookedAttributes struct {
	RankID *int64 `json:"RankID"`
	VM     *int   `json:"VM"`
	VR     *int   `json:"VR"`
}

func (a *internalAttributes) Convert() (*BookedAttributes, error) {
	attributes := &BookedAttributes{
		RankID: utils.Int64I(a.Attributes["RankID"]),
		VM:     utils.IntI(a.Attributes["VM"]),
		VR:     utils.IntI(a.Attributes["VR"]),
	}
	return attributes, nil
}

type ProMaster struct {
	ProID     *int      `json:"ProID"`
	PropName  *string   `json:"PropName"`
	PropValue *string   `json:"PropValue"`
	Timestamp time.Time `json:"Timestamp"`
}

func (p *internalP) Convert() (*ProMaster, error) {
	timestamp := time.Now()
	pro := &ProMaster{
		ProID:     utils.IntI(p.P["ProID"]),
		PropName:  utils.StringI(p.P["PropName"]),
		PropValue: utils.StringI(p.P["PropValue"]),
		Timestamp: timestamp,
	}
	return pro, nil
}

type BlockForecast struct {
	TgrID            *int      `json:"TgrID"`
	BlockID          *int64    `json:"BlockID"`
	Forecast         *float64  `json:"Forecast"`
	InternetForecast *float64  `json:"InternetForecast"`
	Fact             *float64  `json:"Fact"`
	ForecastQuality  *string   `json:"ForecastQuality"`
	Timestamp        time.Time `json:"Timestamp"`
}

func (bb *internalBB) Convert() (*BlockForecast, error) {
	timestamp := time.Now()
	block := &BlockForecast{
		TgrID:            utils.IntI(bb.B["TgrID"]),
		BlockID:          utils.Int64I(bb.B["BlockID"]),
		Forecast:         utils.FloatI(bb.B["Forecast"]),
		InternetForecast: utils.FloatI(bb.B["InternetForecast"]),
		Fact:             utils.FloatI(bb.B["Fact"]),
		ForecastQuality:  utils.StringI(bb.B["ForecastQuality"]),
		Timestamp:        timestamp,
	}
	return block, nil
}

type BlockForecastTgr struct {
	ID        *int      `json:"ID"`
	Name      *string   `json:"Name"`
	Timestamp time.Time `json:"Timestamp"`
}

func (t *internalTgr) Convert() (*BlockForecastTgr, error) {
	timestamp := time.Now()
	block := &BlockForecastTgr{
		ID:        utils.IntI(t.Tgr["ID"]),
		Name:      utils.StringI(t.Tgr["Name"]),
		Timestamp: timestamp,
	}
	return block, nil
}

func (program *ProgramBreaks) Key() string {
	return fmt.Sprintf("%d", *program.BlockID)
}
func (pro *ProMaster) Key() string {
	return fmt.Sprintf("%d-%s-%s", *pro.ProID, *pro.PropName, *pro.PropValue)
}
func (block *BlockForecast) Key() string {
	return fmt.Sprintf("%d", *block.BlockID)
}
func (trg *BlockForecastTgr) Key() string {
	return fmt.Sprintf("%d", *trg.ID)
}

func (b *internalB) ConvertProgramBreaks() (*ProgramBreaks, error) {
	timestamp := time.Now()
	var blockID *int64
	var vm *int
	var vr *int
	var attribute internalAttr
	var attributes []BookedAttributes
	if _, ok := b.B["Booked"]; ok {
		marshalData, err := json.Marshal(b.B["Booked"])
		if err != nil {
			return nil, err
		}
		switch reflect.TypeOf(b.B["Booked"]).Kind() {
		case reflect.Array, reflect.Slice:
			var internalAttributesData []internalAttributes
			err = json.Unmarshal(marshalData, &internalAttributesData)
			if err != nil {
				return nil, err
			}
			for _, qualityItem := range internalAttributesData {
				attr, err := qualityItem.Convert()
				if err != nil {
					return nil, err
				}
				attributes = append(attributes, *attr)
			}
		case reflect.Map, reflect.Struct:
			var internalAttributesData internalAttributes
			err = json.Unmarshal(marshalData, &internalAttributesData)
			if err != nil {
				return nil, err
			}
			attr, err := internalAttributesData.Convert()
			if err != nil {
				return nil, err
			}
			attributes = append(attributes, *attr)
		}
	}
	if _, ok := b.B["attributes"]; ok {
		marshalData, err := json.Marshal(b.B["attributes"])
		if err != nil {
			return nil, err
		}
		switch reflect.TypeOf(b.B["attributes"]).Kind() {
		case reflect.Map, reflect.Struct:
			err = json.Unmarshal(marshalData, &attribute)
			if err != nil {
				return nil, err
			}
			blockID = utils.Int64I(attribute.BlockID)
			vm = utils.IntI(attribute.VM)
			vr = utils.IntI(attribute.VR)
		}
	} else {
		blockID = utils.Int64I(b.B["BlockID"])
		vm = utils.IntI(b.B["VM"])
		vr = utils.IntI(b.B["VR"])
	}
	items := &ProgramBreaks{
		CnlID:                  utils.IntI(b.B["CnlID"]),
		ProgID:                 utils.IntI(b.B["ProgID"]),
		RCID:                   utils.IntI(b.B["RCID"]),
		RPID:                   utils.IntI(b.B["RPID"]),
		IssTime:                utils.IntI(b.B["IssTime"]),
		IssDuration:            utils.IntI(b.B["IssDuration"]),
		BlockTime:              utils.IntI(b.B["BlockTime"]),
		BlockNumber:            utils.IntI(b.B["BlockNumber"]),
		IsPrime:                utils.IntI(b.B["IsPrime"]),
		ProID:                  utils.IntI(b.B["ProID"]),
		ProOriginalPTR:         utils.IntI(b.B["ProOriginalPTR"]),
		BlockDistr:             utils.IntI(b.B["BlockDistr"]),
		SptOptions:             utils.IntI(b.B["SptOptions"]),
		TNSBlockFactTime:       utils.IntI(b.B["TNSBlockFactTime"]),
		TNSBlockFactDur:        utils.IntI(b.B["TNSBlockFactDur"]),
		Pro2:                   utils.IntI(b.B["Pro2"]),
		BlkAdvertTypePTR:       utils.IntI(b.B["BlkAdvertTypePTR"]),
		WeekDay:                utils.IntI(b.B["WeekDay"]),
		PrgBegTimeL:            utils.IntI(b.B["PrgBegTimeL"]),
		PrgEndMonthL:           utils.IntI(b.B["PrgEndMonthL"]),
		PrgBegMonthL:           utils.IntI(b.B["PrgBegMonthL"]),
		BlkAuc:                 utils.IntI(b.B["BlkAuc"]),
		AucRate:                utils.IntI(b.B["AucRate"]),
		NoRating:               utils.IntI(b.B["NoRating"]),
		AvailableAuctionVolume: utils.IntI(b.B["AvailableAuctionVolume"]),
		IsLocal:                utils.IntI(b.B["IsLocal"]),
		RootRCID:               utils.IntI(b.B["RootRCID"]),
		RootRPID:               utils.IntI(b.B["RootRPID"]),
		IsSpecialProject:       utils.IntI(b.B["IsSpecialProject"]),
		RankID:                 utils.Int64I(b.B["RankID"]),
		IssID:                  utils.Int64I(b.B["IssID"]),
		TNSBlockFactID:         utils.Int64I(b.B["TNSBlockFactID"]),
		ForecastRateBase:       utils.FloatI(b.B["ForecastRateBase"]),
		FactRateBase:           utils.FloatI(b.B["FactRateBase"]),
		RateAll:                utils.FloatI(b.B["RateAll"]),
		AuctionStepPrice:       utils.FloatI(b.B["AuctionStepPrice"]),
		PrgName:                utils.StringI(b.B["PrgName"]),
		CnlName:                utils.StringI(b.B["CnlName"]),
		PrgNameShort:           utils.StringI(b.B["PrgNameShort"]),
		TgrID:                  utils.StringI(b.B["TgrID"]),
		TgrName:                utils.StringI(b.B["TgrName"]),
		BlockDate:              utils.IntI(b.B["BlockDate"]),
		BlockID:                blockID,
		VM:                     vm,
		VR:                     vr,
		Booked:                 attributes,
		DLDate:                 utils.TimeI(b.B["DLDate"], `2006-01-02T15:04:05`),
		DLTrDate:               utils.TimeI(b.B["DLTrDate"], `2006-01-02T15:04:05`),
		Timestamp:              timestamp,
	}
	return items, nil
}

func ProgramBreaksStartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := ProgramBreaksUpdateQueue
		amqpConfig := mq_broker.InitConfig()
		err := amqpConfig.DeclareSimpleQueue(qName)
		if err != nil {
			errorCh <- err
		}
		ch, err := amqpConfig.Channel()
		if err != nil {
			errorCh <- err
		}
		err = ch.Qos(1, 0, false)
		messages, err := ch.Consume(qName, "",
			false,
			false,
			false,
			false,
			nil)
		for msg := range messages {
			var bodyJson MqUpdateMessage
			err := json.Unmarshal(msg.Body, &bodyJson)
			if err != nil {
				errorCh <- err
			}
			/*read from s3 by s3Key*/
			req := ProgramBreaksUpdateRequest{
				S3Key: bodyJson.Key,
				Month: bodyJson.Month,
			}
			err = req.Update()
			if err != nil {
				errorCh <- err
			}
			msg.Ack(false)
		}
		defer close(errorCh)
	}()
	return errorCh
}

func (request *ProgramBreaksUpdateRequest) Update() error {
	var err error
	request.S3Key, err = s3.Download(request.S3Key)
	if err != nil {
		return err
	}
	err = request.loadFromFile()
	if err != nil {
		return err
	}
	return nil
}

func (request *ProgramBreaksUpdateRequest) loadFromFile() error {
	resp := utils.VimbResponse{FilePath: request.S3Key}
	/*BreakList*/
	convertData, err := resp.Convert("BreakList")
	if err != nil {
		return err
	}
	marshalData, err := json.Marshal(convertData)
	if err != nil {
		return err
	}
	var internalData []internalB
	err = json.Unmarshal(marshalData, &internalData)
	if err != nil {
		return err
	}
	badgerProgramBreaks := storage.Open(DbProgramBreaks)
	month, err := strconv.Atoi(request.Month)
	if err != nil {
		return err
	}
	for _, dataB := range internalData {
		programBreaks, err := dataB.ConvertProgramBreaks()
		if err != nil {
			return err
		}
		programBreaks.Month = &month
		err = badgerProgramBreaks.Upsert(programBreaks.Key(), programBreaks)
		if err != nil {
			return err
		}
	}
	/*ProMaster*/
	convertProMasterData, err := resp.Convert("ProMaster")
	if err != nil {
		return err
	}
	marshalProMasterData, err := json.Marshal(convertProMasterData)
	if err != nil {
		return err
	}
	switch reflect.TypeOf(convertProMasterData).Kind() {
	case reflect.Array, reflect.Slice:
		var internalData []internalP
		err = json.Unmarshal(marshalProMasterData, &internalData)
		if err != nil {
			return err
		}
		badgerProgramBreaksProMaster := storage.Open(DbProgramBreaksProMaster)
		for _, dataItem := range internalData {
			data, err := dataItem.Convert()
			if err != nil {
				return err
			}
			err = badgerProgramBreaksProMaster.Upsert(data.Key(), data)
			if err != nil {
				return err
			}
		}
	case reflect.Map, reflect.Struct:
		var internalData internalP
		err = json.Unmarshal(marshalProMasterData, &internalData)
		if err != nil {
			return err
		}
		badgerProgramBreaksProMaster := storage.Open(DbProgramBreaksProMaster)
		data, err := internalData.Convert()
		if err != nil {
			return err
		}
		err = badgerProgramBreaksProMaster.Upsert(data.Key(), data)
		if err != nil {
			return err
		}
	}
	/*BlockForecast*/
	convertBlockForecastData, err := resp.Convert("BlockForecast")
	if err != nil {
		return err
	}
	marshalBlockForecastData, err := json.Marshal(convertBlockForecastData)
	if err != nil {
		return err
	}
	switch reflect.TypeOf(convertBlockForecastData).Kind() {
	case reflect.Array, reflect.Slice:
		var internalData []internalBB
		err = json.Unmarshal(marshalBlockForecastData, &internalData)
		if err != nil {
			return err
		}
		badgerProgramBreaksBlockForecast := storage.Open(DbProgramBreaksBlockForecast)
		for _, dataItem := range internalData {
			data, err := dataItem.Convert()
			if err != nil {
				return err
			}
			err = badgerProgramBreaksBlockForecast.Upsert(data.Key(), data)
			if err != nil {
				return err
			}
		}
	case reflect.Map, reflect.Struct:
		var internalData internalBB
		err = json.Unmarshal(marshalBlockForecastData, &internalData)
		if err != nil {
			return err
		}
		badgerProgramBreaksBlockForecast := storage.Open(DbProgramBreaksBlockForecast)
		data, err := internalData.Convert()
		if err != nil {
			return err
		}
		err = badgerProgramBreaksBlockForecast.Upsert(data.Key(), data)
		if err != nil {
			return err
		}
	}
	/*BlockForecastTgr*/
	convertBlockForecastTgrData, err := resp.Convert("BlockForecastTgr")
	if err != nil {
		return err
	}
	marshalBlockForecastTgrData, err := json.Marshal(convertBlockForecastTgrData)
	if err != nil {
		return err
	}
	switch reflect.TypeOf(convertBlockForecastTgrData).Kind() {
	case reflect.Array, reflect.Slice:
		var internalData []internalTgr
		err = json.Unmarshal(marshalBlockForecastTgrData, &internalData)
		if err != nil {
			return err
		}
		badgerProgramBreaksBlockForecastTgr := storage.Open(DbProgramBreaksBlockForecastTgr)
		for _, dataItem := range internalData {
			data, err := dataItem.Convert()
			if err != nil {
				return err
			}
			err = badgerProgramBreaksBlockForecastTgr.Upsert(data.Key(), data)
			if err != nil {
				return err
			}
		}
	case reflect.Map, reflect.Struct:
		var internalData internalTgr
		err = json.Unmarshal(marshalBlockForecastTgrData, &internalData)
		if err != nil {
			return err
		}
		badgerProgramBreaksBlockForecastTgr := storage.Open(DbProgramBreaksBlockForecastTgr)
		data, err := internalData.Convert()
		if err != nil {
			return err
		}
		err = badgerProgramBreaksBlockForecastTgr.Upsert(data.Key(), data)
		if err != nil {
			return err
		}
	}
	return nil
}

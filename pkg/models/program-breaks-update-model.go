package models

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/badgerhold"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/storage/badger"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"reflect"
	"time"
)

type ProgramBreaksUpdateRequest struct {
	S3Key string
	Month string
}

type ProgramBreaks struct {
	Month                  *int64             `json:"Month"`
	CnlID                  *int64             `json:"CnlID"`
	ProgID                 *int64             `json:"ProgID"`
	RCID                   *int64             `json:"RCID"`
	RPID                   *int64             `json:"RPID"`
	IssTime                *int64             `json:"IssTime"`
	IssDuration            *int64             `json:"IssDuration"`
	BlockTime              *int64             `json:"BlockTime"`
	BlockNumber            *int64             `json:"BlockNumber"`
	IsPrime                *int64             `json:"IsPrime"`
	ProID                  *int64             `json:"ProID"`
	ProOriginalPTR         *int64             `json:"ProOriginalPTR"`
	BlockDistr             *int64             `json:"BlockDistr"`
	SptOptions             *int64             `json:"SptOptions"`
	TNSBlockFactTime       *int64             `json:"TNSBlockFactTime"`
	TNSBlockFactDur        *int64             `json:"TNSBlockFactDur"`
	Pro2                   *int64             `json:"Pro2"`
	BlkAdvertTypePTR       *int64             `json:"BlkAdvertTypePTR"`
	WeekDay                *int64             `json:"WeekDay"`
	PrgBegTimeL            *int64             `json:"PrgBegTimeL"`
	PrgEndMonthL           *int64             `json:"PrgEndMonthL"`
	PrgBegMonthL           *int64             `json:"PrgBegMonthL"`
	BlkAuc                 *int64             `json:"BlkAuc"`
	AucRate                *int64             `json:"AucRate"`
	NoRating               *int64             `json:"NoRating"`
	AvailableAuctionVolume *int64             `json:"AvailableAuctionVolume"`
	IsLocal                *int64             `json:"IsLocal"`
	RootRCID               *int64             `json:"RootRCID"`
	RootRPID               *int64             `json:"RootRPID"`
	IsSpecialProject       *int64             `json:"IsSpecialProject"`
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
	BlockDate              *int64             `json:"BlockDate"`
	Booked                 []BookedAttributes `json:"Booked"`
	BlockID                *int64             `json:"BlockID"`
	VM                     *int64             `json:"VM"`
	VR                     *int64             `json:"VR"`
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
	VM      *int64 `json:"VM"`
	VR      *int64 `json:"VR"`
}

func (a *internalAttr) Convert() (*Attributes, error) {
	attributes := &Attributes{
		BlockID: utils.Int64I(a.BlockID),
		VM:      utils.Int64I(a.VM),
		VR:      utils.Int64I(a.VR),
	}
	return attributes, nil
}

type BookedAttributes struct {
	RankID *int64 `json:"RankID"`
	VM     *int64 `json:"VM"`
	VR     *int64 `json:"VR"`
}

func (a *internalAttributes) Convert() (*BookedAttributes, error) {
	attributes := &BookedAttributes{
		RankID: utils.Int64I(a.Attributes["RankID"]),
		VM:     utils.Int64I(a.Attributes["VM"]),
		VR:     utils.Int64I(a.Attributes["VR"]),
	}
	return attributes, nil
}

type ProMaster struct {
	ProID     *int64    `json:"ProID"`
	PropName  *string   `json:"PropName"`
	PropValue *string   `json:"PropValue"`
	Timestamp time.Time `json:"Timestamp"`
}

func (p *internalP) Convert() (*ProMaster, error) {
	timestamp := time.Now()
	pro := &ProMaster{
		ProID:     utils.Int64I(p.P["ProID"]),
		PropName:  utils.StringI(p.P["PropName"]),
		PropValue: utils.StringI(p.P["PropValue"]),
		Timestamp: timestamp,
	}
	return pro, nil
}

type BlockForecast struct {
	TgrID            *int64    `json:"TgrID"`
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
		TgrID:            utils.Int64I(bb.B["TgrID"]),
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
	ID        *int64    `json:"ID"`
	Name      *string   `json:"Name"`
	Timestamp time.Time `json:"Timestamp"`
}

func (t *internalTgr) Convert() (*BlockForecastTgr, error) {
	timestamp := time.Now()
	block := &BlockForecastTgr{
		ID:        utils.Int64I(t.Tgr["ID"]),
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
	var vm *int64
	var vr *int64
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
			vm = utils.Int64I(attribute.VM)
			vr = utils.Int64I(attribute.VR)
		}
	} else {
		blockID = utils.Int64I(b.B["BlockID"])
		vm = utils.Int64I(b.B["VM"])
		vr = utils.Int64I(b.B["VR"])
	}
	items := &ProgramBreaks{
		CnlID:                  utils.Int64I(b.B["CnlID"]),
		ProgID:                 utils.Int64I(b.B["ProgID"]),
		RCID:                   utils.Int64I(b.B["RCID"]),
		RPID:                   utils.Int64I(b.B["RPID"]),
		IssTime:                utils.Int64I(b.B["IssTime"]),
		IssDuration:            utils.Int64I(b.B["IssDuration"]),
		BlockTime:              utils.Int64I(b.B["BlockTime"]),
		BlockNumber:            utils.Int64I(b.B["BlockNumber"]),
		IsPrime:                utils.Int64I(b.B["IsPrime"]),
		ProID:                  utils.Int64I(b.B["ProID"]),
		ProOriginalPTR:         utils.Int64I(b.B["ProOriginalPTR"]),
		BlockDistr:             utils.Int64I(b.B["BlockDistr"]),
		SptOptions:             utils.Int64I(b.B["SptOptions"]),
		TNSBlockFactTime:       utils.Int64I(b.B["TNSBlockFactTime"]),
		TNSBlockFactDur:        utils.Int64I(b.B["TNSBlockFactDur"]),
		Pro2:                   utils.Int64I(b.B["Pro2"]),
		BlkAdvertTypePTR:       utils.Int64I(b.B["BlkAdvertTypePTR"]),
		WeekDay:                utils.Int64I(b.B["WeekDay"]),
		PrgBegTimeL:            utils.Int64I(b.B["PrgBegTimeL"]),
		PrgEndMonthL:           utils.Int64I(b.B["PrgEndMonthL"]),
		PrgBegMonthL:           utils.Int64I(b.B["PrgBegMonthL"]),
		BlkAuc:                 utils.Int64I(b.B["BlkAuc"]),
		AucRate:                utils.Int64I(b.B["AucRate"]),
		NoRating:               utils.Int64I(b.B["NoRating"]),
		AvailableAuctionVolume: utils.Int64I(b.B["AvailableAuctionVolume"]),
		IsLocal:                utils.Int64I(b.B["IsLocal"]),
		RootRCID:               utils.Int64I(b.B["RootRCID"]),
		RootRPID:               utils.Int64I(b.B["RootRPID"]),
		IsSpecialProject:       utils.Int64I(b.B["IsSpecialProject"]),
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
		BlockDate:              utils.Int64I(b.B["BlockDate"]),
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
		defer ch.Close()
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
	badgerProgramBreaks := badger.Open(DbProgramBreaks)
	month := utils.Int64(request.Month)
	for _, dataB := range internalData {
		programBreaks, err := dataB.ConvertProgramBreaks()
		if err != nil {
			return err
		}
		programBreaks.Month = &month
		var networks []ProgramBreaks
		if programBreaks.BlockID != nil {
			err = badgerProgramBreaks.Find(&networks, badgerhold.Where("BlockID").Eq(*programBreaks.BlockID))
			if err != nil {
				return err
			}
		}
		for _, item := range networks {
			err = badgerProgramBreaks.Delete(item.Key(), item)
			if err != nil {
				return err
			}
		}
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
		badgerProgramBreaksProMaster := badger.Open(DbProgramBreaksProMaster)
		for _, dataItem := range internalData {
			data, err := dataItem.Convert()
			if err != nil {
				return err
			}
			var proMasters []ProMaster
			if data.ProID != nil {
				err = badgerProgramBreaksProMaster.Find(&proMasters, badgerhold.Where("BlockID").Eq(*data.ProID))
				if err != nil {
					return err
				}
			}
			for _, item := range proMasters {
				err = badgerProgramBreaksProMaster.Delete(item.Key(), item)
				if err != nil {
					return err
				}
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
		badgerProgramBreaksProMaster := badger.Open(DbProgramBreaksProMaster)
		data, err := internalData.Convert()
		if err != nil {
			return err
		}
		var proMasters []ProMaster
		if data.ProID != nil {
			err = badgerProgramBreaksProMaster.Find(&proMasters, badgerhold.Where("ProID").Eq(*data.ProID))
			if err != nil {
				return err
			}
		}
		for _, item := range proMasters {
			err = badgerProgramBreaksProMaster.Delete(item.Key(), item)
			if err != nil {
				return err
			}
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
		badgerProgramBreaksBlockForecast := badger.Open(DbProgramBreaksBlockForecast)
		for _, dataItem := range internalData {
			data, err := dataItem.Convert()
			if err != nil {
				return err
			}
			var blockForecasts []BlockForecast
			if data.BlockID != nil {
				err = badgerProgramBreaksBlockForecast.Find(&blockForecasts, badgerhold.Where("BlockID").Eq(*data.BlockID))
				if err != nil {
					return err
				}
			}
			for _, item := range blockForecasts {
				err = badgerProgramBreaksBlockForecast.Delete(item.Key(), item)
				if err != nil {
					return err
				}
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
		badgerProgramBreaksBlockForecast := badger.Open(DbProgramBreaksBlockForecast)
		data, err := internalData.Convert()
		if err != nil {
			return err
		}
		var blockForecasts []BlockForecast
		if data.BlockID != nil {
			err = badgerProgramBreaksBlockForecast.Find(&blockForecasts, badgerhold.Where("BlockID").Eq(*data.BlockID))
			if err != nil {
				return err
			}
		}
		for _, item := range blockForecasts {
			err = badgerProgramBreaksBlockForecast.Delete(item.Key(), item)
			if err != nil {
				return err
			}
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
		badgerProgramBreaksBlockForecastTgr := badger.Open(DbProgramBreaksBlockForecastTgr)
		for _, dataItem := range internalData {
			data, err := dataItem.Convert()
			if err != nil {
				return err
			}
			var blockForecastTrg []BlockForecastTgr
			if data.ID != nil {
				err = badgerProgramBreaksBlockForecastTgr.Find(&blockForecastTrg, badgerhold.Where("ID").Eq(*data.ID))
				if err != nil {
					return err
				}
			}
			for _, item := range blockForecastTrg {
				err = badgerProgramBreaksBlockForecastTgr.Delete(item.Key(), item)
				if err != nil {
					return err
				}
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
		badgerProgramBreaksBlockForecastTgr := badger.Open(DbProgramBreaksBlockForecastTgr)
		data, err := internalData.Convert()
		if err != nil {
			return err
		}
		var blockForecastTrg []BlockForecastTgr
		if data.ID != nil {
			err = badgerProgramBreaksBlockForecastTgr.Find(&blockForecastTrg, badgerhold.Where("ID").Eq(*data.ID))
			if err != nil {
				return err
			}
		}
		for _, item := range blockForecastTrg {
			err = badgerProgramBreaksBlockForecastTgr.Delete(item.Key(), item)
			if err != nil {
				return err
			}
		}
		err = badgerProgramBreaksBlockForecastTgr.Upsert(data.Key(), data)
		if err != nil {
			return err
		}
	}
	return nil
}

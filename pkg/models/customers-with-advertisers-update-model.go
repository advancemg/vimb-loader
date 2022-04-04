package models

import (
	"encoding/json"
	"fmt"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/storage"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"reflect"
	"time"
)

type CustomersWithAdvertisersUpdateRequest struct {
	S3Key string
}

type CustomersWithAdvertisers struct {
	Timestamp time.Time `json:"Timestamp"`
	ID        *int      `json:"ID"`
	Name      *string   `json:"Name"`
}

type CustomersWithAdvertisersData struct {
	Timestamp time.Time `json:"Timestamp"`
	CustID    *int      `json:"CustID"`
	AdvID     *int      `json:"AdvID"`
}

func (c *CustomersWithAdvertisers) Key() string {
	return fmt.Sprintf("%d", c.ID)
}

func (c *CustomersWithAdvertisersData) Key() string {
	return fmt.Sprintf("%d", c.AdvID)
}

func (item *internalItem) ConvertCustomerAdvertiser() (*CustomersWithAdvertisers, error) {
	timestamp := time.Now()
	items := &CustomersWithAdvertisers{
		Timestamp: timestamp,
		ID:        utils.IntI(item.Item["ID"]),
		Name:      utils.StringI(item.Item["Name"]),
	}
	return items, nil
}

func (item *internalItem) ConvertCustomerAdvertiserData() (*CustomersWithAdvertisersData, error) {
	timestamp := time.Now()
	items := &CustomersWithAdvertisersData{
		Timestamp: timestamp,
		CustID:    utils.IntI(item.Item["CustID"]),
		AdvID:     utils.IntI(item.Item["AdvID"]),
	}
	return items, nil
}

func CustomersWithAdvertisersStartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := CustomersWithAdvertisersUpdateQueue
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
			req := CustomersWithAdvertisersUpdateRequest{
				S3Key: bodyJson.Key,
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

func (request *CustomersWithAdvertisersUpdateRequest) Update() error {
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

func (request *CustomersWithAdvertisersUpdateRequest) loadFromFile() error {
	resp := utils.VimbResponse{FilePath: request.S3Key}
	/*Advertisers*/
	convertData, err := resp.Convert("Advertisers")
	if err != nil {
		return err
	}
	marshalData, err := json.Marshal(convertData)
	if err != nil {
		return err
	}
	switch reflect.TypeOf(convertData).Kind() {
	case reflect.Array, reflect.Slice:
		var internalData []internalItem
		err = json.Unmarshal(marshalData, &internalData)
		if err != nil {
			return err
		}
		badgerCustomersWithAdvertisers := storage.Open(DbCustomersWithAdvertisers)
		for _, dataItem := range internalData {
			data, err := dataItem.ConvertCustomerAdvertiser()
			if err != nil {
				return err
			}
			err = badgerCustomersWithAdvertisers.Upsert(data.Key(), data)
			if err != nil {
				return err
			}
		}
	case reflect.Map, reflect.Struct:
		var internalData internalItem
		err = json.Unmarshal(marshalData, &internalData)
		if err != nil {
			return err
		}
		badgerCustomersWithAdvertisers := storage.Open(DbCustomersWithAdvertisers)
		data, err := internalData.ConvertCustomerAdvertiser()
		if err != nil {
			return err
		}
		err = badgerCustomersWithAdvertisers.Upsert(data.Key(), data)
		if err != nil {
			return err
		}
	}
	/*Data*/
	convertAdvData, err := resp.Convert("Data")
	if err != nil {
		return err
	}
	marshalAdvData, err := json.Marshal(convertAdvData)
	if err != nil {
		return err
	}
	switch reflect.TypeOf(convertAdvData).Kind() {
	case reflect.Array, reflect.Slice:
		var internalData []internalItem
		err = json.Unmarshal(marshalAdvData, &internalData)
		if err != nil {
			return err
		}
		badgerCustomersWithAdvertisersData := storage.Open(DbCustomersWithAdvertisersData)
		for _, dataItem := range internalData {
			data, err := dataItem.ConvertCustomerAdvertiserData()
			if err != nil {
				return err
			}
			err = badgerCustomersWithAdvertisersData.Upsert(data.Key(), data)
			if err != nil {
				return err
			}
		}
	case reflect.Map, reflect.Struct:
		var internalData internalItem
		err = json.Unmarshal(marshalAdvData, &internalData)
		if err != nil {
			return err
		}
		badgerCustomersWithAdvertisersData := storage.Open(DbCustomersWithAdvertisersData)
		data, err := internalData.ConvertCustomerAdvertiserData()
		if err != nil {
			return err
		}
		err = badgerCustomersWithAdvertisersData.Upsert(data.Key(), data)
		if err != nil {
			return err
		}
	}
	return nil
}

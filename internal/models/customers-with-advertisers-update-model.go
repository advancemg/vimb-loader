package models

import (
	"encoding/json"
	"fmt"
	"github.com/advancemg/vimb-loader/internal/store"
	mq_broker "github.com/advancemg/vimb-loader/pkg/mq-broker"
	"github.com/advancemg/vimb-loader/pkg/s3"
	"github.com/advancemg/vimb-loader/pkg/utils"
	"reflect"
	"time"
)

type CustomersWithAdvertisersUpdateRequest struct {
	S3Key string
}

type CustomersWithAdvertisers struct {
	Timestamp time.Time `json:"Timestamp" bson:"Timestamp"`
	ID        *int64    `json:"ID" bson:"ID"`
	Name      *string   `json:"Name" bson:"Name"`
}

type CustomersWithAdvertisersData struct {
	Timestamp time.Time `json:"Timestamp" bson:"Timestamp"`
	CustID    *int64    `json:"CustID" bson:"CustID"`
	AdvID     *int64    `json:"AdvID" bson:"AdvID"`
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
		ID:        utils.Int64I(item.Item["ID"]),
		Name:      utils.StringI(item.Item["Name"]),
	}
	return items, nil
}

func (item *internalItem) ConvertCustomerAdvertiserData() (*CustomersWithAdvertisersData, error) {
	timestamp := time.Now()
	items := &CustomersWithAdvertisersData{
		Timestamp: timestamp,
		CustID:    utils.Int64I(item.Item["CustID"]),
		AdvID:     utils.Int64I(item.Item["AdvID"]),
	}
	return items, nil
}

func CustomersWithAdvertisersStartJob() chan error {
	errorCh := make(chan error)
	go func() {
		qName := CustomersWithAdvertisersUpdateQueue
		amqpConfig := mq_broker.InitConfig()
		defer amqpConfig.Close()
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
		db, table := utils.SplitDbAndTable(DbCustomersWithAdvertisers)
		dbCustomersWithAdvertisers, err := store.OpenDb(db, table)
		if err != nil {
			return err
		}
		for _, dataItem := range internalData {
			data, err := dataItem.ConvertCustomerAdvertiser()
			if err != nil {
				return err
			}
			err = dbCustomersWithAdvertisers.AddOrUpdate(data.Key(), data)
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
		db, table := utils.SplitDbAndTable(DbCustomersWithAdvertisers)
		dbCustomersWithAdvertisers, err := store.OpenDb(db, table)
		if err != nil {
			return err
		}
		data, err := internalData.ConvertCustomerAdvertiser()
		if err != nil {
			return err
		}
		err = dbCustomersWithAdvertisers.AddOrUpdate(data.Key(), data)
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
		db, table := utils.SplitDbAndTable(DbCustomersWithAdvertisersData)
		dbCustomersWithAdvertisersData, err := store.OpenDb(db, table)
		if err != nil {
			return err
		}
		for _, dataItem := range internalData {
			data, err := dataItem.ConvertCustomerAdvertiserData()
			if err != nil {
				return err
			}
			err = dbCustomersWithAdvertisersData.AddOrUpdate(data.Key(), data)
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
		db, table := utils.SplitDbAndTable(DbCustomersWithAdvertisersData)
		dbCustomersWithAdvertisersData, err := store.OpenDb(db, table)
		if err != nil {
			return err
		}
		data, err := internalData.ConvertCustomerAdvertiserData()
		if err != nil {
			return err
		}
		err = dbCustomersWithAdvertisersData.AddOrUpdate(data.Key(), data)
		if err != nil {
			return err
		}
	}
	return nil
}

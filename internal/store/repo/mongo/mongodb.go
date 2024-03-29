package mongo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/advancemg/vimb-loader/pkg/logging/zap"
	mongodb_client "github.com/advancemg/vimb-loader/pkg/storage/mongodb-client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
	"time"
)

const service = "mongodb"

type DbRepo struct {
	*mongo.Client
	table    string
	database string
}

func New(table, database string) (*DbRepo, error) {
	cfg := mongodb_client.InitConfig()
	cfg.DB = database
	client, err := cfg.New()
	if err != nil {
		return nil, fmt.Errorf("mongodb New(): %w", err)
	}
	return &DbRepo{client, table, database}, nil
}

type MongoKeyValue struct {
	Key   interface{}
	Value interface{}
}

func (c *DbRepo) Close() error {
	if c.Client == nil {
		return nil
	}
	err := c.Client.Disconnect(context.TODO())
	if err != nil {
		return fmt.Errorf("mongodb disconnect: %w", err)
	}
	log.PrintLog("vimb-loader", service, "info", map[string]string{"MongoDB": "Connection closed."})
	return nil
}

func (c *DbRepo) AddOrUpdate(key interface{}, data interface{}) error {
	list := []MongoKeyValue{{
		Key:   bson.M{"_id": key},
		Value: bson.M{"$set": data},
	}}
	_, err := c.AddOrUpdateMany(list, true)
	if err != nil {
		return err
	}
	return nil
}

type Timeout struct {
	IsTimeout bool          `json:"id" bson:"_id"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	Ttl       time.Duration `json:"ttl" bson:"ttl"`
}

func (c *DbRepo) Get(key interface{}, result interface{}) error {
	if c.Client == nil {
		err := c.connect()
		if err != nil {
			return fmt.Errorf("mongodb get: %w", err)
		}
		defer c.Close()
	}
	k := bson.M{key.(string): bson.M{"$gt": false}}
	log.PrintLog("vimb-loader", service, "info", map[string]string{"Get": "start"})
	var timeout Timeout
	err := c.Client.Database(c.database).Collection(c.table).FindOne(context.TODO(), k).Decode(&timeout)
	if err != nil {
		return fmt.Errorf("mongodb - Get: %w", err)
	}
	currentTime := time.Now().UTC()
	if currentTime.Sub(timeout.CreatedAt).Seconds() > timeout.Ttl.Seconds() {
		_, err = c.DeleteOne(k)
		if err != nil {
			return err
		}
		timeout.IsTimeout = false
	}
	marshal, err := bson.Marshal(timeout)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(marshal, result)
	if err != nil {
		return err
	}
	return nil
}

func (c *DbRepo) AddWithTTL(key, value interface{}, ttl time.Duration) error {
	ctx := context.Background()
	if c.Client == nil {
		err := c.connect()
		if err != nil {
			return fmt.Errorf("mongodb addwithttl: %w", err)
		}
		defer c.Close()
	}
	marshal, err := bson.Marshal(value)
	if err != nil {
		return err
	}
	var timeout Timeout
	err = bson.Unmarshal(marshal, &timeout)
	if err != nil {
		return err
	}
	timeout.Ttl = ttl
	timeout.CreatedAt = time.Now().UTC()
	log.PrintLog("vimb-loader", service, "info", map[string]string{"AddWithTTL": "start"})
	sessionCollection := c.Client.Database(c.database).Collection(c.table)
	for {
		_, err = sessionCollection.InsertOne(ctx, bson.M{key.(string): timeout.IsTimeout, "created_at": timeout.CreatedAt, "ttl": timeout.Ttl})
		if err != nil {
			if strings.Contains(err.Error(), "E11000 duplicate key error collection") {
				_, err = sessionCollection.DeleteOne(ctx, bson.M{key.(string): timeout.IsTimeout})
				if err != nil {
					return err
				}
				continue
			} else {
				log.PrintLog("vimb-loader", service, "error", map[string]string{"MongoDBClient InsertOne": "error", "error": err.Error()})
				return err
			}
		}
		break
	}
	log.PrintLog("vimb-loader", service, "info", map[string]string{"MongoDBClient InsertOne": "OK"})
	return nil
}

func (c *DbRepo) Delete(key interface{}, dataType interface{}) error {
	k := bson.M{key.(string): bson.M{"$gt": -1}}
	_ = dataType
	_, err := c.DeleteOne(k)
	if err != nil {
		return fmt.Errorf("mongodb - Delete: %w", err)
	}
	return nil
}

func (c *DbRepo) FindJson(result interface{}, filter []byte) error {
	var request map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader(filter))
	decoder.UseNumber()
	if err := decoder.Decode(&request); err != nil {
		return err
	}
	filterNetworks, err := HandleBadgerRequest(request)
	if err != nil {
		return fmt.Errorf("FindJson: %w", err)
	}
	return c.Find(result, filterNetworks)
}

func (c *DbRepo) FindWhereEq(result interface{}, filed string, value interface{}) error {
	filter := bson.M{filed: bson.M{"$eq": value}}
	return c.Find(result, filter)
}

func (c *DbRepo) FindWhereNe(result interface{}, filed string, value interface{}) error {
	filter := bson.M{filed: bson.M{"$ne": value}}
	return c.Find(result, filter)
}

func (c *DbRepo) FindWhereGt(result interface{}, filed string, value interface{}) error {
	filter := bson.M{filed: bson.M{"$gt": value}}
	return c.Find(result, filter)
}

func (c *DbRepo) FindWhereLt(result interface{}, filed string, value interface{}) error {
	filter := bson.M{filed: bson.M{"$lt": value}}
	return c.Find(result, filter)
}

func (c *DbRepo) FindWhereGe(result interface{}, filed string, value interface{}) error {
	filter := bson.M{filed: bson.M{"$gte": value}}
	return c.Find(result, filter)
}

func (c *DbRepo) FindWhereLe(result interface{}, filed string, value interface{}) error {
	filter := bson.M{filed: bson.M{"$lte": value}}
	return c.Find(result, filter)
}

func (c *DbRepo) FindWhereAnd2Eq(result interface{}, filed1 string, value1 interface{}, filed2 string, value2 interface{}) error {
	filter := bson.M{"$and": []bson.M{
		{filed1: bson.M{"$eq": value1}},
		{filed2: bson.M{"$eq": value2}},
	}}
	return c.Find(result, filter)
}

func (c *DbRepo) FindWhereAnd4Eq(result interface{}, filed1 string, value1 interface{}, filed2 string, value2 interface{}, filed3 string, value3 interface{}, filed4 string, value4 interface{}) error {
	filter := bson.M{"$and": []bson.M{
		{filed1: bson.M{"$eq": value1}},
		{filed2: bson.M{"$eq": value2}},
		{filed3: bson.M{"$eq": value3}},
		{filed4: bson.M{"$eq": value4}},
	}}
	return c.Find(result, filter)
}

func (c *DbRepo) Find(result interface{}, filter bson.M) error {
	if c.Client == nil {
		err := c.connect()
		if err != nil {
			return fmt.Errorf("mongodb find: %w", err)
		}
		defer c.Close()
	}
	log.PrintLog("vimb-loader", service, "info", map[string]string{"Find": "start"})
	coll := c.Client.Database(c.database).Collection(c.table)
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"Find": "error", "error": err.Error()})
		return err
	}
	defer cursor.Close(context.TODO())
	err = cursor.All(context.TODO(), result)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"Find": "error", "error": err.Error()})
		return err
	}
	return nil
}

func (c *DbRepo) DeleteOne(document interface{}) (int64, error) {
	if c.Client == nil {
		err := c.connect()
		if err != nil {
			return 0, fmt.Errorf("mongodb delete one: %w", err)
		}
		defer c.Close()
	}
	log.PrintLog("vimb-loader", service, "info", map[string]string{"DeleteOne": "start"})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	count, err := c.Client.Database(c.database).Collection(c.table).DeleteOne(ctx, document)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"DeleteOne": "error", "error": err.Error()})
		return 0, err
	}
	return count.DeletedCount, err
}

func (c *DbRepo) AddOrUpdateMany(list []MongoKeyValue, upsert bool) ([]byte, error) {
	if c.Client == nil {
		err := c.connect()
		if err != nil {
			return nil, fmt.Errorf("mongodb addorupdate: %w", err)
		}
		defer c.Close()
	}
	log.PrintLog("vimb-loader", service, "info", map[string]string{"AddOrUpdateMany": "start"})
	var models []mongo.WriteModel
	for _, item := range list {
		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(item.Key).
			SetUpdate(item.Value).
			SetUpsert(upsert))
	}
	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(true)
	data, err := c.Client.Database(c.database).
		Collection(c.table).
		BulkWrite(context.TODO(), models, &bulkOption)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"MongoDBClient bulkOption": "error", "error": err.Error()})
		return nil, err
	}
	log.PrintLog("vimb-loader", service, "info", map[string]string{"MongoDBClient InsertOne": "OK"})
	return bson.Marshal(data)
}

func (c *DbRepo) connect() error {
	cfg := mongodb_client.InitConfig()
	cfg.DB = c.database
	client, err := cfg.New()
	if err != nil {
		return fmt.Errorf("mongodb connect: %w", err)
	}
	c.Client = client
	return nil
}

func HandleBadgerRequest(request map[string]interface{}) (bson.M, error) {
	var filters []bson.M
	for field, value := range request {
		for key, val := range value.(map[string]interface{}) {
			filter, err := makeFilter(key, field, val)
			if err != nil {
				return nil, fmt.Errorf("HandleBadgerRequest: %w", err)
			}
			filters = append(filters, filter)
		}
	}
	return bson.M{"$and": filters}, nil
}

func makeFilter(key, filed string, value interface{}) (bson.M, error) {
	value, err := jsonNumber(value)
	if err != nil {
		return nil, fmt.Errorf("makeFilter: %w", err)
	}
	var filter bson.M
	switch key {
	case "eq":
		filter = bson.M{filed: bson.M{"$eq": value}}
	case "ne":
		filter = bson.M{filed: bson.M{"$ne": value}}
	case "gt":
		filter = bson.M{filed: bson.M{"$gt": value}}
	case "lt":
		filter = bson.M{filed: bson.M{"$lt": value}}
	case "ge":
		filter = bson.M{filed: bson.M{"$gte": value}}
	case "le":
		filter = bson.M{filed: bson.M{"$lte": value}}
	case "in":
		filter = bson.M{filed: bson.M{"$in": value}}
	case "isnil":
		filter = bson.M{filed: bson.M{"$nin": value}}
	}
	return filter, nil
}

func jsonNumber(value interface{}) (interface{}, error) {
	if number, ok := value.(json.Number); ok {
		strconv.ParseInt(string(number), 10, 64)
		dot := strings.Contains(number.String(), ".")
		if dot {
			i, err := number.Float64()
			if err != nil {
				return nil, fmt.Errorf("jsonNumber: %w", err)
			}
			return i, nil
		} else {
			i, err := number.Int64()
			if err != nil {
				return nil, fmt.Errorf("jsonNumber: %w", err)
			}
			return i, nil
		}
	}
	return value, nil
}

package mongo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/advancemg/vimb-loader/pkg/logging"
	mongodb_client "github.com/advancemg/vimb-loader/pkg/storage/mongodb-client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
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

func New(db *mongo.Client, table, database string) *DbRepo {
	return &DbRepo{db, database, table}
}

type MongoKeyValue struct {
	Key   interface{}
	Value interface{}
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

func (c *DbRepo) Get(key interface{}, result interface{}) error {
	k := bson.M{key.(string): bson.M{"$gt": -1}}
	c.connect()
	log.PrintLog("vimb-loader", service, "info", map[string]string{"Get": "start"})
	defer c.Client.Disconnect(context.Background())
	col := c.Client.Database(c.database).Collection(c.table)
	err := col.FindOne(context.TODO(), k).Decode(result)
	if err != nil {
		return fmt.Errorf("mongodb - Get: %w", err)
	}
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
	filterNetworks := HandleBadgerRequest(request)
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
	c.connect()
	log.PrintLog("vimb-loader", service, "info", map[string]string{"Find": "start"})
	defer c.Client.Disconnect(context.Background())
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
	c.connect()
	log.PrintLog("vimb-loader", service, "info", map[string]string{"DeleteOne": "start"})
	defer c.Client.Disconnect(context.Background())
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
	c.connect()
	log.PrintLog("vimb-loader", service, "info", map[string]string{"AddOrUpdateMany": "start"})
	defer c.Client.Disconnect(context.Background())
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

func (c *DbRepo) connect() {
	cfg := loadConfig()
	client, err := mongodb_client.New(
		cfg.Mongo.Host,
		cfg.Mongo.Port,
		c.database,
		cfg.Mongo.Username,
		cfg.Mongo.Password)
	if err != nil {
		panic(err)
	}
	c.Client = client
}

type config struct {
	Mongo    mongodb_client.Config `json:"mongodb"`
	Database string                `json:"database"`
}

func loadConfig() *config {
	var cfg config
	configFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&cfg)
	return &cfg
}

func HandleBadgerRequest(request map[string]interface{}) bson.M {
	var filters []bson.M
	for field, value := range request {
		for key, val := range value.(map[string]interface{}) {
			filters = append(filters, makeFilter(key, field, val))
		}
	}
	return bson.M{"$and": filters}
}

func makeFilter(key, filed string, value interface{}) bson.M {
	value = jsonNumber(value)
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
	return filter
}

func jsonNumber(value interface{}) interface{} {
	if number, ok := value.(json.Number); ok {
		strconv.ParseInt(string(number), 10, 64)
		dot := strings.Contains(number.String(), ".")
		if dot {
			i, err := number.Float64()
			if err != nil {
				panic(err)
			}
			return i
		} else {
			i, err := number.Int64()
			if err != nil {
				panic(err)
			}
			return i
		}
	}
	return value
}

package mongo

import (
	"context"
	"encoding/json"
	log "github.com/advancemg/vimb-loader/pkg/logging"
	"github.com/advancemg/vimb-loader/pkg/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DbRepo struct {
	*mongodb.Mongo
}

func (c *DbRepo) Find(result interface{}, filter interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *DbRepo) FindJson(result interface{}, filter []byte) error {
	//TODO implement me
	panic("implement me")
}

func (c *DbRepo) AddOrUpdate(key interface{}, data interface{}) error {
	c.AddOrUpdateMany()
	//func (cfg *MongoDB) AddOrUpdateMany(database, table string, list []mongo_client.MongoKeyValue, upsert bool) ([]byte, error) {
	bytes, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}
	data, err := cfg.Config.AddOrUpdateMany(database, table, bytes, upsert)
	if err != nil {
		return nil, err
	}
	return data, nil
	//}
}

func (c *DbRepo) Get(key interface{}, result interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *DbRepo) Delete(key interface{}, dataType interface{}) error {
	//TODO implement me
	panic("implement me")
}

func New(db *mongodb.Mongo) *DbRepo {
	return &DbRepo{db}
}

const service = "mongodb"

type MongoKeyValue struct {
	Key   interface{}
	Value interface{}
}

func (c *DbRepo) Exist(database, table string, document []byte) (bool, error) {
	var doc interface{}
	err := json.Unmarshal(document, &doc)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"MongodbClient connect": "error", "error": err.Error()})
		return false, err
	}
	log.PrintLog("vimb-loader", service, "info", map[string]string{"Exist": "start"})
	//_, err = c.connect()
	//if err != nil {
	//	log.PrintLog("vimb-loader", service, "error", map[string]string{"MongodbClient connect": "error", "error": err.Error()})
	//	return false, err
	//}
	findOptions := options.FindOne()
	one := c.Client.Database(database).
		Collection(table).
		FindOne(context.Background(), doc, findOptions)
	err = one.Err()
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"Exist": "error", "error": err.Error()})
		return false, err
	}
	return true, nil
}

func (c *DbRepo) CountDocuments(database, table string, document []byte) (int64, error) {
	var doc interface{}
	err := json.Unmarshal(document, &doc)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"Unmarshal": "error", "error": err.Error()})
		return 0, err
	}
	log.PrintLog("vimb-loader", service, "info", map[string]string{"CountDocuments": "start"})
	//_, err = c.connect()
	//if err != nil {
	//	log.PrintLog("vimb-loader", service, "error", map[string]string{"MongodbClient connect": "error", "error": err.Error()})
	//	return 0, err
	//}
	result, err := c.Client.Database(database).
		Collection(table).CountDocuments(context.Background(), doc)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"CountDocuments": "error", "error": err.Error()})
		return 0, err
	}
	return result, nil
}

func (c *DbRepo) DeleteOne(database, table string, document []byte) (int64, error) {
	var doc interface{}
	err := json.Unmarshal(document, &doc)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"Unmarshal": "error", "error": err.Error()})
		return 0, err
	}
	log.PrintLog("vimb-loader", service, "info", map[string]string{"DeleteOne": "start"})
	//_, err = c.connect()
	//if err != nil {
	//	log.PrintLog("vimb-loader", service, "error", map[string]string{"MongodbClient connect": "error", "error": err.Error()})
	//	return 0, err
	//}
	defer c.Client.Disconnect(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	count, err := c.Client.Database(database).Collection(table).DeleteOne(ctx, doc)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"DeleteOne": "error", "error": err.Error()})
		return 0, err
	}
	return count.DeletedCount, err
}
func (c *DbRepo) FindOne(database, table string, document []byte) ([]byte, error) {
	var doc interface{}
	err := json.Unmarshal(document, &doc)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"Unmarshal": "error", "error": err.Error()})
		return nil, err
	}
	log.PrintLog("vimb-loader", service, "info", map[string]string{"FindOne": "start"})
	//_, err = c.connect()
	//if err != nil {
	//	log.PrintLog("vimb-loader", service, "error", map[string]string{"MongodbClient connect": "error", "error": err.Error()})
	//	return nil, err
	//}
	defer c.Client.Disconnect(context.Background())
	findOptions := options.FindOne()
	one := c.Client.Database(database).
		Collection(table).
		FindOne(context.Background(), doc, findOptions)
	bytes, err := one.DecodeBytes()
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"FindOne": "error", "error": err.Error()})
		return nil, err
	}
	return bson.Marshal(bytes)
}
func (c *DbRepo) Listing(database, table string, limit, skip int64, filter []byte) ([][]byte, error) {
	var f interface{}
	err := json.Unmarshal(filter, &f)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"Unmarshal": "error", "error": err.Error()})
		return nil, err
	}
	log.PrintLog("vimb-loader", service, "info", map[string]string{"Listing": "start"})
	//_, err = c.connect()
	//if err != nil {
	//	log.PrintLog("vimb-loader", service, "error", map[string]string{"MongodbClient connect": "error", "error": err.Error()})
	//	return nil, err
	//}
	defer c.Client.Disconnect(context.Background())
	coll := c.Client.Database(database).Collection(table)
	opts := options.Find().SetSort(bson.D{{"_id", 1}}).SetSkip(skip).SetLimit(limit)
	cursor, err := coll.Find(context.TODO(), f, opts)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"Listing": "error", "error": err.Error()})
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"Listing": "error", "error": err.Error()})
		return nil, err
	}
	result := [][]byte{}
	for _, v := range results {
		m, err := bson.Marshal(v)
		if err != nil {
			log.PrintLog("vimb-loader", service, "error", map[string]string{"Listing": "error", "error": err.Error()})
			return nil, err
		}
		result = append(result, m)
	}
	return result, nil
}

func (c *DbRepo) ListingProjection(database, table string, limit, skip int64, filter, projection []byte) ([][]byte, error) {
	var f interface{}
	err := json.Unmarshal(filter, &f)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"Unmarshal": "error", "error": err.Error()})
		return nil, err
	}
	var p bson.D
	err = bson.Unmarshal(projection, &p)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"Unmarshal": "error", "error": err.Error()})
		return nil, err
	}
	log.PrintLog("vimb-loader", service, "info", map[string]string{"ListingProjection": "start"})
	//_, err = c.connect()
	//if err != nil {
	//	log.PrintLog("vimb-loader", service, "error", map[string]string{"MongodbClient connect": "error", "error": err.Error()})
	//	return nil, err
	//}
	defer c.Client.Disconnect(context.Background())
	coll := c.Client.Database(database).Collection(table)
	opts := options.Find().SetSort(bson.D{{"_id", 1}}).SetSkip(skip).SetLimit(limit).SetProjection(p)
	cursor, err := coll.Find(context.TODO(), f, opts)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"ListingProjection": "error", "error": err.Error()})
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"ListingProjection": "error", "error": err.Error()})
		return nil, err
	}
	var result [][]byte
	for _, v := range results {
		m, err := bson.Marshal(v)
		if err != nil {
			log.PrintLog("vimb-loader", service, "error", map[string]string{"ListingProjection": "error", "error": err.Error()})
			return nil, err
		}
		result = append(result, m)
	}
	return result, nil
}

func (c *DbRepo) ListingAggProjection(database, table string, limit, skip int64, aggStage, groupStage, projectionStage []byte) ([][]byte, error) {
	var aggregateStage bson.D
	err := bson.Unmarshal(aggStage, &aggregateStage)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"Unmarshal": "error", "error": err.Error()})
		return nil, err
	}
	var aggregateGroupStage bson.D
	err = bson.Unmarshal(groupStage, &aggregateGroupStage)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"Unmarshal": "error", "error": err.Error()})
		return nil, err
	}
	var aggregateProjectionStage bson.D
	err = bson.Unmarshal(projectionStage, &aggregateProjectionStage)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"Unmarshal": "error", "error": err.Error()})
		return nil, err
	}
	log.PrintLog("vimb-loader", service, "info", map[string]string{"ListingAggProjection": "start"})
	//_, err = c.connect()
	//if err != nil {
	//	log.PrintLog("vimb-loader", service, "error", map[string]string{"MongodbClient connect": "error", "error": err.Error()})
	//	return nil, err
	//}
	aggSort := bson.D{{"$sort", bson.D{{"_id", 1}}}}
	aggLimit := bson.D{{"$limit", limit}}
	aggSkip := bson.D{{"$skip", skip}}
	defer c.Client.Disconnect(context.Background())
	coll := c.Client.Database(database).Collection(table)
	opts := options.Aggregate()
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{
		aggregateStage,
		aggSort,
		aggregateGroupStage,
		aggregateProjectionStage,
		aggLimit,
		aggSkip,
	}, opts)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"ListingAggProjection": "error", "error": err.Error()})
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"ListingAggProjection": "error", "error": err.Error()})
		return nil, err
	}
	var result [][]byte
	for _, v := range results {
		m, err := bson.Marshal(v)
		if err != nil {
			log.PrintLog("vimb-loader", service, "error", map[string]string{"ListingAggProjection": "error", "error": err.Error()})
			return nil, err
		}
		result = append(result, m)
	}
	return result, nil
}

func (c *DbRepo) AddOrUpdateMany(database, table string, list []byte, upsert bool) ([]byte, error) {
	log.PrintLog("vimb-loader", service, "info", map[string]string{"AddOrUpdateMany": "start"})
	//_, err := c.connect()
	//if err != nil {
	//	log.PrintLog("vimb-loader", service, "error", map[string]string{"MongodbClient connect": "error", "error": err.Error()})
	//	return nil, err
	//}
	defer c.Client.Disconnect(context.Background())
	var models []mongo.WriteModel
	var l []MongoKeyValue
	err := json.Unmarshal(list, &l)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"Unmarshal": "error", "error": err.Error()})
		return nil, err
	}
	for _, item := range l {
		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(item.Key).
			SetUpdate(item.Value).
			SetUpsert(upsert))
	}
	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(true)
	data, err := c.Client.Database(database).
		Collection(table).
		BulkWrite(context.TODO(), models, &bulkOption)
	if err != nil {
		log.PrintLog("vimb-loader", service, "error", map[string]string{"MongoDBClient bulkOption": "error", "error": err.Error()})
		return nil, err
	}
	log.PrintLog("vimb-loader", service, "info", map[string]string{"MongoDBClient InsertOne": "OK"})
	return bson.Marshal(data)
}

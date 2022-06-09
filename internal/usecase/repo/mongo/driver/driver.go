package driver

//
//import (
//	"encoding/json"
//	mongo_client "github.com/advancemg/vimb-loader/pkg/storage/mongo"
//	"go.mongodb.org/mongo-driver/bson"
//)
//
//type MongoDB struct {
//	Config *mongo_client.Config
//}
//
//func New(host, port, database, user, pwd, debug string) *MongoDB {
//	config := mongo_client.InitConfig()
//	config.Host = host
//	config.DB = database
//	config.Port = port
//	config.Username = user
//	config.Password = pwd
//	config.Debug = debug
//	return &MongoDB{Config: config}
//}
//
//func (cfg *MongoDB) AddOrUpdateMany(database, table string, list []mongo_client.MongoKeyValue, upsert bool) ([]byte, error) {
//	bytes, err := json.Marshal(list)
//	if err != nil {
//		return nil, err
//	}
//	data, err := cfg.Config.AddOrUpdateMany(database, table, bytes, upsert)
//	if err != nil {
//		return nil, err
//	}
//	return data, nil
//}
//
//func (cfg *MongoDB) Listing(database, table string, limit, skip int64, filter interface{}) ([][]byte, error) {
//	bytes, err := json.Marshal(filter)
//	if err != nil {
//		return nil, err
//	}
//	data, err := cfg.Config.Listing(database, table, limit, skip, bytes)
//	if err != nil {
//		return nil, err
//	}
//	return data, nil
//}
//
//func (cfg *MongoDB) ListingProjection(database, table string, limit, skip int64, filter, projection interface{}) ([][]byte, error) {
//	bytesFilter, err := json.Marshal(filter)
//	if err != nil {
//		return nil, err
//	}
//	bytesProjection, err := bson.Marshal(projection)
//	if err != nil {
//		return nil, err
//	}
//	data, err := cfg.Config.ListingProjection(database, table, limit, skip, bytesFilter, bytesProjection)
//	if err != nil {
//		return nil, err
//	}
//	return data, nil
//}
//
//func (cfg *MongoDB) ListingAggProjection(database, table string, limit, skip int64, aggStage, groupStage, projectStage interface{}) ([][]byte, error) {
//	bytesAggStage, err := bson.Marshal(aggStage)
//	if err != nil {
//		return nil, err
//	}
//	bytesGroupStage, err := bson.Marshal(groupStage)
//	if err != nil {
//		return nil, err
//	}
//	bytesProjectStage, err := bson.Marshal(projectStage)
//	if err != nil {
//		return nil, err
//	}
//	data, err := cfg.Config.ListingAggProjection(database, table, limit, skip, bytesAggStage, bytesGroupStage, bytesProjectStage)
//	if err != nil {
//		return nil, err
//	}
//	return data, nil
//}
//
//func (cfg *MongoDB) FindOne(database, table string, document interface{}) ([]byte, error) {
//	bytesFilter, err := json.Marshal(document)
//	if err != nil {
//		return nil, err
//	}
//	data, err := cfg.Config.FindOne(database, table, bytesFilter)
//	if err != nil {
//		return nil, err
//	}
//	return data, nil
//}
//
//func (cfg *MongoDB) DeleteOne(database, table string, document interface{}) (int64, error) {
//	bytesFilter, err := json.Marshal(document)
//	if err != nil {
//		return 0, err
//	}
//	data, err := cfg.Config.DeleteOne(database, table, bytesFilter)
//	if err != nil {
//		return 0, err
//	}
//	return data, nil
//}
//
//func (cfg *MongoDB) CountDocuments(database, table string, document interface{}) (int64, error) {
//	bytesFilter, err := json.Marshal(document)
//	if err != nil {
//		return 0, err
//	}
//	count, err := cfg.Config.CountDocuments(database, table, bytesFilter)
//	if err != nil {
//		return 0, err
//	}
//	return count, nil
//}
//
//func (cfg *MongoDB) Exist(database, table string, document interface{}) (bool, error) {
//	bytes, err := json.Marshal(document)
//	if err != nil {
//		return false, err
//	}
//	exist, err := cfg.Config.Exist(database, table, bytes)
//	if err != nil {
//		return false, err
//	}
//	return exist, nil
//}

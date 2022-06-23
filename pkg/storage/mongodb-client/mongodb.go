package mongodb_client

import (
	"context"
	"fmt"
	log "github.com/advancemg/vimb-loader/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const service = "mongodb"

type Config struct {
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	DB       string `json:"Db"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func New(host, port, db, username, password string) (*mongo.Client, error) {
	return connect(host, port, db, username, password)
}

func connect(host, port, database, username, password string) (*mongo.Client, error) {
	log.PrintLog("vimb-loader", service, "info", map[string]string{"MongodbClient": "start"})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	credential := options.Credential{
		Username: username,
		Password: password,
	}
	var url = fmt.Sprintf(`mongodb://%s:%s/%s`, host, port, database)
	clientOptions := options.Client().
		SetAuth(credential).
		ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)
	return client, err
}

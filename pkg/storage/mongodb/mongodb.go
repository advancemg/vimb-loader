package mongodb

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
	Debug    string `json:"Debug"`
}

type Mongo struct {
	Client *mongo.Client
}

func New(host, port, db, username, password, debug string) (*Mongo, error) {
	return connect(host, port, db, username, password, debug)
}

func connect(host, port, db, username, password, debug string) (*Mongo, error) {
	log.PrintLog("vimb-loader", service, "info", map[string]string{"MongodbClient": "start"})
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	credential := options.Credential{
		AuthSource: db,
		Username:   username,
		Password:   password,
	}
	clientOptions := options.Client().
		//SetAuth(credential).
		ApplyURI(fmt.Sprintf(`mongodb://%s:%s/%s`, host, port, db)).
		SetReplicaSet("rs0")
	/*prod*/
	if debug == `false` {
		var url = fmt.Sprintf(`mongodb://%s:%s/%s`,
			host, port, db)
		clientOptions = options.Client().
			SetAuth(credential).
			ApplyURI(url)
	}
	client, err := mongo.Connect(ctx, clientOptions)
	return &Mongo{Client: client}, err
}

package mongodb_client

import (
	"context"
	"fmt"
	"github.com/advancemg/vimb-loader/pkg/logging/zap"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
)

func Test_connect(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	zap.Init()
	client, err := connect("localhost", "27017", "db", "root", "qwerty")
	if err != nil {
		panic(err)
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	usersCollection := client.Database("db").Collection("100500")
	user := bson.D{{"fullName", "User 1"}, {"age", 30}}
	result, err := usersCollection.InsertOne(context.Background(), user)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.InsertedID)
}

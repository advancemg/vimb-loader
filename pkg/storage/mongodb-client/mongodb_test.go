package mongodb_client

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
)

func Test_connect(t *testing.T) {
	client, err := connect("localhost", "27017", "db", "root", "qwerty")
	if err != nil {
		panic(err)
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	usersCollection := client.Database("db").Collection("budgets")
	user := bson.D{{"fullName", "User 1"}, {"age", 30}}
	result, err := usersCollection.InsertOne(context.Background(), user)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.InsertedID)
}

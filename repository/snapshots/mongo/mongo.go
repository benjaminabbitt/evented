package mongo

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Mongo struct {
	client   mongo.Client
	Database string
}

func eventPagesToInterface(pages []*evented_core.EventPage) []interface{} {
	s := make([]interface{}, len(pages))
	for i, v := range pages {
		s[i] = v
	}
	return s
}

func Get(id string) (snap *evented_core.Snapshot, err error) {

}
func Put(id string, snap *evented_core.Snapshot) (err error) {

}

func NewMongoClient(uri string, databaseName string) (client Mongo) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	return Mongo{client: *mongoClient, Database: databaseName}
}

package mongo

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Mongo struct {
	client mongo.Client
	Database string
}

func eventPagesToInterface(pages []*evented_core.EventPage) []interface{} {
	s := make([]interface{}, len(pages))
	for i, v := range pages {
		s[i] = v
	}
	return s
}

func (m *Mongo) Add(id string, evt []*evented_core.EventPage) (err error){
	interfaces := eventPagesToInterface(evt)
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	collection := m.client.Database(m.Database).Collection(id)
	collection.InsertMany(ctx, interfaces)
	return nil
}

func (m *Mongo) Get(id string) (evt []*evented_core.EventPage, err error){
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := m.client.Database(m.Database).Collection(id)
	cur, err := collection.Find(ctx, nil)
	var results []*evented_core.EventPage
	for cur.Next(ctx){
		var elem evented_core.EventPage
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil{
		log.Fatal(err)
	}

	cur.Close(ctx)
	return results, nil
}
func (m *Mongo) GetTo(id string, to uint32) (evt []*evented_core.EventPage, err error){
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := m.client.Database(m.Database).Collection(id)
	cur, err := collection.Find(ctx, bson.D{
		{"Sequence", bson.D{{"$lt", to}}},
	})
	var results []*evented_core.EventPage
	for cur.Next(ctx){
		var elem evented_core.EventPage
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil{
		log.Fatal(err)
	}

	cur.Close(ctx)
	return results, nil
}
func (m *Mongo) GetFrom(id string, from uint32) (evt []*evented_core.EventPage, err error){
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := m.client.Database(m.Database).Collection(id)
	cur, err := collection.Find(ctx, bson.D{
		{"Sequence", bson.D{{"$gt", from}}},
	})
	var results []*evented_core.EventPage
	for cur.Next(ctx){
		var elem evented_core.EventPage
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil{
		log.Fatal(err)
	}

	cur.Close(ctx)
	return results, nil
}
func (m *Mongo) GetFromTo(id string, from uint32, to uint32) (evt []*evented_core.EventPage, err error){
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := m.client.Database(m.Database).Collection(id)
	cur, err := collection.Find(ctx, bson.D{
		{"Sequence", bson.D{{"$lt", to}}},
		{"Sequence", bson.D{{"$gt", from}}},
	})
	var results []*evented_core.EventPage
	for cur.Next(ctx){
		var elem evented_core.EventPage
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil{
		log.Fatal(err)
	}

	cur.Close(ctx)
	return results, nil
}

func NewMongoClient(uri string, databaseName string)(client Mongo){
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	return Mongo{client: *mongoClient, Database:databaseName}
}
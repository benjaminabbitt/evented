package repository

import (
	"context"
	"github.com/benjaminabbitt/evented"
	mongosupport "github.com/benjaminabbitt/evented/support/mongo"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

type Processed struct {
	errh           *evented.ErrLogger
	log            *zap.SugaredLogger
	client         mongo.Client
	Database       string
	Collection     mongo.Collection
	CollectionName string
}

func (o Processed) Received(ctx context.Context, id uuid.UUID, sequence uint32) (err error) {
	idBytes, err := mongosupport.RootToMongo(id)
	record := MongoEventTrackRecord{
		MongoId:  idBytes,
		Root:     id.String(),
		Sequence: sequence,
	}
	if sequence == 0 {
		_, err := o.Collection.InsertOne(ctx, record)
		if err != nil {
			return err
		}
	} else {
		_, err := o.Collection.ReplaceOne(ctx, bson.D{{"_id", idBytes}}, record)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o Processed) LastReceived(ctx context.Context, id uuid.UUID) (sequence uint32, err error) {
	idBytes, err := mongosupport.RootToMongo(id)
	singleResult := o.Collection.FindOne(ctx, bson.D{{"_id", idBytes}})
	record := &MongoEventTrackRecord{}
	err = singleResult.Decode(record)
	if err != nil {
		return sequence, err
	}
	return record.Sequence, nil
}

type MongoEventTrackRecord struct {
	MongoId  [12]byte `bson:"_id"`
	Root     string
	Sequence uint32
}

func NewProcessedClient(uri string, databaseName string, log *zap.SugaredLogger, errh *evented.ErrLogger) (client *Processed) {
	mongoClient, err := mongo.Connect(nil, options.Client().ApplyURI(uri))
	errh.LogIfErr(err, "")
	err = mongoClient.Ping(nil, readpref.Primary())
	errh.LogIfErr(err, "")
	collection := mongoClient.Database(databaseName).Collection("processtracking")
	client = &Processed{client: *mongoClient, Database: databaseName, Collection: *collection, log: log, errh: errh}
	return client
}

package mongo

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/repository/snapshots"
	mongosupport "github.com/benjaminabbitt/evented/support/mongo"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

type SnapshotMongoRepo struct {
	log        *zap.SugaredLogger
	client     *mongo.Client
	collection *mongo.Collection
}

type snapshot struct {
	MongoId  [12]byte `bson:"_id"`
	Root     string
	Sequence uint32
	state    *any.Any
}

func coreToStorage(root uuid.UUID, snap *evented_core.Snapshot) *snapshot {
	mongoId, _ := mongosupport.RootToMongo(root)
	return &snapshot{
		MongoId:  mongoId,
		Root:     root.String(),
		Sequence: snap.Sequence,
		state:    snap.State,
	}

}

func storageToCore(storage *snapshot) (root uuid.UUID, snapshot *evented_core.Snapshot, err error) {
	root, err = uuid.Parse(storage.Root)
	if err != nil {
		return uuid.New(), nil, err
	}
	return root, &evented_core.Snapshot{
		Sequence: storage.Sequence,
		State:    storage.state,
	}, nil
}

func (o SnapshotMongoRepo) Get(ctx context.Context, root uuid.UUID) (snap *evented_core.Snapshot, err error) {
	idBytes, err := mongosupport.RootToMongo(root)
	singleResult := o.collection.FindOne(ctx, bson.D{{"_id", idBytes}})
	record := &snapshot{}
	err = singleResult.Decode(record)
	if err != nil {
		o.log.Error(err)
		return nil, err
	}
	_, coreRecord, err := storageToCore(record)
	if err != nil {
		o.log.Error(err)
		return nil, err
	}
	return coreRecord, nil
}

func (o SnapshotMongoRepo) Put(ctx context.Context, root uuid.UUID, snap *evented_core.Snapshot) (err error) {
	record := coreToStorage(root, snap)
	idBytes, err := mongosupport.RootToMongo(root)
	if snap.Sequence == 0 {
		_, err := o.collection.InsertOne(ctx, record)
		if err != nil {
			o.log.Error(err)
			return err
		}
	} else {
		_, err := o.collection.ReplaceOne(ctx, bson.D{{"_id", idBytes}}, record)
		if err != nil {
			o.log.Error(err)
			return err
		}
	}
	return nil
}

func NewSnapshotMongoRepo(uri string, databaseName string, log *zap.SugaredLogger) (client snapshots.SnapshotStorer) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Error(err)
	}
	collection := mongoClient.Database(databaseName).Collection("snapshots")
	if collection != nil {
		if err != nil {
			log.Fatal(err)
		}
		return SnapshotMongoRepo{client: mongoClient, collection: collection, log: log}
	}
	return nil
}

package repository

import (
	"context"
	"github.com/benjaminabbitt/evented"
	mongoSupport "github.com/benjaminabbitt/evented/support/mongo"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type Processed struct {
	errh           *evented.ErrLogger
	log            *zap.SugaredLogger
	client         mongoSupport.Client
	Database       string
	Collection     mongoSupport.Collection
	CollectionName string
}

func (o Processed) Received(id uuid.UUID, sequence uint32) (err error) {
	idBytes, err := o.uuidToIdBytes(id)
	record := MongoEventTrackRecord{
		MongoId:  idBytes,
		Root:     id.String(),
		Sequence: 0,
	}
	ctx := context.Background()
	if sequence == 0 {
		_, err := o.Collection.InsertOne(ctx, record)
		if err != nil {
			return err
		}
	} else {
		_, err := o.Collection.UpdateOne(ctx, bson.D{{"_id", idBytes}}, record)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o Processed) LastReceived(id uuid.UUID) (sequence uint32, err error) {
	idBytes, err := o.uuidToIdBytes(id)
	singleResult := o.Collection.FindOne(context.Background(), bson.D{{"_id", idBytes}})
	record := &MongoEventTrackRecord{}
	err = singleResult.Decode(record)
	if err != nil {
		return sequence, err
	}
	return record.Sequence, nil
}

func (o Processed) uuidToIdBytes(id uuid.UUID) (idBytes [12]byte, err error) {
	idByteSlice, err := id.MarshalBinary()
	if err != nil {
		return idBytes, err
	}
	copy(idBytes[:], idByteSlice[:12])
	return idBytes, nil
}

type MongoEventTrackRecord struct {
	MongoId  [12]byte `bson:"_id"`
	Root     string
	Sequence uint32
}

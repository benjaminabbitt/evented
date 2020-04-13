package mongo

import (
	"context"
	"encoding/binary"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/repository/events"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.uber.org/zap"
)

type EventRepoMongo struct {
	log            *zap.SugaredLogger
	client         mongo.Client
	Database       string
	Collection     *mongo.Collection
	CollectionName string
}

type mongoEvent struct {
	MongoId     [12]byte `bson:"_id"`
	Sequence    uint32
	CreatedAt   *timestamp.Timestamp
	Event       *any.Any
	Synchronous bool
	Root        string
}

func (m EventRepoMongo) pageToMEP(root uuid.UUID, page *evented_core.EventPage) (r mongoEvent) {
	mongoId := m.generateId(root, page)

	return mongoEvent{
		MongoId:     mongoId,
		Sequence:    m.getSequence(page),
		CreatedAt:   page.CreatedAt,
		Event:       page.Event,
		Synchronous: page.Synchronous,
		Root:        root.String(),
	}
}

func (m EventRepoMongo) pageToMEPWithSequence(root uuid.UUID, sequence uint32, page *evented_core.EventPage) (r mongoEvent) {
	page.Sequence = &evented_core.EventPage_Num{Num: sequence}
	return m.pageToMEP(root, page)
}

func (m EventRepoMongo) getSequence(page *evented_core.EventPage) uint32 {
	var sequence uint32
	switch s := page.Sequence.(type) {
	case *evented_core.EventPage_Num:
		sequence = s.Num
	default:
		m.log.Error("Attempted to retreive sequence from event without sequence set.  This should not happen")
	}
	return sequence
}

func (m EventRepoMongo) generateId(root uuid.UUID, page *evented_core.EventPage) [12]byte {
	var mongoId [12]byte
	rootBin, _ := root.MarshalBinary()
	for i, v := range rootBin[0:7] {
		mongoId[i] = v
	}

	var sequenceBytes [4]byte
	binary.BigEndian.PutUint32(sequenceBytes[:], m.getSequence(page))

	for i, v := range sequenceBytes {
		mongoId[i+8] = v
	}
	return mongoId
}

func (EventRepoMongo) mepToPage(m mongoEvent) (root uuid.UUID, page *evented_core.EventPage) {
	page = &evented_core.EventPage{
		Sequence:    &evented_core.EventPage_Num{Num: m.Sequence},
		CreatedAt:   m.CreatedAt,
		Event:       m.Event,
		Synchronous: m.Synchronous,
	}
	root, _ = uuid.Parse(m.Root)
	return root, page
}

func (m EventRepoMongo) eventPagesToInterface(root uuid.UUID, pages []*evented_core.EventPage) []interface{} {
	s := make([]interface{}, len(pages))
	for k, v := range pages {
		s[k] = m.pageToMEP(root, v)
	}
	return s
}

//Adds an array of events to the data store
func (m EventRepoMongo) Add(ctx context.Context, id uuid.UUID, events []*evented_core.EventPage) (err error) {
	var numbered []*evented_core.EventPage
	var forced *evented_core.EventPage
	remainingEvents := events
	for {
		numbered, forced, remainingEvents = m.extractUntilFirstForced(remainingEvents)
		if len(numbered) > 0 {
			err := m.insert(ctx, id, numbered)
			if err != nil {
				return err
			}
		}
		if forced != nil {
			err := m.insertForced(ctx, id, forced)
			if err != nil {
				return err
			}
		}
		if len(remainingEvents) == 0 {
			break
		}
	}
	return nil
}

func (m EventRepoMongo) extractUntilFirstForced(events []*evented_core.EventPage) (numbered []*evented_core.EventPage, forced *evented_core.EventPage, remainder []*evented_core.EventPage) {
	for idx, page := range events {
		switch page.GetSequence().(type) {
		case *evented_core.EventPage_Force:
			return events[:idx], page, events[idx+1:]
		}
	}
	return events, nil, nil
}

func (m EventRepoMongo) insertForced(ctx context.Context, id uuid.UUID, event *evented_core.EventPage) error {
	var err error
	for {
		var seq uint32
		seq, err = m.getNextSequence(ctx, id)
		if err != nil {
			return err
		}
		mep := m.pageToMEPWithSequence(id, seq, event)
		_, err = m.Collection.InsertOne(ctx, mep)
		if err != nil {
			if Any(err.(mongo.BulkWriteException).WriteErrors, isKeyConflict) {
				continue
			} else {
				return err
			}
		} else {
			break
		}
	}
	return nil
}

func isKeyConflict(err mongo.BulkWriteError) bool {
	return err.Code == 11000
}

func Any(vs []mongo.BulkWriteError, f func(writeError mongo.BulkWriteError) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

func (m EventRepoMongo) insert(ctx context.Context, id uuid.UUID, events []*evented_core.EventPage) error {
	_, err := m.Collection.InsertMany(ctx, m.eventPagesToInterface(id, events))
	return err
}

func (m EventRepoMongo) getNextSequence(ctx context.Context, id uuid.UUID) (uint32, error) {
	idStr := id.String()
	opts := options.FindOne()
	opts.SetSort(bson.D{{"sequence", -1}})
	result := m.Collection.FindOne(ctx, bson.D{{"root", idStr}}, opts)
	if result.Err() != nil {
		err := result.Err()
		//XXX: find some better way to do this
		if err.Error() == "mongo: no documents in result" {
			return 0, nil
		} else {
			return 0, result.Err()
		}
	}
	var resultModel mongoEvent
	err := result.Decode(&resultModel)
	if err != nil {
		return 0, err
	}
	return resultModel.Sequence + 1, nil
}

// Gets the events related to the provided ID
func (m EventRepoMongo) Get(ctx context.Context, id uuid.UUID) (evt []*evented_core.EventPage, err error) {
	cur, err := m.Collection.Find(ctx, bson.D{{"root", id.String()}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var results []*evented_core.EventPage
	for cur.Next(ctx) {
		var elem mongoEvent
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		_, page := m.mepToPage(elem)
		results = append(results, page)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// Gets the events related to the provided ID
// To provides an inclusive limit to the events fetched
func (m EventRepoMongo) GetTo(ctx context.Context, id uuid.UUID, to uint32) (evt []*evented_core.EventPage, err error) {
	cur, err := m.Collection.Find(ctx, bson.D{
		{"root", id.String()},
		{"sequence", bson.D{{"$lte", to}}},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var results []*evented_core.EventPage
	for cur.Next(ctx) {
		var elem mongoEvent
		err := cur.Decode(&elem)
		if err != nil {
			m.log.Fatal(err)
		}
		_, page := m.mepToPage(elem)
		results = append(results, page)
	}

	if err := cur.Err(); err != nil {
		m.log.Fatal(err)
	}

	return results, nil
}

// Gets the events related to the provided ID
// From provides an inclusive limit to the events fetched
func (m EventRepoMongo) GetFrom(ctx context.Context, id uuid.UUID, from uint32) (evt []*evented_core.EventPage, err error) {
	cur, err := m.Collection.Find(ctx, bson.D{
		{"root", id.String()},
		{"sequence", bson.D{{"$gte", from}}},
	})
	if err != nil {
		return nil, err
	}
	var results []*evented_core.EventPage
	for cur.Next(ctx) {
		var elem mongoEvent
		err := cur.Decode(&elem)
		if err != nil {
			m.log.Fatal(err)
		}
		_, page := m.mepToPage(elem)
		results = append(results, page)
	}

	if err := cur.Err(); err != nil {
		m.log.Fatal(err)
	}

	cur.Close(ctx)
	return results, nil
}

// Gets the events related to the provided ID
// From and To provide an inclusive limit to the events fetched
func (m EventRepoMongo) GetFromTo(ctx context.Context, id uuid.UUID, from uint32, to uint32) (evt []*evented_core.EventPage, err error) {
	cur, err := m.Collection.Find(ctx, bson.D{
		{"root", id.String()},
		{"sequence", bson.D{{"$lte", to}}},
		{"sequence", bson.D{{"$gte", from}}},
	})
	if err != nil {
		return nil, err
	}
	var results []*evented_core.EventPage
	for cur.Next(ctx) {
		var elem mongoEvent
		err := cur.Decode(&elem)
		if err != nil {
			m.log.Fatal(err)
		}
		_, page := m.mepToPage(elem)
		results = append(results, page)
	}

	if err := cur.Err(); err != nil {
		m.log.Fatal(err)
	}

	cur.Close(ctx)
	return results, nil
}

func (m EventRepoMongo) establishIndices() error {
	sequenceModel := mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "root", Value: bsonx.Int32(1)},
			{Key: "sequence", Value: bsonx.Int32(1)},
		},
	}
	indices := m.Collection.Indexes()
	_, err := indices.CreateOne(context.Background(), sequenceModel)
	if err != nil {
		return err
	}
	return nil
}

func NewEventRepoMongo(uri string, databaseName string, eventCollectionName string, log *zap.SugaredLogger) (client events.EventRepository, err error) {
	mongoClient, err := mongo.Connect(nil, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	err = mongoClient.Ping(nil, readpref.Primary())
	collection := mongoClient.Database(databaseName).Collection(eventCollectionName)
	client = &EventRepoMongo{client: *mongoClient, Database: databaseName, Collection: collection, CollectionName: eventCollectionName, log: log}
	return client, nil
}

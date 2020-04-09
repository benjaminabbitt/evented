package mongo

import (
	"context"
	"encoding/binary"
	"github.com/benjaminabbitt/evented"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_query "github.com/benjaminabbitt/evented/proto/query"
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

type Mongo struct {
	errh           *evented.ErrLogger
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

func (m Mongo) pageToMEP(root uuid.UUID, page evented_core.EventPage) (r mongoEvent) {
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

func (m Mongo) pageToMEPWithSequence(root uuid.UUID, sequence uint32, page evented_core.EventPage) (r mongoEvent) {
	page.Sequence = &evented_core.EventPage_Num{Num: sequence}
	return m.pageToMEP(root, page)
}
func (m Mongo) getSequence(page evented_core.EventPage) uint32 {
	var sequence uint32
	switch s := page.Sequence.(type) {
	case *evented_core.EventPage_Num:
		sequence = s.Num
	default:
		m.log.Error("Attempted to retreive sequence from event without sequence set.  This should not happen")
	}
	return sequence
}

func (m Mongo) generateId(root uuid.UUID, page evented_core.EventPage) [12]byte {
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

func (Mongo) mepToPage(m mongoEvent) (root uuid.UUID, page evented_core.EventPage) {
	page = evented_core.EventPage{
		Sequence:    &evented_core.EventPage_Num{Num: m.Sequence},
		CreatedAt:   m.CreatedAt,
		Event:       m.Event,
		Synchronous: m.Synchronous,
	}
	root, _ = uuid.Parse(m.Root)
	return root, page
}

func (m Mongo) eventPagesToInterface(root uuid.UUID, pages []*evented_core.EventPage) []interface{} {
	s := make([]interface{}, len(pages))
	for k, v := range pages {
		s[k] = m.pageToMEP(root, *v)
	}
	return s
}

//Adds an array of events to the data store
func (m Mongo) Add(ctx context.Context, id uuid.UUID, events []*evented_core.EventPage) (err error) {
	var numbered []*evented_core.EventPage
	var forced *evented_core.EventPage
	remainingEvents := events
	for {
		numbered, forced, remainingEvents = m.extractUntilFirstForced(remainingEvents)
		if len(numbered) > 0 {
			m.insert(ctx, id, numbered)
		}
		if forced != nil {
			m.insertForced(ctx, id, forced)
		}
		if len(remainingEvents) == 0 {
			break
		}
	}
	return nil
}

func (m Mongo) extractUntilFirstForced(events []*evented_core.EventPage) (numbered []*evented_core.EventPage, forced *evented_core.EventPage, remainder []*evented_core.EventPage) {
	for idx, page := range events {
		switch page.GetSequence().(type) {
		case *evented_core.EventPage_Force:
			return events[:idx], page, events[idx+1:]
		}
	}
	return events, nil, nil
}

func (m Mongo) insertForced(ctx context.Context, id uuid.UUID, event *evented_core.EventPage) {
	for {
		seq := m.getNextSequence(ctx, id)
		mep := m.pageToMEPWithSequence(id, seq, *event)
		_, err := m.Collection.InsertOne(nil, mep)
		if err == nil {
			break
		}
	}
}

func (m Mongo) insert(ctx context.Context, id uuid.UUID, events []*evented_core.EventPage) {
	m.Collection.InsertMany(ctx, m.eventPagesToInterface(id, events))
}

func (m Mongo) getNextSequence(ctx context.Context, id uuid.UUID) uint32 {
	idStr := id.String()
	options := options.FindOne()
	options.SetSort(bson.D{{"sequence", -1}})
	result := m.Collection.FindOne(ctx, bson.D{{"root", idStr}}, options)
	if result.Err() != nil {
		// XXX: the only way to identify what the error is here is via a string comparison, eww.  Working with the assumption that any error here is a no documents in result.
		return 0
	}
	var resultModel mongoEvent
	result.Decode(&resultModel)
	return resultModel.Sequence + 1
}

// Gets the next available sequence for a provided ID
func (m Mongo) GetNextSequence(ctx context.Context, id uuid.UUID) (nextSequence *evented_query.NextSequence, err error) {
	nextSequence = &evented_query.NextSequence{
		Sequence: m.getNextSequence(ctx, id),
	}
	return nextSequence, nil
}

// Gets the events related to the provided ID
func (m Mongo) Get(ctx context.Context, id uuid.UUID) (evt []*evented_core.EventPage, err error) {
	cur, err := m.Collection.Find(ctx, bson.D{{"root", id.String()}})
	var results []*evented_core.EventPage
	for cur.Next(ctx) {
		var elem mongoEvent
		err := cur.Decode(&elem)
		if err != nil {
			m.log.Fatal(err)
		}
		_, page := m.mepToPage(elem)
		results = append(results, &page)
	}

	if err := cur.Err(); err != nil {
		m.log.Fatal(err)
	}

	cur.Close(ctx)
	return results, nil
}

// Gets the events related to the provided ID
// To provides an inclusive limit to the events fetched
func (m Mongo) GetTo(ctx context.Context, id uuid.UUID, to uint32) (evt []*evented_core.EventPage, err error) {
	cur, err := m.Collection.Find(ctx, bson.D{
		{"root", id.String()},
		{"sequence", bson.D{{"$lte", to}}},
	})
	var results []*evented_core.EventPage
	for cur.Next(ctx) {
		var elem mongoEvent
		err := cur.Decode(&elem)
		if err != nil {
			m.log.Fatal(err)
		}
		_, page := m.mepToPage(elem)
		results = append(results, &page)
	}

	if err := cur.Err(); err != nil {
		m.log.Fatal(err)
	}

	cur.Close(ctx)
	return results, nil
}

// Gets the events related to the provided ID
// From provides an inclusive limit to the events fetched
func (m Mongo) GetFrom(ctx context.Context, id uuid.UUID, from uint32) (evt []*evented_core.EventPage, err error) {
	cur, err := m.Collection.Find(ctx, bson.D{
		{"root", id.String()},
		{"sequence", bson.D{{"$gte", from}}},
	})
	var results []*evented_core.EventPage
	for cur.Next(ctx) {
		var elem mongoEvent
		err := cur.Decode(&elem)
		if err != nil {
			m.log.Fatal(err)
		}
		_, page := m.mepToPage(elem)
		results = append(results, &page)
	}

	if err := cur.Err(); err != nil {
		m.log.Fatal(err)
	}

	cur.Close(ctx)
	return results, nil
}

// Gets the events related to the provided ID
// From and To provide an inclusive limit to the events fetched
func (m Mongo) GetFromTo(ctx context.Context, id uuid.UUID, from uint32, to uint32) (evt []*evented_core.EventPage, err error) {
	cur, err := m.Collection.Find(ctx, bson.D{
		{"root", id.String()},
		{"sequence", bson.D{{"$lte", to}}},
		{"sequence", bson.D{{"$gte", from}}},
	})
	var results []*evented_core.EventPage
	for cur.Next(ctx) {
		var elem mongoEvent
		err := cur.Decode(&elem)
		if err != nil {
			m.log.Fatal(err)
		}
		_, page := m.mepToPage(elem)
		results = append(results, &page)
	}

	if err := cur.Err(); err != nil {
		m.log.Fatal(err)
	}

	cur.Close(ctx)
	return results, nil
}

func (m Mongo) establishIndices() {
	sequenceModel := mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "root", Value: bsonx.Int32(1)},
			{Key: "sequence", Value: bsonx.Int32(1)},
		},
	}
	indices := m.Collection.Indexes()
	indices.CreateOne(context.Background(), sequenceModel)
}

func NewMongoClient(uri string, databaseName string, eventCollectionName string, log *zap.SugaredLogger, errh *evented.ErrLogger) (client *Mongo) {
	mongoClient, err := mongo.Connect(nil, options.Client().ApplyURI(uri))
	errh.LogIfErr(err, "")
	err = mongoClient.Ping(nil, readpref.Primary())
	errh.LogIfErr(err, "")
	collection := mongoClient.Database(databaseName).Collection(eventCollectionName)
	client = &Mongo{client: *mongoClient, Database: databaseName, Collection: collection, CollectionName: eventCollectionName, log: log, errh: errh}
	client.establishIndices()
	return client
}

package repository

import (
	"github.com/benjaminabbitt/evented"
	"github.com/benjaminabbitt/evented/support"
	mongoSupport "github.com/benjaminabbitt/evented/support/mongo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"testing"
)

func TestLastProcessedRepoSuite(t *testing.T) {
	suite.Run(t, new(LastProcessedRepo))
}

type LastProcessedRepo struct {
	suite.Suite
	log             *zap.SugaredLogger
	errh            *evented.ErrLogger
	client          *mongoSupport.MockMongoClient
	collection      *mongoSupport.MockCollection
	systemUnderTest Processed
	sampleId        uuid.UUID
	mongoId         [12]byte
}

func (o *LastProcessedRepo) SetupTest() {
	o.log, o.errh = support.Log()
	defer o.log.Sync()

	o.client = &mongoSupport.MockMongoClient{}
	o.collection = &mongoSupport.MockCollection{}

	o.systemUnderTest = Processed{
		errh:           o.errh,
		log:            o.log,
		client:         o.client,
		Database:       "",
		Collection:     o.collection,
		CollectionName: "",
	}

	id, _ := uuid.Parse("c5c10714-2272-4329-809c-38344e318279")
	o.sampleId = id
	o.mongoId = [12]byte{197, 193, 7, 20, 34, 114, 67, 41, 128, 156, 56, 52}
}

func (o *LastProcessedRepo) Test_Received_Sequence_Zero() {
	result := &mongo.InsertOneResult{}
	o.collection.On("InsertOne", mock.Anything, mock.Anything, mock.Anything).Return(result, nil)

	_ = o.systemUnderTest.Received(o.sampleId, 0)
	o.collection.AssertExpectations(o.T())
}

func (o *LastProcessedRepo) Test_Received_Sequence_NonZero() {
	result := &mongo.UpdateResult{}
	o.collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(result, nil)

	_ = o.systemUnderTest.Received(o.sampleId, 1)
	o.collection.AssertExpectations(o.T())
}

func (o *LastProcessedRepo) Test_LastReceived() {
	result := &mongoSupport.MockSingleResult{}
	o.collection.On("FindOne", mock.Anything, bson.D{{"_id", o.mongoId}}, mock.Anything).Return(result)

	metr := MongoEventTrackRecord{
		MongoId:  o.mongoId,
		Root:     o.sampleId.String(),
		Sequence: 0,
	}
	result.On("Decode", &MongoEventTrackRecord{}).Return(metr)
	_, _ = o.systemUnderTest.LastReceived(o.sampleId)
}

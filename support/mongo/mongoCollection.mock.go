package mongosupport

import (
	"context"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockMongoCollection struct {
	mock.Mock
}

func (o *MockMongoCollection) Clone(opts ...*options.CollectionOptions) (*mongo.Collection, error) {
	args := o.Called(opts)
	return args.Get(0).(*mongo.Collection), args.Error(1)
}

func (o *MockMongoCollection) Name() string {
	args := o.Called()
	return args.String(0)
}

func (o *MockMongoCollection) Database() *mongo.Database {
	args := o.Called()
	return args.Get(0).(*mongo.Database)
}

func (o *MockMongoCollection) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	args := o.Called(ctx, models, opts)
	return args.Get(0).(*mongo.BulkWriteResult), args.Error(1)
}

func (o *MockMongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := o.Called(ctx, document, opts)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (o *MockMongoCollection) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	args := o.Called(ctx, documents, opts)
	return args.Get(0).(*mongo.InsertManyResult), args.Error(1)
}

func (o *MockMongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func (o *MockMongoCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func (o *MockMongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := o.Called(ctx, filter, update, opts)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (o *MockMongoCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := o.Called(ctx, filter, update, opts)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (o *MockMongoCollection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	args := o.Called(ctx, filter, replacement, opts)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (o *MockMongoCollection) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	args := o.Called(ctx, pipeline, opts)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (o *MockMongoCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).(int64), args.Error(1)
}

func (o *MockMongoCollection) EstimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	args := o.Called(ctx, opts)
	return args.Get(0).(int64), args.Error(1)
}

func (o *MockMongoCollection) Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	args := o.Called(ctx, fieldName, filter, opts)
	return args.Get(0).([]interface{}), args.Error(1)
}

func (o *MockMongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (o *MockMongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (o *MockMongoCollection) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (o *MockMongoCollection) FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	args := o.Called(ctx, filter, replacement, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (o *MockMongoCollection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	args := o.Called(ctx, filter, update, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (o *MockMongoCollection) Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	args := o.Called(ctx, pipeline, opts)
	return args.Get(0).(*mongo.ChangeStream), args.Error(1)
}

func (o *MockMongoCollection) Indexes() mongo.IndexView {
	args := o.Called()
	return args.Get(0).(mongo.IndexView)
}

func (o *MockMongoCollection) Drop(ctx context.Context) error {
	args := o.Called(ctx)
	return args.Error(0)
}

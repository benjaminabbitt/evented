package mongo

import (
	"context"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockCollection struct{ mock.Mock }

func (o MockCollection) Clone(opts ...*options.CollectionOptions) (*mongo.Collection, error) {
	args := o.Called(opts)
	return args.Get(0).(*mongo.Collection), args.Error(1)
}

func (o MockCollection) Name() string {
	args := o.Called()
	return args.String(0)
}

func (o MockCollection) Database() *mongo.Database {
	args := o.Called()
	return args.Get(0).(*mongo.Database)
}

func (o MockCollection) BulkWrite(ctx context.Context, models []mongo.WriteModel,
	opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	args := o.Called(ctx, models, opts)
	return args.Get(0).(*mongo.BulkWriteResult), args.Error(1)
}

func (o MockCollection) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := o.Called(ctx, document, opts)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (o MockCollection) InsertMany(ctx context.Context, documents []interface{},
	opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	args := o.Called(ctx, documents, opts)
	return args.Get(0).(*mongo.InsertManyResult), args.Error(1)
}

func (o MockCollection) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func (o MockCollection) DeleteMany(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func (o MockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (o MockCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (o MockCollection) ReplaceOne(ctx context.Context, filter interface{},
	replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (o MockCollection) Aggregate(ctx context.Context, pipeline interface{},
	opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	args := o.Called(ctx, pipeline, opts)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (o MockCollection) CountDocuments(ctx context.Context, filter interface{},
	opts ...*options.CountOptions) (int64, error) {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).(int64), args.Error(1)
}

func (o MockCollection) EstimatedDocumentCount(ctx context.Context,
	opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	args := o.Called(ctx, opts)
	return args.Get(0).(int64), args.Error(1)
}

func (o MockCollection) Distinct(ctx context.Context, fieldName string, filter interface{},
	opts ...*options.DistinctOptions) ([]interface{}, error) {
	args := o.Called(ctx, fieldName, filter, opts)
	return args.Get(0).([]interface{}), args.Error(1)
}

func (o MockCollection) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (*mongo.Cursor, error) {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.Cursor), args.Error(1)
}

func (o MockCollection) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (o MockCollection) FindOneAndDelete(ctx context.Context, filter interface{},
	opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (o MockCollection) FindOneAndReplace(ctx context.Context, filter interface{},
	replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (o MockCollection) FindOneAndUpdate(ctx context.Context, filter interface{},
	update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func (o MockCollection) Watch(ctx context.Context, pipeline interface{},
	opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	args := o.Called(ctx, pipeline, opts)
	return args.Get(0).(*mongo.ChangeStream), args.Error(1)
}

func (o MockCollection) Indexes() mongo.IndexView {
	args := o.Called()
	return args.Get(0).(mongo.IndexView)
}

func (o MockCollection) Drop(ctx context.Context) error {
	args := o.Called(ctx)
	return args.Error(0)
}

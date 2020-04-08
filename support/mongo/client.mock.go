package mongo

import (
	"context"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MockMongoClient struct{ mock.Mock }

func (o MockMongoClient) Connect(ctx context.Context) error {
	args := o.Called(ctx)
	return args.Error(0)
}

func (o MockMongoClient) Disconnect(ctx context.Context) error {
	args := o.Called(ctx)
	return args.Error(0)
}

func (o MockMongoClient) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	args := o.Called(ctx, rp)
	return args.Error(0)
}

func (o MockMongoClient) StartSession(opts ...*options.SessionOptions) (mongo.Session, error) {
	args := o.Called(opts)
	return args.Get(0).(mongo.Session), args.Error(1)
}

func (o MockMongoClient) Database(name string, opts ...*options.DatabaseOptions) *mongo.Database {
	args := o.Called(name, opts)
	return args.Get(0).(*mongo.Database)
}

func (o MockMongoClient) ListDatabases(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) (mongo.ListDatabasesResult, error) {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).(mongo.ListDatabasesResult), args.Error(1)
}

func (o MockMongoClient) ListDatabaseNames(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) ([]string, error) {
	args := o.Called(ctx, filter, opts)
	return args.Get(0).([]string), args.Error(1)
}

func (o MockMongoClient) UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error {
	args := o.Called(ctx, fn)
	return args.Error(0)
}

func (o MockMongoClient) UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(mongo.SessionContext) error) error {
	args := o.Called(ctx, opts, fn)
	return args.Error(0)
}

func (o MockMongoClient) Watch(ctx context.Context, pipeline interface{},
	opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	args := o.Called(ctx, pipeline, opts)
	return args.Get(0).(*mongo.ChangeStream), args.Error(1)
}

func (o MockMongoClient) NumberSessionsInProgress() int {
	args := o.Called()
	return args.Get(0).(int)
}

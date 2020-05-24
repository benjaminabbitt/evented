package main

import (
	"fmt"
	"github.com/benjaminabbitt/evented/applications/commandHandler/business/client"
	"github.com/benjaminabbitt/evented/applications/commandHandler/configuration"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework/transport"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	"github.com/benjaminabbitt/evented/repository/events"
	eventmongo "github.com/benjaminabbitt/evented/repository/events/mongo"
	"github.com/benjaminabbitt/evented/repository/snapshots"
	snapshotmongo "github.com/benjaminabbitt/evented/repository/snapshots/mongo"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/benjaminabbitt/evented/transport/async/amqp/sender"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	opentracing "github.com/opentracing/opentracing-go"
	config "github.com/uber/jaeger-client-go/config"
	zap2jaeger "github.com/uber/jaeger-client-go/log/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"io"
	"net"
)

var log *zap.SugaredLogger

func main() {
	log = support.Log()

	config := configuration.Configuration{}
	config.Initialize("commandHandler", log)

	tracer, closer := setupJaeger(*config.AppName)
	initSpan := tracer.StartSpan("Init")
	defer closer.Close()

	businessAddress := config.BusinessURL()
	commandHandlerPort := config.Port()
	log.Infow("Starting Command Handler", "port", commandHandlerPort)
	businessClient, _ := client.NewBusinessClient(businessAddress, log)
	log.Infow("Command Handler Started", "port", commandHandlerPort)

	eventRepo, _ := setupEventRepo(config, log, initSpan)
	ssRepo := setupSnapshotRepo(config, initSpan)

	repo := eventBook.MakeRepositoryBasic(eventRepo, ssRepo, config.Domain(), log)

	handlers := transport.NewTransportHolder(log)

	for _, url := range config.SagaURLs() {
		log.Infow("Connecting with Saga... ", "url", url)
		sagaConn := grpcWithInterceptors.GenerateConfiguredConn(url, log, tracer)
		handlers.Add(saga.NewGRPCSagaClient(sagaConn))
		log.Infow("Connection with Saga Successful", "url", url)
	}

	for _, url := range config.ProjectorURLs() {
		log.Infow("Connecting with Projector... ", "url", url)
		projectorConn := grpcWithInterceptors.GenerateConfiguredConn(url, log, tracer)
		handlers.Add(projector.NewGRPCProjector(projectorConn))
		log.Infow("Connection with Projector Successful.", "url", url)
	}

	err := handlers.Add(setupServiceBus(config, initSpan))
	if err != nil {
		log.Error(err)
	}

	server := framework.NewServer(
		repo,
		handlers,
		businessClient,
		log,
	)

	log.Infow("Opening port", "port", config.Port())
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port()))
	if err != nil {
		log.Error(err)
	}
	log.Infow("Creating GRPC Server")
	rpc := grpcWithInterceptors.GenerateConfiguredServer(log.Desugar(), tracer)
	log.Infow("Registering Command Handler with GRPC")
	eventedcore.RegisterCommandHandlerServer(rpc, server)

	health := health.NewServer()
	grpc_health_v1.RegisterHealthServer(rpc, health)
	health.Resume()
	log.Infow("Handler registered.")
	log.Infow("Serving...")
	initSpan.Finish()
	err = rpc.Serve(lis)
	if err != nil {
		log.Error(err)
	}
}

func setupSnapshotRepo(config configuration.Configuration, span opentracing.Span) (repo snapshots.SnapshotStorer) {
	childSpan := span.Tracer().StartSpan("Snapshot Repo Initialization", opentracing.ChildOf(span.Context()))
	defer childSpan.Finish()
	return snapshotmongo.NewSnapshotMongoRepo(config.SnapshotStoreURL(), config.SnapshotStoreDatabaseName(), log)
}

func setupServiceBus(config configuration.Configuration, span opentracing.Span) (ch chan *eventedcore.EventBook) {
	childSpan := span.Tracer().StartSpan("Service Bus Initialization", opentracing.ChildOf(span.Context()))
	defer childSpan.Finish()
	ch = make(chan *eventedcore.EventBook)
	trans := sender.NewAMQPSender(ch, config.TransportURL(), config.TransportExchange(), log)
	err := trans.Connect()
	if err != nil {
		log.Error(err)
	}
	trans.Run()
	return ch
}

func setupEventRepo(config configuration.Configuration, log *zap.SugaredLogger, span opentracing.Span) (repo events.EventStorer, err error) {
	childSpan := span.Tracer().StartSpan("Event Repo Initialization", opentracing.ChildOf(span.Context()))
	defer childSpan.Finish()
	defer span.Finish()
	repo, err = eventmongo.NewEventRepoMongo(config.EventStoreURL(), config.EventStoreDatabaseName(), config.EventStoreCollectionName(), log)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func setupJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.New(service, config.Logger(zap2jaeger.NewLogger(log.Desugar())))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

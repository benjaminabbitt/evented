package main

import (
	"context"
	"fmt"
	actx2 "github.com/benjaminabbitt/evented/applications/command/command-handler/actx"
	"github.com/benjaminabbitt/evented/applications/command/command-handler/business/client"
	"github.com/benjaminabbitt/evented/applications/command/command-handler/configuration"
	"github.com/benjaminabbitt/evented/applications/command/command-handler/framework"
	"github.com/benjaminabbitt/evented/applications/command/command-handler/framework/transport"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	"github.com/benjaminabbitt/evented/repository/events"
	"github.com/benjaminabbitt/evented/repository/events/memory"
	eventmongo "github.com/benjaminabbitt/evented/repository/events/mongo"
	"github.com/benjaminabbitt/evented/repository/snapshots"
	memory2 "github.com/benjaminabbitt/evented/repository/snapshots/memory"
	snapshotmongo "github.com/benjaminabbitt/evented/repository/snapshots/mongo"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/actx"
	"github.com/benjaminabbitt/evented/support/consul"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/benjaminabbitt/evented/support/jaeger"
	"github.com/benjaminabbitt/evented/support/network"
	"github.com/benjaminabbitt/evented/transport/async/amqp/sender"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"io"
	"net"
)

var log *zap.SugaredLogger

type SyncItem struct {
	Url string
}

func main() {
	log = support.Log()
	support.LogStartup(log, "Command Handler")
	appCtx := &actx2.BasicCommandHandlerApplicationContext{
		Actx: actx.Actx{
			Log:    log,
			Tracer: nil,
		},
		RetryStrategy: backoff.NewExponentialBackOff(),
	}

	v, err := support.Initialize(log, viper.New())
	appCtx.Config = v

	if viper.GetString(support.ConfigMgmtType) == string(support.Consul) {
		setupConsul(log, v)
	}

	tracer, closer := jaeger.SetupJaeger(fmt.Sprintf("%s-%s", v.GetString(configuration.Domain), v.GetString(configuration.Name)), log)
	appCtx.SetTracer(tracer)
	initSpan := tracer.StartSpan("Init")
	defer func(closer io.Closer) {
		err := closer.Close()
		if err != nil {
			log.Error(err)
		}
	}(closer)

	businessAddress := appCtx.Config.GetString(configuration.BusinessUrl)
	businessClient, _ := client.NewBusinessClient(appCtx, businessAddress)

	commandHandlerPort := appCtx.Config.GetUint32(configuration.Port)
	log.Infow("Starting Business Logic Coordinator", "port", commandHandlerPort)

	eventRepo, err := setupEventRepo(appCtx, initSpan)
	if err != nil {
		log.Fatal(err)
	}
	log.Infow("Set up Event Repository")

	ssRepo, err := setupSnapshotRepo(appCtx, initSpan)
	if err != nil {
		log.Fatal(err)
	}
	log.Infow("Set up Snapshot Repository")

	repo := eventBook.MakeRepositoryBasic(appCtx, eventRepo, ssRepo)

	handlers := transport.NewTransportHolder(appCtx)

	log.Infow("Setting up Synchronous Sagas")
	syncCount := 0
	for {
		url := appCtx.Config.GetString(fmt.Sprintf("%s.%d.%s", configuration.SyncSaga, syncCount, configuration.Url))
		if url == "" {
			break
		}
		log.Infow("Connecting with Saga... ", "url", url)
		sagaConn := grpcWithInterceptors.GenerateConfiguredConn(url, log, tracer)
		handlers.AddSagaTransporter(saga.NewGRPCSagaClient(sagaConn))
		log.Infow("Connection with Saga Successful", "url", url)
	}
	log.Infow("Synchronous Sagas done")

	log.Infow("Setting up Synchronous Projectors")
	syncCount = 0
	for {
		url := appCtx.Config.GetString(fmt.Sprintf("%s.%d.%s", configuration.SyncProj, syncCount, configuration.Url))
		if url == "" {
			break
		}
		log.Infow("Connecting with evented... ", "url", url)
		projectorConn := grpcWithInterceptors.GenerateConfiguredConn(url, log, tracer)
		handlers.AddProjectorClient(projector.NewGRPCProjector(projectorConn))
		log.Infow("Connection with Projector Successful.", "url", url)
	}
	log.Infow("Synchronous Projectors done")

	handlers.AddEventBookChan(setupTransport(appCtx, initSpan))

	server := framework.NewServer(
		appCtx,
		repo,
		handlers,
		businessClient,
	)

	var addrs []string
	if viper.GetBool("bindLocal") {
		addrs = network.GetExternalAddrs()
	}
	log.Infow("Opening port on addresses", "port", appCtx.Config.GetUint32(configuration.Port), "addrs", addrs)
	listeners := listen(addrs, appCtx.Config.GetUint(configuration.Port))
	log.Infow("Creating GRPC Server")
	rpc := grpcWithInterceptors.GenerateConfiguredServer(log.Desugar(), tracer)
	log.Infow("Registering Command Handler with GRPC")
	evented.RegisterBusinessCoordinatorServer(rpc, server)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(rpc, healthServer)
	healthServer.Resume()
	log.Infow("Handler registered.")
	log.Infow("Serving...")
	initSpan.Finish()
	for _, listener := range listeners {
		err := rpc.Serve(listener)
		if err != nil {
			log.Error(err)
		}
	}
}

func listen(externalAddrs []string, port uint) (listeners []net.Listener) {
	if externalAddrs == nil {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			listeners = append(listeners, listener)
		} else {
			log.Error(err)
		}
	} else {
		for _, addr := range externalAddrs {
			listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
			if err == nil {
				listeners = append(listeners, listener)
			} else {
				log.Error(err)
			}
		}
	}
	return listeners
}

func setupSnapshotRepo(actx *actx2.BasicCommandHandlerApplicationContext, span opentracing.Span) (repo snapshots.SnapshotStorer, err error) {
	childSpan := actx.Tracer().StartSpan("Snapshot Repo Initialization", opentracing.ChildOf(span.Context()))
	defer childSpan.Finish()
	snapshotKind := actx.Config.GetString(configuration.SnapshotKind)
	if snapshotKind == configuration.MongoKind {
		return snapshotmongo.NewSnapshotMongoRepo(&actx.Actx, actx.Config.GetString(configuration.SnapshotMongoUrl), actx.Config.GetString(configuration.SnapshotMongoCollection), actx.Config.GetString(configuration.SnapshotMongoName))
	} else if snapshotKind == configuration.MemoryKind {
		return memory2.NewSnapshotRepoMemory(actx.Log())
	} else {
		return nil, fmt.Errorf("specified snapshot repository type %s is invalid", snapshotKind)
	}
}

func setupTransport(appCtx *actx2.BasicCommandHandlerApplicationContext, span opentracing.Span) (ch chan *evented.EventBook) {
	childSpan := span.Tracer().StartSpan("Service Bus Initialization", opentracing.ChildOf(span.Context()))
	defer childSpan.Finish()
	ch = make(chan *evented.EventBook)
	var trans sender.EventSender

	transportKind := appCtx.Config.GetString(configuration.TransportKind)
	if transportKind == configuration.TransportRabbitKind {
		trans = sender.NewAMQPSender(ch, appCtx.Config.GetString(configuration.TransportRabbitUrl), appCtx.Config.GetString(configuration.TransportRabbitExchange), log)
		err := trans.(*sender.AMQPSender).Connect()
		if err != nil {
			log.Error(err)
		}
	} else if transportKind == configuration.TransportNoOpKind {
		trans = sender.NoOp{}
	}
	trans.Run()
	return ch
}

func setupEventRepo(actx *actx2.BasicCommandHandlerApplicationContext, span opentracing.Span) (repo events.EventStorer, err error) {
	childSpan := span.Tracer().StartSpan("Event Repo Initialization", opentracing.ChildOf(span.Context()))
	defer childSpan.Finish()
	eventsKind := actx.Config.GetString(configuration.RepoKind)
	if eventsKind == configuration.MemoryKind {
		repo, err = memory.NewEventRepoMemory(log)
		log.Debug("Memory event repository initialized")
	} else if eventsKind == configuration.MongoKind {
		url := actx.Config.GetString(configuration.RepoMongoUrl)
		name := actx.Config.GetString(configuration.RepoMongoName)
		collection := actx.Config.GetString(configuration.RepoMongoCollection)
		repo, err = eventmongo.NewEventRepoMongo(context.Background(), url, name, collection, log)
		log.Debug("MongoDB event repository initialized")
	} else {
		log.Error("Specified Event Repository %s does not match one of recognized: %s", eventsKind, []string{configuration.MongoKind, configuration.MemoryKind})
	}
	if err != nil {
		log.Error(err)
	}
	return repo, nil
}

func setupConsul(log *zap.SugaredLogger, config *viper.Viper) {
	c := consul.NewEventedConsul(config.GetString(configuration.ConsulHost), config.GetUint(configuration.ConsulPort))
	id, err := uuid.NewRandom()
	if err != nil {
		log.Error(err)
	}
	err = c.Register(config.GetString(configuration.Name), id.String())
	if err != nil {
		log.Error("Error registering with Consul", err)
	}
}

package main

import (
	"context"
	"fmt"
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
	"github.com/benjaminabbitt/evented/support/network"
	"github.com/benjaminabbitt/evented/transport/async/amqp/sender"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go/config"
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
	support.LogStartup(log, "Command Handler")
	appCtx := &framework.BasicCommandHandlerApplicationContext{
		Actx: actx.Actx{
			Log:    log,
			Tracer: nil,
		},
		RetryStrategy: backoff.NewExponentialBackOff(),
	}

	conf := &configuration.Configuration{}
	conf = support.Initialize(log, conf).(*configuration.Configuration)
	appCtx.Config = conf

	setupConsul(log, conf)

	tracer, closer := setupJaeger(fmt.Sprintf("%s-%s", conf.Domain, conf.Name))
	appCtx.SetTracer(tracer)
	initSpan := tracer.StartSpan("Init")
	defer func(closer io.Closer) {
		err := closer.Close()
		if err != nil {
			log.Error(err)
		}
	}(closer)

	businessAddress := appCtx.Config.Business.Url
	businessClient, _ := client.NewBusinessClient(appCtx, businessAddress)

	commandHandlerPort := conf.Port
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
	for _, ea := range conf.Sync.Sagas {
		log.Infow("Connecting with Saga... ", "url", ea.Url)
		sagaConn := grpcWithInterceptors.GenerateConfiguredConn(ea.Url, log, tracer)
		handlers.AddSagaTransporter(saga.NewGRPCSagaClient(sagaConn))
		log.Infow("Connection with Saga Successful", "url", ea.Url)
	}
	log.Infow("Synchronous Sagas done")

	log.Infow("Setting up Synchronous Projectors")
	for _, ea := range conf.Sync.Projectors {
		log.Infow("Connecting with evented... ", "url", ea.Url)
		projectorConn := grpcWithInterceptors.GenerateConfiguredConn(ea.Url, log, tracer)
		handlers.AddProjectorClient(projector.NewGRPCProjector(projectorConn))
		log.Infow("Connection with Projector Successful.", "url", ea.Url)
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
	log.Infow("Opening port on addresses", "port", conf.Port, "addrs", addrs)
	listeners := listen(addrs, conf.Port)
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

func setupSnapshotRepo(actx *framework.BasicCommandHandlerApplicationContext, span opentracing.Span) (repo snapshots.SnapshotStorer, err error) {
	childSpan := actx.Tracer().StartSpan("Snapshot Repo Initialization", opentracing.ChildOf(span.Context()))
	defer childSpan.Finish()
	if actx.Config.Snapshots.Kind == mongodbName {
		return snapshotmongo.NewSnapshotMongoRepo(actx.Config.Snapshots.Mongodb.Url, actx.Config.Snapshots.Mongodb.Name, log)
	} else if actx.Config.Snapshots.Kind == memoryName {
		return memory2.NewSnapshotRepoMemory(actx.Log())
	} else {
		return nil, fmt.Errorf("specified snapshot repository type %s is invalid", actx.Config.Snapshots.Kind)
	}
}

const noopName = "noop"
const amqpName = "amqp"

func setupTransport(appCtx *framework.BasicCommandHandlerApplicationContext, span opentracing.Span) (ch chan *evented.EventBook) {
	childSpan := span.Tracer().StartSpan("Service Bus Initialization", opentracing.ChildOf(span.Context()))
	defer childSpan.Finish()
	ch = make(chan *evented.EventBook)
	var trans sender.EventSender
	if appCtx.Config.Transport.Kind == amqpName {
		trans = sender.NewAMQPSender(ch, appCtx.Config.Transport.Rabbitmq.Url, appCtx.Config.Transport.Rabbitmq.Exchange, log)
		err := trans.(*sender.AMQPSender).Connect()
		if err != nil {
			log.Error(err)
		}
	} else if appCtx.Config.Transport.Kind == noopName {
		trans = sender.NoOp{}
	}
	trans.Run()
	return ch
}

const memoryName = "memory"
const mongodbName = "mongodb"

func setupEventRepo(actx *framework.BasicCommandHandlerApplicationContext, span opentracing.Span) (repo events.EventStorer, err error) {
	childSpan := span.Tracer().StartSpan("Event Repo Initialization", opentracing.ChildOf(span.Context()))
	defer childSpan.Finish()
	var eventRepoTypes = []string{"memory", "mongodb"}
	if actx.Config.Events.Kind == memoryName {
		repo, err = memory.NewEventRepoMemory(log)
		log.Debug("Memory event repository initialized")
	} else if actx.Config.Events.Kind == mongodbName {
		repo, err = eventmongo.NewEventRepoMongo(context.Background(), actx.Config.Events.Mongodb.Url, actx.Config.Events.Mongodb.Name, actx.Config.Events.Mongodb.Collection, log)
		log.Debug("MongoDB event repository initialized")
	} else {
		log.Error("Specified Event Repository %s does not match one of recognized: ", eventRepoTypes)
	}
	if err != nil {
		log.Error(err)
	}
	return repo, nil
}

func setupJaeger(serviceName string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(zap2jaeger.NewLogger(log.Desugar())))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func setupConsul(log *zap.SugaredLogger, config *configuration.Configuration) {
	c := consul.NewEventedConsul(config.ConsulHost, config.Port)
	id, err := uuid.NewRandom()
	if err != nil {
		log.Error(err)
	}
	err = c.Register(config.Name, id.String())
	if err != nil {
		log.Error("Error registering with Consul", err)
	}
}

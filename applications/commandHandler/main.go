package main

import (
	"context"
	"fmt"
	"github.com/benjaminabbitt/evented/applications/commandHandler/business/client"
	"github.com/benjaminabbitt/evented/applications/commandHandler/configuration"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework/transport"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/business"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/core"

	"github.com/benjaminabbitt/evented/repository/eventBook"
	"github.com/benjaminabbitt/evented/repository/events"
	"github.com/benjaminabbitt/evented/repository/events/memory"
	eventmongo "github.com/benjaminabbitt/evented/repository/events/mongo"
	"github.com/benjaminabbitt/evented/repository/snapshots"
	snapshotmongo "github.com/benjaminabbitt/evented/repository/snapshots/mongo"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/consul"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/benjaminabbitt/evented/transport/async/amqp/sender"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
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

	conf := configuration.Configuration{}
	conf.Initialize(log)

	setupConsul(log, conf)

	tracer, closer := setupJaeger(conf.AppName())
	initSpan := tracer.StartSpan("Init")
	defer closer.Close()

	businessAddress := conf.BusinessURL()
	businessClient, _ := client.NewBusinessClient(businessAddress, log)

	commandHandlerPort := conf.Port()
	log.Infow("Starting Business Logic Coordinator", "port", commandHandlerPort)

	eventRepo, _ := setupEventRepo(conf, log, initSpan)
	ssRepo := setupSnapshotRepo(conf, initSpan)

	repo := eventBook.MakeRepositoryBasic(eventRepo, ssRepo, conf.Domain(), log)

	handlers := transport.NewTransportHolder(log)

	for _, url := range conf.SagaURLs() {
		log.Infow("Connecting with Saga... ", "url", url)
		sagaConn := grpcWithInterceptors.GenerateConfiguredConn(url, log, tracer)
		err := handlers.Add(saga.NewGRPCSagaClient(sagaConn))
		if err != nil {
			log.Error(err)
		}
		log.Infow("Connection with Saga Successful", "url", url)
	}

	for _, url := range conf.ProjectorURLs() {
		log.Infow("Connecting with Projector... ", "url", url)
		projectorConn := grpcWithInterceptors.GenerateConfiguredConn(url, log, tracer)
		err := handlers.Add(projector.NewGRPCProjector(projectorConn))
		if err != nil {
			log.Error(err)
		}
		log.Infow("Connection with Projector Successful.", "url", url)
	}

	err := handlers.Add(setupServiceBus(conf, initSpan))
	if err != nil {
		log.Error(err)
	}

	server := framework.NewServer(
		repo,
		handlers,
		businessClient,
		log,
	)

	var addrs []string
	if viper.GetBool("bindLocal") {
		addrs = getExternalAddrs()
	}
	log.Infow("Opening port on addresses", "port", conf.Port(), "addrs", addrs)
	listeners := listen(addrs, conf.Port())
	log.Infow("Creating GRPC Server")
	rpc := grpcWithInterceptors.GenerateConfiguredServer(log.Desugar(), tracer)
	log.Infow("Registering Command Handler with GRPC")
	business.RegisterBusinessCoordinatorServer(rpc, server)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(rpc, healthServer)
	healthServer.Resume()
	log.Infow("Handler registered.")
	log.Infow("Serving...")
	initSpan.Finish()
	for _, listener := range listeners {
		err = rpc.Serve(listener)
		if err != nil {
			log.Error(err)
		}
	}
}

func listen(externalAddrs []string, port uint) (listeners []net.Listener) {
	if externalAddrs == nil {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			listeners = append(listeners, listener)
		} else {
			log.Error(err)
		}
	} else {
		for _, addr := range externalAddrs {
			listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
			if err != nil {
				listeners = append(listeners, listener)
			} else {
				log.Error(err)
			}
		}
	}
	return listeners
}

func getExternalAddrs() (externalAddrs []string) {
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, i := range ifaces {
			addrs, _ := i.Addrs()
			for _, addr := range addrs {
				switch v := addr.(type) {
				case *net.IPNet:
					externalAddrs = addIfNotLoopback(v.IP, externalAddrs)
				case *net.IPAddr:
					externalAddrs = addIfNotLoopback(v.IP, externalAddrs)
				}
			}
		}
	}
	return externalAddrs
}

func addIfNotLoopback(addr net.IP, externalAddrs []string) (rExternalAddrs []string) {
	rExternalAddrs = externalAddrs
	if !addr.IsLoopback() {
		rExternalAddrs = append(rExternalAddrs, addr.String())
	}
	return rExternalAddrs
}

func setupSnapshotRepo(config configuration.Configuration, span opentracing.Span) (repo snapshots.SnapshotStorer) {
	childSpan := span.Tracer().StartSpan("Snapshot Repo Initialization", opentracing.ChildOf(span.Context()))
	defer childSpan.Finish()
	return snapshotmongo.NewSnapshotMongoRepo(config.SnapshotStoreURL(), config.SnapshotStoreDatabaseName(), log)
}

func setupServiceBus(config configuration.Configuration, span opentracing.Span) (ch chan core.EventBook) {
	childSpan := span.Tracer().StartSpan("Service Bus Initialization", opentracing.ChildOf(span.Context()))
	defer childSpan.Finish()
	ch = make(chan core.EventBook)
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
	var eventRepoTypes = []string{"Memory", "MongoDb"}
	if config.EventRepoType() == eventRepoTypes[0] {
		repo, err = memory.NewEventRepoMemory(log)
	} else if config.EventRepoType() == eventRepoTypes[1] {
		repo, err = eventmongo.NewEventRepoMongo(context.Background(), config.EventStoreURL(), config.EventStoreDatabaseName(), config.EventStoreCollectionName(), log)
	} else {
		log.Error("Specified Event Repository %s does not match one of recognized, \"MongoDb\"")
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

func setupConsul(log *zap.SugaredLogger, config configuration.Configuration) {
	c := consul.EventedConsul{Log: log, ConsulHost: config.ConsulHost()}
	id, err := uuid.NewRandom()
	if err != nil {
		log.Error(err)
	}
	err = c.Register(config.AppName(), id.String(), config.Port())
	if err != nil {
		log.Error(err)
	}
}

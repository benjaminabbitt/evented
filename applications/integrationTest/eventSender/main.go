package main

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_saga_coordinator "github.com/benjaminabbitt/evented/proto/sagaCoordinator"
	"github.com/benjaminabbitt/evented/support"
	"github.com/google/uuid"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

var log *zap.SugaredLogger

func main() {
	log := support.Log()
	defer log.Sync()

	var name *string = flag.String("appName", "", "The name of the application.  This is used in a number of places, from configuration file name, to queue names.")
	var configPath *string = flag.String("configPath", ".", "The configuration path of the application.  Full config will be located at $configpath/$appName.yaml")
	flag.Parse()

	err := support.SetupConfig(name, configPath, flag.CommandLine)
	if err != nil {
		log.Error(err)
	}

	log.Info("Starting...")
	target := viper.GetString("eventHandlerURL")
	log.Info(target)

	var customFunc grpc_zap.CodeToLevel
	zapLogOpts := []grpc_zap.Option{
		grpc_zap.WithLevels(customFunc),
	}
	grpc_zap.ReplaceGrpcLogger(log.Desugar())

	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Millisecond)),
		grpc_retry.WithCodes(codes.NotFound, codes.Aborted, codes.Unavailable, codes.Unimplemented, codes.Unknown),
	}
	conn, err := grpc.Dial(target,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(grpc_zap.UnaryClientInterceptor(log.Desugar(), zapLogOpts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
		grpc.WithStreamInterceptor(grpc_zap.StreamClientInterceptor(log.Desugar(), zapLogOpts...)),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(retryOpts...)),
	)

	log.Infof("Connected to remote %s", target)
	if err != nil {
		log.Error(err)
		stat, _ := status.FromError(err)
		log.Error(stat)
	}
	sh := evented_saga_coordinator.NewSagaCoordinatorClient(conn)
	log.Info("Client Created...")
	id, err := uuid.NewRandom()
	protoId := evented_proto.UUIDToProto(id)

	var pages []*evented_core.EventPage
	for i := 0; i <= 1; i++ {
		pages = append(pages, framework.NewEventPage(uint32(i), false, nil))
	}
	eventBook := &evented_core.EventBook{
		Cover: &evented_core.Cover{
			Domain: viper.GetString("domain"),
			Root:   &protoId,
		},
		Pages: pages,
	}
	res, err := sh.HandleSync(context.Background(), eventBook)
	log.Info(res)
	if err != nil {

	}
}

package cmd

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework"
	"github.com/benjaminabbitt/evented/applications/integrationTest/sender/configuration"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/benjaminabbitt/evented/support/jaeger"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"strconv"
)

func init() {
	sendEvent.Flags().String("host", "localhost", "The host with which to connect")
	sendEvent.Flags().Int("port", 1738, "The port on which to connect")
	rootCmd.AddCommand(sendEvent)
}

var sendEvent = &cobra.Command{
	Use:   "event",
	Short: "Sends an evented event",
	Long:  `Sends an evented event to the location and with the data specified`,
	Run: func(cmd *cobra.Command, args []string) {
		log := support.Log()
		defer log.Sync()
		host := args[0]
		port, err := strconv.Atoi(args[1])
		if err != nil {
			log.Error("Error converting port (2nd parameter) to integer port")
		}
		domain := args[2]
		id, err := uuid.Parse(args[3])
		if err != nil {
			log.Error("Error converting id (4th parameter) to UUID")
		}

		config := configuration.Configuration{}
		SendEvent(host, port, domain, id, config, log)

	},
}

func SendEvent(host string, port int, domain string, id uuid.UUID, config configuration.Configuration, log *zap.SugaredLogger) {
	log.Info("Starting...")
	target := config.EventHandlerURL()
	log.Info(target)

	tracer, closer := jaeger.SetupJaeger(config.AppName(), log)
	defer jaeger.CloseJaeger(closer, log)

	conn := grpcWithInterceptors.GenerateConfiguredConn(target, log, tracer)

	sh := evented.NewSagaCoordinatorClient(conn)
	log.Info("Client Created...")

	id, err := uuid.NewRandom()
	protoId := evented_proto.UUIDToProto(id)

	var pages []*evented.EventPage
	for i := 0; i <= 1; i++ {
		pages = append(pages, framework.NewEventPage(uint32(i), false, nil))
	}
	eventBook := &evented.EventBook{
		Cover: &evented.Cover{
			Domain: config.Domain(),
			Root:   &protoId,
		},
		Pages: pages,
	}
	res, err := sh.HandleSync(context.Background(), eventBook)
	log.Info(res)
	if err != nil {
	}
}

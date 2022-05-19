package root

import (
	"fmt"
	todoACtx "github.com/benjaminabbitt/evented/applications/todo/actx"
	"github.com/benjaminabbitt/evented/applications/todo/commands/send"
	"github.com/benjaminabbitt/evented/applications/todo/commands/serve"
	"github.com/benjaminabbitt/evented/applications/todo/configuration"
	"github.com/benjaminabbitt/evented/support"
	eventedACtx "github.com/benjaminabbitt/evented/support/actx"
	"github.com/benjaminabbitt/evented/support/consul"
	"github.com/benjaminabbitt/evented/support/jaeger"
	"github.com/benjaminabbitt/evented/support/serpent"
	"github.com/dsnet/try"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

var Cmd = &cobra.Command{
	Use:   "sender",
	Short: "root command",
	Long:  "long root command",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		try.E(serpent.BindFlags(cmd, viper.GetViper(), "todo"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		try.E1(fmt.Fprintln(os.Stdout, "In Command Run"))
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context().(*todoACtx.TodoSendContext)
		log := ctx.Actx.Log
		log.Info("In Root Command PostRun")
		try.E(log.Sync())
		try.E(ctx.Actx.Log.Sync())
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	Cmd.AddCommand(send.Cmd)
	Cmd.AddCommand(serve.Cmd)
}

const NAME = "todo"

func Execute() error {
	log := support.Log()

	tracer, closer := jaeger.SetupJaeger(fmt.Sprintf("%s-%s", NAME, "test"), log)
	defer closer.Close()

	config := &configuration.Configuration{}
	config.SetName(NAME)

	var hold interface{}
	hold = try.E1(support.Initialize(log, config))
	config = hold.(*configuration.Configuration)

	setupConsul(config, log)

	tsc := &todoACtx.TodoSendContext{
		Actx: eventedACtx.Actx{
			Log:    log,
			Tracer: tracer,
		},
		Configuration: config,
	}
	return Cmd.ExecuteContext(tsc)
}

func setupConsul(config *configuration.Configuration, log *zap.SugaredLogger) {
	c := consul.NewEventedConsul(config.ConsulHost, config.Port)
	id, err := uuid.NewRandom()
	if err != nil {
		log.Error(err)
	}
	try.E(c.Register(config.Name, id.String()))
}

func initConfig() {
	try.E1(fmt.Fprintln(os.Stdout, "In initConfig"))
}

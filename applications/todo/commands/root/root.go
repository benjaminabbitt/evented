package root

import (
	"fmt"
	todoACtx "github.com/benjaminabbitt/evented/applications/todo/actx"
	"github.com/benjaminabbitt/evented/applications/todo/commands/send"
	"github.com/benjaminabbitt/evented/applications/todo/commands/serve"
	"github.com/benjaminabbitt/evented/applications/todo/configuration"
	"github.com/benjaminabbitt/evented/support"
	eventedACtx "github.com/benjaminabbitt/evented/support/actx"
	"github.com/benjaminabbitt/evented/support/jaeger"
	"github.com/benjaminabbitt/evented/support/serpent"
	"github.com/dsnet/try"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	Cmd.AddCommand(send.Cmd)
	Cmd.AddCommand(serve.Cmd)
}

func Execute() error {
	log := support.Log()
	v := viper.New()
	v = try.E1(support.Initialize(log, v))
	tracer, closer := jaeger.SetupJaeger(fmt.Sprintf("%s-%s", v.GetString(configuration.Domain), "test"), log)
	defer try.E(closer.Close())

	tsc := &todoACtx.TodoSendContext{
		Actx: eventedACtx.Actx{
			Log:    log,
			Tracer: tracer,
			Config: v,
		},
	}
	return Cmd.ExecuteContext(tsc)
}

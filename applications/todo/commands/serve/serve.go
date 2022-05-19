package serve

import (
	"github.com/benjaminabbitt/evented/applications/todo/actx"
	"github.com/benjaminabbitt/evented/applications/todo/business"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcHealth"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "serve",
	Run: serve,
}

func serve(command *cobra.Command, args []string) {
	ctx := command.Context().(*actx.TodoSendContext)
	log := ctx.Log
	config := ctx.Configuration
	tracer := ctx.Tracer

	support.LogStartup(log, "")

	lis, err := support.OpenPort(config.Port, log)

	rpc := grpcWithInterceptors.GenerateConfiguredServer(log.Desugar(), tracer)

	server := business.NewTodoBusinessLogicServer(ctx)
	evented.RegisterBusinessLogicServer(rpc, server)

	grpcHealth.RegisterHealthChecks(rpc, config.Name, log)

	log.Infow("Starting Business Server...")
	err = rpc.Serve(lis)
	log.Infow("Serving...")
	if err != nil {
		log.Error(err)
	}
}

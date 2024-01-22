package actx

import (
	"github.com/benjaminabbitt/evented/applications/command/command-handler/configuration"
	"github.com/benjaminabbitt/evented/support"
	"github.com/opentracing/opentracing-go"
)

type CommandHandlerContext struct {
	support.BasicApplicationContext
	opentracing.Tracer
	Configuration *configuration.Configuration
}

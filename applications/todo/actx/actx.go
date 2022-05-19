package actx

import (
	"github.com/benjaminabbitt/evented/applications/todo/configuration"
	"github.com/benjaminabbitt/evented/support/actx"
)

type TodoSendContext struct {
	actx.Actx
	Configuration *configuration.Configuration
}

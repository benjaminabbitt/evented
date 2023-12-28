package sender

import (
	"github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
)

type NoOp struct {
}

func (n NoOp) Handle(_ *evented.EventBook) (err error) {
	return nil
}

func (n NoOp) Run() {
}

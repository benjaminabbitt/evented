package console

import (
	"fmt"
	"github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/transport"
)

type Sender struct {
}

func (sender Sender) Send(eventBook evented_core.EventBook){
		fmt.Printf("%+v", eventBook)
}

func NewConsoleSender() transport.Projection {
	return &Sender{}
}

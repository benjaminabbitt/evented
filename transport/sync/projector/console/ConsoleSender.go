package console

import (
	"fmt"
	"github.com/benjaminabbitt/evented/transport"
	"github.com/benjaminabbitt/evented/proto/core"
)

type Sender struct {
}

func (sender Sender) Send(eventBook evented_core.EventBook){
		fmt.Printf("%+v", eventBook)
}

func NewConsoleSender() transport.Projection {
	return &Sender{}
}

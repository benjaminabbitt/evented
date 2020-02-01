package evented_consoleTransport

import (
	"fmt"
	"github.com/benjaminabbitt/evented/transport"
)

type ConsoleSender struct {
}

func (sender ConsoleSender) Send(evt transport.TransportEvent){
	fmt.Printf("%+v", evt)
}


func NewConsoleSender() transport.EventTransportSender {
	return &ConsoleSender{}
}

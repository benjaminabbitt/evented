package create

import (
	"fmt"
	"github.com/benjaminabbitt/evented/applications/todo/actx"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/dsnet/try"
)

func Send(actx actx.TodoSendContext, host string, port uint16, command *evented.CommandBook) (A *evented.SynchronousProcessingResponse) {
	client := CreateClient(actx, host, port)
	result := try.E1(client.Handle(actx, command))
	return result
}

func CreateClient(actx actx.TodoSendContext, host string, port uint16) evented.BusinessCoordinatorClient {
	target := fmt.Sprintf("%s:%d", host, port)
	conn := grpcWithInterceptors.GenerateConfiguredConn(target, actx.Log, actx.Tracer)
	client := evented.NewBusinessCoordinatorClient(conn)
	return client
}

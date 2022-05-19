package business

import (
	"github.com/benjaminabbitt/evented/applications/todo/actx"
	"github.com/benjaminabbitt/evented/applications/todo/business/states"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/todo/proto"
	"github.com/benjaminabbitt/evented/support"
	"github.com/dsnet/try"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewTodoBusinessLogicServer(actx *actx.TodoSendContext) *TodoBusinessLogicServer {
	return &TodoBusinessLogicServer{
		actx: actx,
	}
}

type TodoBusinessLogicServer struct {
	evented.UnimplementedBusinessLogicServer
	actx *actx.TodoSendContext
}

func (o TodoBusinessLogicServer) Handle(ctx context.Context, in *evented.ContextualCommand) (eb *evented.EventBook, err error) {
	o.actx.Log.Infow("Business Logic Handle", "contextualCommand", in)
	var eventPages []*evented.EventPage
	var todo states.Todo = states.NonextantTodo{}
	//TODO: harden
	ts := timestamppb.Now()
	if in.Events != nil {
		tmp := NilClearSlice(in.Events.Pages)
		for _, event := range tmp {
			event := try.E1(event.Event.UnmarshalNew())
			switch event := event.(type) {
			case *proto.TodoCreated:
				todo = todo.Created(event)
			case *proto.TodoEdited:
				todo = todo.Edited(event)
			case *proto.StatusSet:
				todo = todo.StatusSet(event)
			}
		}
	}
	for _, commandPage := range in.Command.Pages {
		var result *anypb.Any
		command := try.E1(commandPage.Command.UnmarshalNew())
		switch command := command.(type) {
		case *proto.CreateTodo:
			result = try.E1(anypb.New(try.E1(todo.Create(command))))
		case *proto.EditTodo:
			result = try.E1(anypb.New(try.E1(todo.Edit(command))))
		case *proto.SetStatus:
			result = try.E1(anypb.New(try.E1(todo.SetStatus(command))))
		}
		eventPage := &evented.EventPage{
			Sequence:    &evented.EventPage_Num{Num: commandPage.Sequence},
			CreatedAt:   ts,
			Event:       result,
			Synchronous: false,
		}
		eventPages = append(eventPages, eventPage)
	}

	eventBook := &evented.EventBook{
		Cover:    in.Command.Cover,
		Pages:    eventPages,
		Snapshot: nil,
	}

	o.actx.Log.Infow("Todo Handle", "eventBook", support.StringifyEventBook(eventBook))

	return eventBook, nil
}

func NilClearSlice[T any](slice []T) []T {
	if slice == nil {
		return []T{}
	} else {
		return slice
	}
}

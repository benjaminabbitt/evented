package states

import (
	"fmt"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/todo/proto"
)

type IncompleteTodo struct {
}

func (i *IncompleteTodo) Create(command *proto.CreateTodo) (event *proto.TodoCreated, err error) {
	return nil, fmt.Errorf("todo already created")
}

func (i *IncompleteTodo) Created(event *proto.TodoCreated) Todo {
	return i
}

func (i *IncompleteTodo) Edit(command *proto.EditTodo) (event *proto.TodoEdited, err error) {
	return &proto.TodoEdited{
		Payload:  command.Payload,
		Extended: command.Extended,
	}, nil
}

func (i *IncompleteTodo) Edited(event *proto.TodoEdited) Todo {
	return i
}

func (i *IncompleteTodo) SetStatus(command *proto.SetStatus) (event *proto.StatusSet, err error) {
	if command.Done {
		return &proto.StatusSet{Done: true}, nil
	} else {
		return nil, fmt.Errorf("todo already incomplete")
	}
}

func (i *IncompleteTodo) StatusSet(event *proto.StatusSet) Todo {
	if event.Done {
		return &CompleteTodo{}
	} else {
		return i
	}
}

package states

import (
	"fmt"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/todo/proto"
)

type CompleteTodo struct{}

func (c CompleteTodo) Create(command *proto.CreateTodo) (event *proto.TodoCreated, err error) {
	return nil, fmt.Errorf("todo already created")
}

func (c CompleteTodo) Created(event *proto.TodoCreated) Todo {
	return c
}

func (c CompleteTodo) Edit(command *proto.EditTodo) (event *proto.TodoEdited, err error) {
	return &proto.TodoEdited{
		Payload:  command.Payload,
		Extended: command.Extended,
	}, nil
}

func (c CompleteTodo) Edited(event *proto.TodoEdited) Todo {
	return c
}

func (c CompleteTodo) SetStatus(command *proto.SetStatus) (event *proto.StatusSet, err error) {
	if !command.Done {
		return &proto.StatusSet{Done: false}, nil
	} else {
		return nil, fmt.Errorf("todo already complete")
	}
}

func (c CompleteTodo) StatusSet(event *proto.StatusSet) Todo {
	if event.Done {
		return c
	} else {
		return &NonextantTodo{}
	}
}

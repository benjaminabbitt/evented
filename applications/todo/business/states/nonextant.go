package states

import (
	"fmt"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/todo/proto"
)

type NonextantTodo struct {
}

func (n NonextantTodo) Create(command *proto.CreateTodo) (event *proto.TodoCreated, err error) {
	return &proto.TodoCreated{
		Todo:     command.Todo,
		Extended: command.Extended,
	}, nil
}

func (n NonextantTodo) Created(event *proto.TodoCreated) Todo {
	return NonextantTodo{}
}

func (n NonextantTodo) Edit(command *proto.EditTodo) (event *proto.TodoEdited, err error) {
	return nil, fmt.Errorf("todo doesn't exist")
}

func (n NonextantTodo) Edited(event *proto.TodoEdited) Todo {
	return &n
}

func (n NonextantTodo) SetStatus(command *proto.SetStatus) (event *proto.StatusSet, err error) {
	return nil, fmt.Errorf("todo doesn't exist")
}

func (n NonextantTodo) StatusSet(event *proto.StatusSet) Todo {
	return &n
}

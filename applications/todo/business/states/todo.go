package states

import (
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/todo/proto"
)

type Todo interface {
	Create(command *proto.CreateTodo) (event *proto.TodoCreated, err error)
	Created(event *proto.TodoCreated) Todo
	Edit(command *proto.EditTodo) (event *proto.TodoEdited, err error)
	Edited(event *proto.TodoEdited) Todo
	SetStatus(command *proto.SetStatus) (event *proto.StatusSet, err error)
	StatusSet(event *proto.StatusSet) Todo
}

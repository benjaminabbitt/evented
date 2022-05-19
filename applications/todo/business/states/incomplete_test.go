package states

import (
	"fmt"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/todo/proto"
	"github.com/stretchr/testify/suite"
	"testing"
)

type IncompleteSuite struct {
	suite.Suite
	sut *IncompleteTodo
}

func (s *IncompleteSuite) TestCreation() {
	s.sut = &IncompleteTodo{}
	todo := &proto.MinimumTodo{Title: "test"}
	extended := &proto.ExtendedTodo{
		Body:      StringPointer("body"),
		Due:       nil,
		Duration:  nil,
		Important: BoolPointer(true),
		RemindAt:  nil,
		Priority:  nil,
	}
	_, err := s.sut.Create(&proto.CreateTodo{
		Todo:     todo,
		Extended: extended,
	})

	s.Assert().Error(fmt.Errorf("todo already created"), err)
}

func (s *IncompleteSuite) TestCreated() {
	s.sut = &IncompleteTodo{}
	todo := &proto.MinimumTodo{Title: "test"}
	extended := &proto.ExtendedTodo{
		Body:      StringPointer("body"),
		Due:       nil,
		Duration:  nil,
		Important: BoolPointer(true),
		RemindAt:  nil,
		Priority:  nil,
	}
	created := s.sut.Created(&proto.TodoCreated{
		Todo:     todo,
		Extended: extended,
	})

	expected := &IncompleteTodo{}

	s.Assert().Equal(expected, created)

}

func (s *IncompleteSuite) TestModify() {
	s.sut = &IncompleteTodo{}
	todo := &proto.MinimumTodoOptions{}
	extended := &proto.ExtendedTodo{}
	newState, err := s.sut.Edit(&proto.EditTodo{
		Payload:  todo,
		Extended: extended,
	})
	expected := &proto.TodoEdited{
		Payload:  todo,
		Extended: extended,
	}
	s.Assert().Equal(expected, newState)
	s.Assert().Nil(err)
}

func (s *IncompleteSuite) TestModified() {
	s.sut = &IncompleteTodo{}
	newState := s.sut.Edited(&proto.TodoEdited{})
	s.Assert().Equal(s.sut, newState)
}

func (s *IncompleteSuite) TestComplete() {
	s.sut = &IncompleteTodo{}
	status := &proto.SetStatus{Done: false}
	_, err := s.sut.SetStatus(status)
	s.Assert().Errorf(err, "todo doesn't exist")
}

func (s *IncompleteSuite) TestCompleted() {
	s.sut = &IncompleteTodo{}
	newState := s.sut.StatusSet(&proto.StatusSet{Done: true})
	expected := &CompleteTodo{}
	s.Assert().Equal(expected, newState)
}

func TestIncompleteSuite(t *testing.T) {
	suite.Run(t, new(IncompleteSuite))
}

package states

import (
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/todo/proto"
	"github.com/dsnet/try"
	"github.com/stretchr/testify/suite"
	"testing"
)

type NonExtantSuite struct {
	suite.Suite
	sut *NonextantTodo
}

func (s *NonExtantSuite) TestCreation() {
	s.sut = &NonextantTodo{}
	todo := &proto.MinimumTodo{Title: "test"}
	extended := &proto.ExtendedTodo{
		Body:      StringPointer("body"),
		Due:       nil,
		Duration:  nil,
		Important: BoolPointer(true),
		RemindAt:  nil,
		Priority:  nil,
	}
	event := try.E1(s.sut.Create(&proto.CreateTodo{
		Todo:     todo,
		Extended: extended,
	}))

	expected := &proto.TodoCreated{
		Todo:     todo,
		Extended: extended,
	}

	s.Assert().Equal(event, expected)
}

func (s *NonExtantSuite) TestCreated() {
	s.sut = &NonextantTodo{}
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

	expected := NonextantTodo{}

	s.Assert().Equal(expected, created)

}

func (s *NonExtantSuite) TestModify() {
	s.sut = &NonextantTodo{}
	todo := &proto.MinimumTodoOptions{}
	extended := &proto.ExtendedTodo{}
	_, err := s.sut.Edit(&proto.EditTodo{
		Payload:  todo,
		Extended: extended,
	})
	s.Assert().Errorf(err, "todo doesn't exist")
}

func (s *NonExtantSuite) TestModified() {
	s.sut = &NonextantTodo{}
	newState := s.sut.Edited(&proto.TodoEdited{})
	s.Assert().Equal(s.sut, newState)
}

func (s *NonExtantSuite) TestComplete() {
	s.sut = &NonextantTodo{}
	status := &proto.SetStatus{Done: false}
	_, err := s.sut.SetStatus(status)
	s.Assert().Errorf(err, "todo doesn't exist")
}

func (s *NonExtantSuite) TestCompleted() {
	s.sut = &NonextantTodo{}
	newState := s.sut.StatusSet(&proto.StatusSet{})
	s.Assert().Equal(s.sut, newState)
}

func StringPointer(str string) *string {
	return &str
}

func BoolPointer(b bool) *bool {
	return &b
}

func TestNonExtantSuite(t *testing.T) {
	suite.Run(t, new(NonExtantSuite))
}

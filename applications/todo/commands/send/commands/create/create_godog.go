package create

import (
	"context"
	"fmt"
	todoACtx "github.com/benjaminabbitt/evented/applications/todo/actx"
	"github.com/benjaminabbitt/evented/applications/todo/business"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/todo/proto"
	"github.com/benjaminabbitt/evented/support/actx"
	"github.com/cucumber/godog"
	"github.com/dsnet/try"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type testSuite struct {
	Title    string
	Sequence uint32
	Domain   string
	Id       uuid.UUID
	result   *evented.EventBook
	Sut      *business.TodoBusinessLogicServer
}

func (td *testSuite) aTitleOf(arg1 string) error {
	td.Title = arg1
	return nil
}

func (td *testSuite) iRunThis() (err error) {
	cmd := &evented.ContextualCommand{
		Events:  nil,
		Command: CreateCommmand(td.Id, td.Title),
	}
	eb := try.E1(td.Sut.Handle(context.Background(), cmd))
	try.Handle(&err)
	td.result = eb
	return err
}

func (td *testSuite) theDomainShouldBe(arg1 string) error {
	if td.result.Cover.Domain == arg1 {
		return nil
	} else {
		return fmt.Errorf("domain didn't match")
	}
}

func (td *testSuite) theIdShouldBeSet() (err error) {
	if try.E1(evented_proto.ProtoToUUID(td.result.Cover.Root)) == td.Id {
		try.Handle(&err)
		return err
	} else {
		return fmt.Errorf("ID was not set or didn't match, expected %s; found: %s",
			td.Id, td.result.Cover.Root)
	}
}

func (td *testSuite) theSequenceShouldBe(arg1 int) error {
	var seq uint32
	switch protoSeq := td.result.Pages[0].Sequence.(type) {
	case *evented.EventPage_Num:
		seq = protoSeq.Num
	case *evented.EventPage_Force:
		return fmt.Errorf("found a forced event page, expected unforced sequence %d", arg1)
	}
	if seq != uint32(arg1) {
		return fmt.Errorf("sequences didn't match, expected %d; found %d", arg1, seq)
	}
	return nil
}

func (td *testSuite) thereShouldBeAnEventCreatedWithTheTitle(arg1 string) (err error) {
	created := new(proto.TodoCreated)
	try.E(td.result.Pages[0].Event.UnmarshalTo(created))
	try.Handle(&err)
	if created.Todo.Title != td.Title {
		return fmt.Errorf("titles didn't match.  expected \"%s\", received \"%s\"", td.Title, created.Todo.Title)
	}
	return err

}

func Test() {

}

func InitializeScenario(ctx *godog.ScenarioContext) {
	id, _ := uuid.NewRandom()
	todoCtx := &todoACtx.TodoSendContext{
		Actx: actx.Actx{
			Context: context.Background(),
			Log:     zap.S(),
		},
	}
	sut := business.NewTodoBusinessLogicServer(todoCtx)
	var td = testSuite{
		Id:  id,
		Sut: sut,
	}
	ctx.Step(`^a title of "([^"]*)"$`, td.aTitleOf)
	ctx.Step(`^I run this$`, td.iRunThis)
	ctx.Step(`^the domain should be "([^"]*)"$`, td.theDomainShouldBe)
	ctx.Step(`^the id should be set$`, td.theIdShouldBeSet)
	ctx.Step(`^the sequence should be (\d+)$`, td.theSequenceShouldBe)
	ctx.Step(`^there should be an event created with the title "([^"]*)"$`, td.thereShouldBeAnEventCreatedWithTheTitle)
}

package business

import (
	"github.com/benjaminabbitt/evented/applications/todo/actx"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/todo/proto"
	actx2 "github.com/benjaminabbitt/evented/support/actx"
	"github.com/dsnet/try"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	log2 "log"
	"testing"
)

func TestNonExtantSuite(t *testing.T) {
	suite.Run(t, new(BusinessTestSuite))
}

type BusinessTestSuite struct {
	suite.Suite
	sut *TodoBusinessLogicServer
}

func (o *BusinessTestSuite) TestBasic() {
	try.F(log2.Fatal)
	log := try.E1(zap.NewDevelopment())
	slog := log.Sugar()
	ctx := &actx.TodoSendContext{
		Actx: actx2.Actx{
			Log:    slog,
			Tracer: nil,
		},
		Configuration: nil,
	}

	o.sut = NewTodoBusinessLogicServer(ctx)

	id := evented_proto.UUIDToProto(try.E1(uuid.NewRandom()))

	cmd := &evented.ContextualCommand{
		Events: &evented.EventBook{},
		Command: &evented.CommandBook{
			Cover: &evented.Cover{
				Domain: "todo",
				Root:   &id,
			},
			Pages: []*evented.CommandPage{{
				Sequence:    0,
				Synchronous: false,
				Command: try.E1(anypb.New(
					&proto.CreateTodo{
						Todo:     &proto.MinimumTodo{Title: "title"},
						Extended: &proto.ExtendedTodo{},
					},
				)),
			}},
		},
	}
	events, err := o.sut.Handle(ctx, cmd)
	println(events)
	assert.Nil(o.T(), err)
}

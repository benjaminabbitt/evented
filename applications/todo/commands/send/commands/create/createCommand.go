package create

import (
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/todo/proto"
	"github.com/dsnet/try"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"
)

func CreateCommmand(id uuid.UUID, title string) (book *evented.CommandBook) {
	pid := evented_proto.UUIDToProto(id)
	book = &evented.CommandBook{
		Cover: &evented.Cover{
			Domain: "todo",
			Root:   &pid,
		},
		Pages: []*evented.CommandPage{{
			Sequence:    0,
			Synchronous: false,
			Command: try.E1(anypb.New(&proto.CreateTodo{
				Todo: &proto.MinimumTodo{Title: title},
				Extended: &proto.ExtendedTodo{
					Body:      nil,
					Due:       nil,
					Duration:  nil,
					Important: nil,
					RemindAt:  nil,
					Priority:  nil,
				},
			})),
		}},
	}
	return book
}

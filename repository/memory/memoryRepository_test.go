package memoryRepository

import (
	"github.com/benjaminabbitt/evented/framework"
	eventedproto "github.com/benjaminabbitt/evented/proto/core"
	"github.com/golang/protobuf/ptypes"
)
import (
	"testing"
)

func TestIdentity(t *testing.T) {
}

func TestAdd(t *testing.T) {
	mr := NewMemoryRepository()
	any, _ := ptypes.MarshalAny(&eventedproto.Empty{})
	fevent := framework.Event{
		Id:       "",
		Sequence: 0,
		Details:  *any,
	}
	_ = mr.Add(fevent)
}

func TestFiltering(t *testing.T) {
	mr := NewMemoryRepository()
	nonEvent, err := ptypes.MarshalAny(&eventedproto.Empty{})
	if err == nil {
		se0 := framework.Event{
			Id:       "a",
			Sequence: 0,
			Details:  *nonEvent,
		}
		se1 := framework.Event{
			Id:       "a",
			Sequence: 1,
			Details:  *nonEvent,
		}
		se2 := framework.Event{
			Id:       "a",
			Sequence: 2,
			Details:  *nonEvent,
		}
		mr.Add(se0)
		mr.Add(se1)
		mr.Add(se2)
	}

	t.Run("TestGetEventsTo", func(t *testing.T) {
		events, _ := mr.GetTo("a", 1)
		if len(events) != 2 {
			t.Fail()
		}
	})

	t.Run("TestGetEventsFromTo", func(t *testing.T) {
		events, _ := mr.GetFromTo("a", 1, 1)
		if len(events) != 1 {
			t.Fail()
		}
	})
}

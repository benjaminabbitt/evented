package memoryRepository


import	"github.com/golang/protobuf/ptypes"
import (
	framework "github.com/benjaminabbitt/evented"
	"testing"
)

func TestIdentity(t *testing.T) {
}

func TestAdd(t *testing.T) {
	mr := NewMemoryRepository()
	businessEvent, err := ptypes.MarshalAny(&framework.EmptyMessage{})

	if err == nil {
		se := framework.StorageEvent{
			Id:       "a",
			Sequence: 0,
			Details:  *businessEvent,
		}
		mr.Add(se)
		events := mr.Get("a")
		if events[0].Sequence != se.Sequence {
			t.Fail()
		}
	}

}

func TestFiltering(t *testing.T){
	mr := NewMemoryRepository()
	nonEvent, err := ptypes.MarshalAny(&framework.EmptyMessage{})
	if err == nil {
		se0 := framework.StorageEvent{
			Id:       "a",
			Sequence: 0,
			Details:  *nonEvent,
		}
		se1 := framework.StorageEvent{
			Id:       "a",
			Sequence: 1,
			Details:  *nonEvent,
		}
		se2 := framework.StorageEvent{
			Id:       "a",
			Sequence: 2,
			Details:  *nonEvent,
		}
		mr.Add(se0)
		mr.Add(se1)
		mr.Add(se2)
	}

	t.Run("TestGetEventsTo", func(t *testing.T){
		events := mr.GetTo("a", 1)
		if len(events) != 2 {
			t.Fail()
		}
	})

	t.Run("TestGetEventsFromTo", func(t *testing.T){
		events := mr.GetFromTo("a", 1, 1)
		if len(events)!= 1 {
			t.Fail()
		}
	})
}

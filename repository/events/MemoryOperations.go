package events

import (
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"go.uber.org/zap"
)

func GetSequence(log *zap.SugaredLogger, page *evented.EventPage) uint32 {
	var sequence uint32
	switch s := page.Sequence.(type) {
	case *evented.EventPage_Num:
		sequence = s.Num
	default:
		log.Error("Attempted to retrieve sequence from event without sequence set.  This should not happen")
	}
	return sequence
}

func SetSequence(page *evented.EventPage, sequence uint32) {
	page.Sequence = &evented.EventPage_Num{Num: sequence}
}

func ExtractUntilFirstForced(events []*evented.EventPage) (numbered []*evented.EventPage, forced *evented.EventPage, remainder []*evented.EventPage) {
	for idx, page := range events {
		switch page.GetSequence().(type) {
		case *evented.EventPage_Force:
			return events[:idx], page, events[idx+1:]
		}
	}
	return events, nil, nil
}

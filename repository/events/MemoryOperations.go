package events

import (
	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
	"go.uber.org/zap"
)

func GetSequence(log *zap.SugaredLogger, page *evented_core.EventPage) uint32 {
	var sequence uint32
	switch s := page.Sequence.(type) {
	case *evented_core.EventPage_Num:
		sequence = s.Num
	default:
		log.Error("Attempted to retrieve sequence from event without sequence set.  This should not happen")
	}
	return sequence
}

func SetSequence(page *evented_core.EventPage, sequence uint32) {
	page.Sequence = &evented_core.EventPage_Num{Num: sequence}
}

func ExtractUntilFirstForced(events []*evented_core.EventPage) (numbered []*evented_core.EventPage, forced *evented_core.EventPage, remainder []*evented_core.EventPage) {
	for idx, page := range events {
		switch page.GetSequence().(type) {
		case *evented_core.EventPage_Force:
			return events[:idx], page, events[idx+1:]
		}
	}
	return events, nil, nil
}

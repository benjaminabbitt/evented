package events

import (
	core "github.com/benjaminabbitt/evented/proto/evented/core"
	"go.uber.org/zap"
)

func GetSequence(log *zap.SugaredLogger, page *core.EventPage) uint32 {
	var sequence uint32
	switch s := page.Sequence.(type) {
	case *core.EventPage_Num:
		sequence = s.Num
	default:
		log.Error("Attempted to retrieve sequence from event without sequence set.  This should not happen")
	}
	return sequence
}

func SetSequence(page *core.EventPage, sequence uint32) {
	page.Sequence = &core.EventPage_Num{Num: sequence}
}

func ExtractUntilFirstForced(events []*core.EventPage) (numbered []*core.EventPage, forced *core.EventPage, remainder []*core.EventPage) {
	for idx, page := range events {
		switch page.GetSequence().(type) {
		case *core.EventPage_Force:
			return events[:idx], page, events[idx+1:]
		}
	}
	return events, nil, nil
}

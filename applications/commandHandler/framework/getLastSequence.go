package framework

import evented_core "github.com/benjaminabbitt/evented/proto/core"

type GetLastSequence struct {
	book evented_core.EventBook
}

func (o *GetLastSequence) GetLastSequence() uint32 {
	if len(o.book.Pages) > 0 {
		return o.book.Pages[len(o.book.Pages)-1].Sequence.(*evented_core.EventPage_Num).Num
	} else {
		return o.book.Snapshot.Sequence
	}
}
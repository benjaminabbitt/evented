package support

import (
	"errors"
	"github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
)

func GetSequence(eb *evented.EventPage) (seq uint32, err error) {
	switch page := eb.Sequence.(type) {
	case *evented.EventPage_Num:
		return page.Num, nil
	default:
		return 0, errors.New("sequence not set")
	}
}

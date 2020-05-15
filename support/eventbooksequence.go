package support

import (
	"errors"
	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
)

func GetSequence(eb *evented_core.EventPage) (seq uint32, err error) {
	switch page := eb.Sequence.(type) {
	case *evented_core.EventPage_Num:
		return page.Num, nil
	default:
		return 0, errors.New("sequence not set")
	}
}

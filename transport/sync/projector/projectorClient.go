package projector

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/support"
	"go.uber.org/zap"
)

type ProjectorClient struct {
	log      *zap.SugaredLogger
	Sequence uint32
}

func (c *ProjectorClient) HandleSync(in *evented_core.EventBook) (*evented_core.Projection, error) {
	c.log.Infow("HandleSync", "eventBook", support.StringifyEventBook(in))
	c.updateSequence(in)
	projection := &evented_core.Projection{
		Cover:      in.Cover,
		Projector:  "simple",
		Sequence:   c.Sequence,
		Projection: nil,
	}
	return projection, nil
}

func (c *ProjectorClient) Handle(in *evented_core.EventBook) error {
	c.log.Infow("Handle", "eventBook", support.StringifyEventBook(in))
	c.updateSequence(in)
	return nil
}

func (c *ProjectorClient) updateSequence(eb *evented_core.EventBook) {
	for _, page := range eb.Pages {
		switch s := page.Sequence.(type){
		case *evented_core.EventPage_Num:
			c.Sequence = s.Num
		default:
			c.log.Warnw("Received sequence without an assigned number.  This should not exist at this point in processing", "eventBook", *eb)
		}

	}
}

func NewProjectorClient(log *zap.SugaredLogger) *ProjectorClient {
	return &ProjectorClient{log, 0}
}

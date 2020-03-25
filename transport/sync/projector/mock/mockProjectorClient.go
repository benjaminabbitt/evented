package mock

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"go.uber.org/zap"
)

type ProjectorClient struct{
	log *zap.SugaredLogger
	Sequence uint32
}

func (c *ProjectorClient) HandleSync(in *evented_core.EventBook)(*evented_core.Projection, error){
	c.log.Infow("HandleSync", "eventBook", in)
	c.updateSequence(in)
	projection := &evented_core.Projection{
		Cover:      in.Cover,
		Projector:  "mock",
		Sequence:   c.Sequence,
		Projection: nil,
	}
	c.log.Infow("ProjectSync - End", "projection", projection)
	return projection, nil
}

func (c *ProjectorClient) Handle(in *evented_core.EventBook) error {
	c.log.Infow("Handle", "eventBook", in)
	c.updateSequence(in)
	return nil
}

func (c *ProjectorClient) updateSequence(eb *evented_core.EventBook){
	for _, page := range eb.Pages {
		c.Sequence = page.Sequence
	}
}

func NewProjectorClient(log *zap.SugaredLogger) *ProjectorClient {
	return &ProjectorClient{log, 0}
}
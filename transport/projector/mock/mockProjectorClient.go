package mock

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/transport"
	"github.com/stretchr/testify/mock"
)

type ProjectorClient struct{
	mock.Mock
	Sequence uint32
}

func (c *ProjectorClient) ProjectSync(in *evented_core.EventBook)(*evented_core.Projection, error){
	c.Called(in)
	c.updateSequence(in)
	return &evented_core.Projection{
		Cover:      in.Cover,
		Projector:  "mock",
		Sequence:   c.Sequence,
		Projection: nil,
	}, nil
}

func (c *ProjectorClient) Project(in *evented_core.EventBook) error {
	c.Called(in)
	c.updateSequence(in)
	return nil
}

func (c *ProjectorClient) updateSequence(eb *evented_core.EventBook){
	for _, page := range eb.Pages {
		c.Sequence = page.Sequence
	}
}

func NewProjectorClient() transport.Projection {
	return &ProjectorClient{}
}
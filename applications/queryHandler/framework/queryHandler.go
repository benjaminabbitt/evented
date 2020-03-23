package framework

import (
	"github.com/benjaminabbitt/evented"
	"github.com/benjaminabbitt/evented/proto/core"
	evented_query "github.com/benjaminabbitt/evented/proto/query"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	"golang.org/x/net/context"
)

type Server struct {
	evented_query.UnimplementedEventQueryServer
	Repo eventBook.Repository
}

func (s Server) GetEventBook(ctx context.Context, req *evented_query.Query) (*evented_core.EventBook, error) {
	var eventBook evented_core.EventBook
	var err error
	if req.UpperBound != 0 && req.LowerBound != 0 {
		eventBook, err = s.Repo.GetFromTo(req.Root, req.LowerBound, req.UpperBound)
		evented.FailOnError(err, "Failed to get events")
	} else {
		eventBook, err = s.Repo.GetFrom(req.Root, req.LowerBound)
		evented.FailOnError(err, "Failed to get events")
	}
	return &eventBook, nil
}


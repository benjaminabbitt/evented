package framework

import (
	"context"
	"errors"
	"github.com/benjaminabbitt/evented/applications/command/command-handler/business/client"
	"github.com/benjaminabbitt/evented/applications/command/command-handler/framework/transport"
	eventedproto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(actx CommandHandlerApplicationContext, eventBookRepository eventBook.Storer, transports transport.Holder, businessClient client.BusinessClient) Server {
	return Server{
		retry:               actx.RetryStrategy(),
		log:                 actx.Log(),
		eventBookRepository: eventBookRepository,
		transports:          transports,
		businessClient:      businessClient,
	}
}

func (o *Server) Shutdown() {
	o.server.GracefulStop()
}

type Server struct {
	evented.UnimplementedBusinessCoordinatorServer
	retry               backoff.BackOff
	log                 *zap.SugaredLogger
	eventBookRepository eventBook.Storer
	transports          transport.Holder
	businessClient      client.BusinessClient
	server              *grpc.Server
}

func (o Server) Handle(ctx context.Context, in *evented.CommandBook) (result *evented.SynchronousProcessingResponse, err error) {
	id, err := eventedproto.ProtoToUUID(in.Cover.Root)
	if err != nil {
		return nil, err
	}
	priorState, err := o.eventBookRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	contextualCommand := &evented.ContextualCommand{
		Events:  priorState,
		Command: in,
	}

	var businessResponse *evented.EventBook
	err = backoff.Retry(func() error {
		businessResponse, err = o.businessClient.Handle(ctx, contextualCommand)
		return err
	}, o.retry)
	if err != nil {
		return nil, err
	}
	result, err = o.handleEventBook(ctx, businessResponse)
	if err != nil {
		return nil, err
	}

	result.Books = append(result.Books, businessResponse)
	return result, err
}

func (o Server) handleEventBook(ctx context.Context, eb *evented.EventBook) (result *evented.SynchronousProcessingResponse, rerr error) {
	result = &evented.SynchronousProcessingResponse{}
	result.Books = []*evented.EventBook{eb}

	err := o.eventBookRepository.Put(ctx, eb)
	if err != nil {
		return nil, err
	}

	sync, _, err := o.extractSynchronous(eb)
	if err != nil {
		return nil, err
	}

	syncResults, err := o.executeSyncProjections(ctx, sync)
	if err != nil {
		return result, err
	}
	result.Projections = append(result.Projections, syncResults...)

	otherDomainEventBooks, otherProjections, err := o.executeSyncSagas(ctx, sync)
	if err != nil {
		return nil, multierror.Append(rerr, err)
	}
	result.Books = append(result.Books, otherDomainEventBooks...)
	result.Projections = append(result.Projections, otherProjections...)

	for _, t := range o.transports.GetTransports() {
		t <- eb
	}

	return result, nil
}

func (o Server) executeSyncSagas(ctx context.Context, sync *evented.EventBook) (eventBooks []*evented.EventBook, projections []*evented.Projection, rerr error) {
	for _, syncSaga := range o.transports.GetSaga() {
		var response *evented.SynchronousProcessingResponse
		var err error
		err = backoff.Retry(func() error {
			response, err = syncSaga.HandleSync(ctx, sync)
			return err
		}, o.retry)

		if err != nil {
			o.log.Error(err)
			rerr = multierror.Append(rerr, err)
			continue
		}
		eventBooks = append(eventBooks, response.Books...)
		projections = append(projections, response.Projections...)
	}
	return eventBooks, projections, rerr
}

func (o Server) executeSyncProjections(ctx context.Context, sync *evented.EventBook) (result []*evented.Projection, rerr error) {
	for _, syncProjector := range o.transports.GetProjectors() {
		var response *evented.Projection
		var err error
		err = backoff.Retry(func() error {
			response, err = syncProjector.HandleSync(ctx, sync)
			return err
		}, o.retry)
		if err != nil {
			o.log.Error(err)
			rerr = multierror.Append(rerr, err)
			continue
		}
		result = append(result, response)
	}
	return result, rerr
}

func (o Server) extractSynchronous(originalBook *evented.EventBook) (synchronous *evented.EventBook, async *evented.EventBook, err error) {
	if len(originalBook.Pages) == 0 {
		return nil, nil, errors.New("event book has no pages -- not correct in this context")
	}
	var lastIdx uint32
	for idx, event := range originalBook.Pages {
		if event.Synchronous {
			lastIdx = uint32(idx) + 1
		}
	}
	synchronous = new(evented.EventBook)
	synchronous.Pages = originalBook.Pages[:lastIdx]
	synchronous.Cover = originalBook.Cover
	synchronous.Snapshot = originalBook.Snapshot

	async = new(evented.EventBook)
	async.Pages = originalBook.Pages[lastIdx:]
	async.Cover = originalBook.Cover
	async.Snapshot = nil

	return synchronous, async, err
}

func (o Server) Record(ctx context.Context, in *evented.EventBook) (response *evented.SynchronousProcessingResponse, err error) {
	return o.handleEventBook(ctx, in)
}

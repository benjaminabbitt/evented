package framework

import (
	"context"
	"errors"
	"github.com/benjaminabbitt/evented/applications/commandHandler/business/client"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework/transport"
	eventedproto "github.com/benjaminabbitt/evented/proto"
	business "github.com/benjaminabbitt/evented/proto/evented/business/coordinator"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(eventBookRepository eventBook.Storer, transports transport.TransportHolder, businessClient client.BusinessClient, log *zap.SugaredLogger) Server {
	return Server{
		log:                 log,
		eventBookRepository: eventBookRepository,
		transports:          transports,
		businessClient:      businessClient,
	}
}

func (o *Server) Shutdown() {
	o.server.GracefulStop()
}

type Server struct {
	business.UnimplementedBusinessCoordinatorServer
	log                 *zap.SugaredLogger
	eventBookRepository eventBook.Storer
	transports          transport.TransportHolder
	businessClient      client.BusinessClient
	server              *grpc.Server
}

func (o Server) Handle(ctx context.Context, in *eventedcore.CommandBook) (result *eventedcore.SynchronousProcessingResponse, err error) {
	uuid, err := eventedproto.ProtoToUUID(in.Cover.Root)
	if err != nil {
		return nil, err
	}
	priorState, err := o.eventBookRepository.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}

	contextualCommand := &eventedcore.ContextualCommand{
		Events:  priorState,
		Command: in,
	}

	var businessResponse *eventedcore.EventBook
	err = backoff.Retry(func() error {
		businessResponse, err = o.businessClient.Handle(ctx, contextualCommand)
		return err
	}, backoff.NewExponentialBackOff())
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

func (o Server) handleEventBook(ctx context.Context, eb *eventedcore.EventBook) (result *eventedcore.SynchronousProcessingResponse, rerr error) {
	result = &eventedcore.SynchronousProcessingResponse{}
	result.Books = []*eventedcore.EventBook{eb}

	err := o.eventBookRepository.Put(ctx, eb)
	if err != nil {
		return nil, err
	}

	sync, _, err := o.extractSynchronous(eb)
	if err != nil {
		return nil, err
	}

	syncResult, err := o.executeSyncProjections(ctx, sync)
	if err != nil {
		return result, err
	}
	result.Projections = append(result.Projections, syncResult...)

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

func (o Server) executeSyncSagas(ctx context.Context, sync *eventedcore.EventBook) (eventBooks []*eventedcore.EventBook, projections []*eventedcore.Projection, rerr error) {
	for _, syncSaga := range o.transports.GetSaga() {
		var response *eventedcore.SynchronousProcessingResponse
		var err error
		backoff.Retry(func() error {
			response, err = syncSaga.HandleSync(ctx, sync)
			return err
		}, backoff.NewExponentialBackOff())
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

func (o Server) executeSyncProjections(ctx context.Context, sync *eventedcore.EventBook) (result []*eventedcore.Projection, rerr error) {
	result = []*eventedcore.Projection{}
	for _, syncProjector := range o.transports.GetProjectors() {
		var response *eventedcore.Projection
		var err error
		backoff.Retry(func() error {
			response, err = syncProjector.HandleSync(ctx, sync)
			return err
		}, backoff.NewExponentialBackOff())
		if err != nil {
			o.log.Error(err)
			rerr = multierror.Append(rerr, err)
			continue
		}
		result = append(result, response)
	}
	return result, rerr
}

func (o Server) extractSynchronous(originalBook *eventedcore.EventBook) (synchronous *eventedcore.EventBook, async *eventedcore.EventBook, err error) {
	if len(originalBook.Pages) == 0 {
		return nil, nil, errors.New("event book has no pages -- not correct in this context")
	}
	var lastIdx uint32
	for idx, event := range originalBook.Pages {
		if event.Synchronous {
			lastIdx = uint32(idx) + 1
		}
	}
	synchronous = new(eventedcore.EventBook)
	synchronous.Pages = originalBook.Pages[:lastIdx]
	synchronous.Cover = originalBook.Cover
	synchronous.Snapshot = originalBook.Snapshot

	async = new(eventedcore.EventBook)
	async.Pages = originalBook.Pages[lastIdx:]
	async.Cover = originalBook.Cover
	async.Snapshot = nil

	return synchronous, async, err
}

func (o Server) Record(ctx context.Context, in *eventedcore.EventBook) (response *eventedcore.SynchronousProcessingResponse, err error) {
	return o.handleEventBook(ctx, in)
}

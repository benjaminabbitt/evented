package framework

import (
	"context"
	"errors"
	"fmt"
	"github.com/benjaminabbitt/evented/applications/commandHandler/business/client"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework/transport"
	eventedproto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func NewServer(eventBookRepository eventBook.EventBookStorer, transports transport.TransportHolder, businessClient client.BusinessClient, log *zap.SugaredLogger) Server {
	return Server{
		log:                 log,
		eventBookRepository: eventBookRepository,
		transports:          transports,
		businessClient:      businessClient,
	}
}

func (o *Server) Listen(port uint) error {
	o.log.Infow("Opening port", "port", port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		o.log.Error(err)
		return err
	}
	o.log.Infow("Creating GRPC Server")
	o.server = grpc.NewServer()
	o.log.Infow("Registering Command Handler with GRPC")
	evented_core.RegisterCommandHandlerServer(o.server, o)
	o.log.Infow("Handler registered.")
	o.log.Infow("Serving...")
	err = o.server.Serve(lis)
	if err != nil {
		o.log.Error(err)
		return err
	}
	return nil
}

func (o *Server) Earmuffs() {
	o.server.GracefulStop()
}

type Server struct {
	evented_core.UnimplementedCommandHandlerServer
	log                 *zap.SugaredLogger
	eventBookRepository eventBook.EventBookStorer
	transports          transport.TransportHolder
	businessClient      client.BusinessClient
	server              *grpc.Server
}

func (o Server) Handle(ctx context.Context, in *evented_core.CommandBook) (result *evented_core.SynchronousProcessingResponse, err error) {
	uuid, err := eventedproto.ProtoToUUID(in.Cover.Root)
	if err != nil {
		return nil, err
	}
	priorState, err := o.eventBookRepository.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}

	contextualCommand := &evented_core.ContextualCommand{
		Events:  priorState,
		Command: in,
	}

	businessResponse, err := o.businessClient.Handle(ctx, contextualCommand)
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

func (o Server) handleEventBook(ctx context.Context, eb *evented_core.EventBook) (result *evented_core.SynchronousProcessingResponse, rerr error) {
	result = &evented_core.SynchronousProcessingResponse{}
	result.Books = []*evented_core.EventBook{eb}

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

func (o Server) executeSyncSagas(ctx context.Context, sync *evented_core.EventBook) (eventBooks []*evented_core.EventBook, projections []*evented_core.Projection, rerr error) {
	for _, syncSaga := range o.transports.GetSaga() {
		response, err := syncSaga.HandleSync(ctx, sync)
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

func (o Server) executeSyncProjections(ctx context.Context, sync *evented_core.EventBook) (result []*evented_core.Projection, rerr error) {
	result = []*evented_core.Projection{}
	for _, syncProjector := range o.transports.GetProjectors() {
		response, err := syncProjector.HandleSync(ctx, sync)
		if err != nil {
			o.log.Error(err)
			rerr = multierror.Append(rerr, err)
			continue
		}
		result = append(result, response)
	}
	return result, rerr
}

func (o Server) extractSynchronous(originalBook *evented_core.EventBook) (synchronous *evented_core.EventBook, async *evented_core.EventBook, err error) {
	if len(originalBook.Pages) == 0 {
		return nil, nil, errors.New("event book has no pages -- not correct in this context")
	}
	var lastIdx uint32
	for idx, event := range originalBook.Pages {
		if event.Synchronous {
			lastIdx = uint32(idx) + 1
		}
	}
	synchronous = new(evented_core.EventBook)
	synchronous.Pages = originalBook.Pages[:lastIdx]
	synchronous.Cover = originalBook.Cover
	synchronous.Snapshot = originalBook.Snapshot

	async = new(evented_core.EventBook)
	async.Pages = originalBook.Pages[lastIdx:]
	async.Cover = originalBook.Cover
	async.Snapshot = nil

	return synchronous, async, err
}

func (o Server) Record(ctx context.Context, in *evented_core.EventBook) (response *evented_core.SynchronousProcessingResponse, err error) {
	return o.handleEventBook(ctx, in)
}

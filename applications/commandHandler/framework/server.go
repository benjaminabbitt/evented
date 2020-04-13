package framework

import (
	"context"
	"fmt"
	"github.com/benjaminabbitt/evented/applications/commandHandler/business/client"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework/transport"
	eventedproto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/repository/eventBook"
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

func (o *Server) Listen(port uint16) {
	o.log.Infow("Opening port", "port", port)
	lis := o.createListener(port)
	o.log.Infow("Port opened", "port", port)
	o.log.Infow("Creating GRPC Server")
	grpcServer := grpc.NewServer()
	o.log.Infow("Registering Command Handler with GRPC")
	evented_core.RegisterCommandHandlerServer(grpcServer, o)
	o.log.Infow("Handler registered.")
	o.log.Infow("Serving...")
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
}

func (o *Server) createListener(port uint16) net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		o.log.Error(err)
	}
	return lis
}

type Server struct {
	evented_core.UnimplementedCommandHandlerServer
	log                 *zap.SugaredLogger
	eventBookRepository eventBook.EventBookStorer
	transports          transport.TransportHolder
	businessClient      client.BusinessClient
}

func (o Server) Handle(ctx context.Context, in *evented_core.CommandBook) (result *evented_core.CommandHandlerResponse, err error) {
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
	response, err := o.handleEventBook(ctx, businessResponse)
	return response, err
}

func (o Server) handleEventBook(ctx context.Context, eb *evented_core.EventBook) (result *evented_core.CommandHandlerResponse, err error) {
	err = o.eventBookRepository.Put(ctx, eb)

	sync, _ := o.extractSynchronous(eb)
	var eventBooks []*evented_core.EventBook
	var projections []*evented_core.Projection

	for _, syncProjector := range o.transports.GetProjections() {
		response, err := syncProjector.HandleSync(ctx, sync)
		if err != nil {
			o.log.Error(err)
			return nil, err
		}
		projections = append(projections, response)
	}

	for _, syncSaga := range o.transports.GetSaga() {
		response, err := syncSaga.HandleSync(ctx, sync)
		if err != nil {
			o.log.Error(err)
			return nil, err
		}

		eventBooks = append(eventBooks, response)
	}

	for _, t := range o.transports.GetTransports() {
		err := t.Handle(ctx, eb)
		if err != nil {
			o.log.Error(err)
			return nil, err
		}
	}

	return result, nil
}

func (o Server) extractSynchronous(originalBook *evented_core.EventBook) (synchronous *evented_core.EventBook, async *evented_core.EventBook) {
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

	return synchronous, async
}

func (o Server) Record(ctx context.Context, in *evented_core.EventBook) (response *evented_core.CommandHandlerResponse, err error) {
	r, err := o.handleEventBook(ctx, in)
	return r, err
}

package framework

import (
	"context"
	"flag"
	"fmt"
	"github.com/benjaminabbitt/evented/proto/business"
	"github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	"github.com/benjaminabbitt/evented/transport"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func NewServer(eventBookRepository eventBook.Repository, syncSagas []transport.SyncSaga, syncProjections []transport.SyncProjection, asyncSagas []transport.Saga, asyncProjections []transport.Projection, businessClient evented_business.BusinessLogicClient) Server {
	return Server{
		eventBookRepository: eventBookRepository,
		syncSagas:           syncSagas,
		syncProjectors:      syncProjections,
		asyncSagas:          asyncSagas,
		asyncProjectors:     asyncProjections,
		businessClient:      businessClient,
	}
}

func Listen(server Server, port uint16){
	lis := createListener(port)
	grpcServer := grpc.NewServer()

	evented_core.RegisterCommandHandlerServer(grpcServer, server)
	err := grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateBusinessClient(address string) (evented_business.BusinessLogicClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := evented_business.NewBusinessLogicClient(conn)
	return client, conn
}

func createListener(port uint16) net.Listener {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return lis
}

type Server struct {
	evented_core.UnimplementedCommandHandlerServer
	eventBookRepository eventBook.Repository
	syncSagas           []transport.SyncSaga
	syncProjectors      []transport.SyncProjection
	asyncSagas          []transport.Saga
	asyncProjectors     []transport.Projection
	businessClient      evented_business.BusinessLogicClient
}

func (s Server) Handle(ctx context.Context, in *evented_core.CommandBook) (result *evented_core.CommandHandlerResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	priorState, err := s.eventBookRepository.Get(in.Cover.Root)

	defer cancel()

	contextualCommand := &evented_core.ContextualCommand{
		Events:  &priorState,
		Command: in,
	}

	businessResponse, err := s.businessClient.Handle(ctx, contextualCommand)
	if err != nil {
		log.Fatal(err)
	}

	err = s.eventBookRepository.Put(*businessResponse)
	if err != nil {
		log.Fatal(err)
	}

	sync, async := s.extractSynchronous(*businessResponse)

	for _, saga := range s.syncSagas{
		book, _ := saga.SendSync(&sync)
		result.Books = append(result.Books, book)
	}

	for _, projector := range s.syncProjectors {
		projection, _ := projector.ProjectSync(&sync)
		result.Projections = append(result.Projections, projection)
	}

	for _, saga := range s.syncSagas {
		saga.SendSync(&async)
	}

	for _, projector := range s.syncProjectors {
		projector.ProjectSync(&async)
	}

	for _, saga := range s.asyncSagas {
		saga.Send(businessResponse)
	}

	for _, projector := range s.asyncProjectors {
		projector.Project(businessResponse)
	}
	return nil, nil
}


func (s Server) extractSynchronous(originalBook evented_core.EventBook) (synchronous evented_core.EventBook, async evented_core.EventBook) {
	var lastIdx uint32
	for idx, event := range originalBook.Pages {
		if event.Synchronous {
			lastIdx = uint32(idx) + 1
		}
	}
	synchronous.Pages = originalBook.Pages[:lastIdx]
	synchronous.Cover = originalBook.Cover
	synchronous.Snapshot = originalBook.Snapshot

	async.Pages = originalBook.Pages[lastIdx:]
	async.Cover = originalBook.Cover
	async.Snapshot = nil

	return synchronous, async
}

func (s Server) Record(ctx context.Context, in *evented_core.EventBook) (*evented_core.CommandHandlerResponse, error) {
	return nil, nil
}

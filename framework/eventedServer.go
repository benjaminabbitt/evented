package framework

import (
	"context"
	"flag"
	"fmt"
	"github.com/benjaminabbitt/evented/protobuf"
	"github.com/benjaminabbitt/evented/repository"
	"github.com/benjaminabbitt/evented/transport"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func NewServer(eventRepos repository.EventRepository, transport transport.EventTransportSender, port uint16, businessAddress string) {
	lis := createListener(port)
	client, conn := createBusinessClient(businessAddress)
	defer conn.Close()
	grpcServer := grpc.NewServer()
	protobuf.RegisterCommandHandlerServer(grpcServer, &server{
		eventRepository:      eventRepos,
		eventTransportSender: transport,
		businessClient: client,
	})
	grpcServer.Serve(lis)
}

func createBusinessClient(address string) (protobuf.BusinessLogicClient, *grpc.ClientConn){
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := protobuf.NewBusinessLogicClient(conn)
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


type server struct{
	protobuf.UnimplementedCommandHandlerServer
	eventRepository repository.EventRepository
	eventTransportSender transport.EventTransportSender
	businessClient protobuf.BusinessLogicClient
}

func (s *server) Handle(ctx context.Context, in *protobuf.Command)(*protobuf.Events, error){
	s.eventRepository.Get(in.Id)

	log.Print("In Handle...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	pc := &protobuf.BusinessCommand{
		Snapshot: nil,
		Events:   nil,
		Command:  in,
	}

	businessResponse, err := s.businessClient.Handle(ctx, pc)

	log.Printf("%+v", businessResponse)

	//s.eventRepository.Add()

	return businessResponse.Events, err
}

func (s *server) Record(ctx context.Context, in *protobuf.Events)(*protobuf.Empty, error){
	return nil, nil
}

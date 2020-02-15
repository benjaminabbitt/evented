package framework

import (
	"context"
	"flag"
	"fmt"
	"github.com/benjaminabbitt/evented/proto/business"
	"github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/transport"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func NewServer(eventRepos EventRepository, transport transport.EventTransportSender, port uint16, businessAddress string) {
	lis := createListener(port)
	client, conn := createBusinessClient(businessAddress)
	defer conn.Close()
	grpcServer := grpc.NewServer()
	evented_core.RegisterCommandHandlerServer(grpcServer, &server{
		eventRepository:      eventRepos,
		eventTransportSender: transport,
		businessClient: client,
	})
	grpcServer.Serve(lis)
}

func createBusinessClient(address string) (evented_business.BusinessLogicClient, *grpc.ClientConn){
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


type server struct{
	evented_core.UnimplementedCommandHandlerServer
	eventRepository      EventRepository
	eventTransportSender transport.EventTransportSender
	businessClient       evented_business.BusinessLogicClient
}

func (s *server) Handle(ctx context.Context, in *evented_core.CommandBook)(*evented_core.EventBook, error){
	commands := CommandBookToCommand(in)
	for _, command := range commands {
		_, err := s.eventRepository.Get(command.Id)
		if(err != nil){

		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	pc := &evented_business.BusinessCommand{
		Snapshot: nil,
		Events:   nil,
		Command:  in,
	}

	businessResponse, err := s.businessClient.Handle(ctx, pc)

	log.Printf("%+v", businessResponse)

	//s.eventRepository.Add()

	return businessResponse.Events, err
}


func (s *server) Record(ctx context.Context, in *evented_core.EventBook)(*evented_core.Empty, error){
	return nil, nil
}

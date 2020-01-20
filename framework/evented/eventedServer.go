package main

import (
	"awesomeProject/framework/generated/pb/evented"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

const (
	address     = "localhost:8081"
)


func main() {
	flag.Parse()
	port := 8080
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	evented.RegisterCommandHandlerServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}

type server struct{
	evented.UnimplementedCommandHandlerServer
}

func (s *server) Handle(ctx context.Context, in *evented.Command)(*evented.Events, error){
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := evented.NewBusinessLogicClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	pc := &evented.BusinessCommand{
		Snapshot: nil,
		Events:   nil,
		Command:  in,
	}

	businssResponse, err := c.Handle(ctx, pc)
	return businssResponse.Events, err
}

func (s *server) Record(ctx context.Context, in *evented.Events)(*evented.Empty, error){

	return nil, nil
}

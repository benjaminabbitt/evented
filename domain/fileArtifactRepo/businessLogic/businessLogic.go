package businessLogic

//go:generate protoc --go_out=./generated/pb/ --proto_path=./ fileArtifactRepo.proto

import (
	far "awesomeProject/domain/fileArtifactRepo/generated/pb"
	evented "awesomeProject/framework/generated/pb"
	"context"
	"flag"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	flag.Parse()
	port := 8081
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	evented.RegisterBusinessLogicServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}

type server struct{
	evented.UnimplementedBusinessLogicServer
}

func (s *server) Handle(ctx context.Context, in *evented.BusinessCommand)(*evented.BusinessResponse, error){
	commandDetails := far.CreateArtifact{}
	ptypes.UnmarshalAny(in.Command.CommandSpecifics, &commandDetails)
	ac, _ := ptypes.MarshalAny(&far.ArtifactCreated{
		Version:              "0",
		Name:                 commandDetails.Name,
		Url:                  commandDetails.Url,
		Description:          commandDetails.Description,
	})

	event := []*evented.Event{&evented.Event{
		Id:     in.Command.Id,
		User:   in.Command.User,
		Ts:     time.Now().UTC().Format(time.RFC3339),
		Bodies: []*evented.EventSequence{
			&evented.EventSequence{
				Sequence: 0,
				EventDetails: ac,
			},
		},
	}}


	events := evented.Events{
		Events: event,
	}

	response := evented.BusinessResponse{
		Snapshot: nil,
		Events:   &events,
	}

	return &response, nil
}


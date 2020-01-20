package main

import (
	far "awesomeProject/domain/fileArtifactRepo/generated/pb"
	evented "awesomeProject/framework/generated/pb"
	"context"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address     = "localhost:8080"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	createArtifactAny, err := ptypes.MarshalAny(&far.CreateArtifact{
		Version:     "",
		Name:        "",
		Url:         "",
		Description: "",
	})
	client := evented.NewCommandHandlerClient(conn)
	r, err := client.Handle(ctx, &evented.Command{
		Id:               "id",
		User:             "user",
		CommandSpecifics: createArtifactAny,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Events[0].Id)
}

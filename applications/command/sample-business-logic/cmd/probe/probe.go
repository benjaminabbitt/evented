package probe

import (
	"context"
	"flag"
	"github.com/benjaminabbitt/evented/applications/command/sample-business-logic/support"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/todo/proto"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
)

func init() {
	Probe.AddCommand(Create)
	Probe.AddCommand(Edit)
	Probe.AddCommand(Complete)
	Probe.AddCommand(UnComplete)
}

var Probe = &cobra.Command{
	Use:   "probe",
	Short: "",
	Run:   nil,
}

var Create = &cobra.Command{
	Use:   "create",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		println("In Create")
		body := "Body"
		due := timestamppb.New(time.Now())
		td, err := time.ParseDuration("1h")
		duration := durationpb.New(td)
		id, err := uuid.NewRandom()
		mid, err := id.MarshalBinary()

		createTodo := &proto.CreateTodo{
			Title:     "Title",
			Body:      &body,
			Due:       due,
			Duration:  duration,
			Important: BoolPtr(false),
			RemindAt:  nil,
		}
		marshalled, err := anypb.New(createTodo)
		if err != nil {
			println(err)
		}
		cp := []*evented.CommandPage{{
			Sequence:    0,
			Synchronous: false,
			Command:     marshalled,
		}}
		cc := &evented.Cover{
			Domain: support.Domain,
			Root:   &evented.UUID{Value: mid},
		}
		cb := &evented.CommandBook{
			Cover: cc,
			Pages: cp,
		}
		Send(cb)
	},
}

var Edit = &cobra.Command{
	Use:   "edit",
	Short: "",
	Run:   nil,
}

var Complete = &cobra.Command{
	Use:   "complete",
	Short: "",
	Run:   nil,
}

var UnComplete = &cobra.Command{
	Use:   "uncomplete",
	Short: "",
	Run:   nil,
}

var (
	addr = "localhost:1313"
)

func Send(book *evented.CommandBook) {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("fail to close: %v", err)
		}
	}(conn)
	client := evented.NewBusinessCoordinatorClient(conn)
	client.Handle(context.Background(), book)
}

func BoolPtr(it bool) *bool {
	return &it
}

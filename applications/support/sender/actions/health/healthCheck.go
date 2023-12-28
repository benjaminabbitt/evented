package health

import (
	"context"
	"fmt"
	"github.com/benjaminabbitt/evented/applications/support/sender/actions/root"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"strconv"
)

func init() {
	root.RootCmd.AddCommand(sendHealth)
}

var sendHealth = &cobra.Command{
	Use:   "health",
	Short: "Sends a GRPC Health Check to the specified host, port, and service name",
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		port, _ := strconv.Atoi(args[1])
		name := args[2]
		conn, _ := grpc.Dial(fmt.Sprintf("%s:%v", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		defer conn.Close()
		client := grpc_health_v1.NewHealthClient(conn)
		req := &grpc_health_v1.HealthCheckRequest{Service: name}
		resp, _ := client.Check(context.Background(), req)
		fmt.Println(resp)
	},
}

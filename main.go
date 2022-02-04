package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/scardozos/esplai-planning/grpc/groups"
	"github.com/scardozos/esplai-planning/models"
	"google.golang.org/grpc"
)

var (
	grpcServerAddr = os.Getenv("GRPC_SERVER_CLIENT_ADDR_SELF")
	grpcServerPort = os.Getenv("GRPC_SERVER_CLIENT_PORT_SELF")
)

func main() {

	// Server logic
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", grpcServerAddr, grpcServerPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	groupServer := models.NewGroupServer()
	pb.RegisterGroupsServer(grpcServer, groupServer)

	/* Testing:
	d := groupServer.dbClient
	go d.UnsetNonWeek(&models.DateTime{Year: 2022, Month: 2, Day: 5})
	go d.UnsetNonWeek(&models.DateTime{Year: 2022, Month: 1, Day: 23})
	go d.IsNonWeek(&models.DateTime{Year: 2022, Month: 1, Day: 22})
	go d.IsNonWeek(&models.DateTime{Year: 2022, Month: 1, Day: 23})
	go d.IsNonWeek(&models.DateTime{Year: 2022, Month: 1, Day: 24})
	go d.IsNonWeek(&models.DateTime{Year: 2022, Month: 1, Day: 25})
	*/

	grpcServer.Serve(lis)

}

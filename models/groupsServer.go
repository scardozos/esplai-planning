package models

import (
	"context"
	"log"
	"time"

	pb "github.com/scardozos/esplai-planning/grpc/groups"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	startYear  = 2022
	startMonth = 1
	startDay   = 17
)

type GroupsServer struct {
	dbClient *GrpcClient
	pb.UnimplementedGroupsServer
}

// TODO: Refactor code
func (s *GroupsServer) GetGroupPlaces(ctx context.Context, dateRequest *pb.DateRequest) (*pb.GroupsPlacesResponse, error) {
	date := dateRequest.Date
	log.Printf("Got new request for %v-%v-%v\n", date.Year, date.Month, date.Day)
	startDate := DateTime{Year: startYear, Month: startMonth, Day: startDay}
	requestedDate := DateTime{Year: date.Year, Month: date.Month, Day: date.Day}

	nonWeeks, err := s.dbClient.GetNonWeeks()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get static weeks: %v", err)
	}

	groups := InitialGroupState()

	reqDateUnmarshaled := requestedDate.ToTime()
	saturday := ChangeWeekDay(reqDateUnmarshaled, time.Saturday)

	futureGroups := IterateNextWeeks(
		CalcWeekNumNoWeeks(startDate.ToTime(), reqDateUnmarshaled, nonWeeks),
		groups,
	)
	futureGroupsApiModel := MarshalGroupModel(futureGroups)
	return &pb.GroupsPlacesResponse{
		Groups:            futureGroupsApiModel,
		RequestedSaturday: &pb.Date{Year: int32(saturday.Year()), Month: int32(saturday.Month()), Day: int32(saturday.Day())},
	}, nil
}

func NewGroupServer() *GroupsServer {
	// Initialize database client config containing static weeks
	clientCtx, err := NewGrpcClientContext("localhost:9001")
	if err != nil {
		log.Fatal(err)
	}
	db := &GrpcClient{Context: clientCtx}
	return &GroupsServer{dbClient: db}
}

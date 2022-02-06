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
	// caching NonWeeks allow for a great alleviation of incoming traffic towards dbServer (ep-weekhandler)
	//
	// Rationale
	//
	// NonWeeks will (should not) be changed regularly
	// and should also fit in memory, for which caching allows for
	// an improved backend performance
	cachedNonWeeks []time.Time
	// we declare dbClient for refreshing cachedNonWeeks
	// TODO: Might add future option to disable caching
	dbClient *GrpcClient
	pb.UnimplementedGroupsServer
}

// TODO: Refactor code
func (s *GroupsServer) GetGroupPlaces(ctx context.Context, dateRequest *pb.DateRequest) (*pb.GroupsPlacesResponse, error) {
	date := dateRequest.Date
	log.Printf("Got new request for %v-%v-%v\n", date.Year, date.Month, date.Day)
	startDate := DateTime{Year: startYear, Month: startMonth, Day: startDay}
	requestedDate := DateTime{Year: date.Year, Month: date.Month, Day: date.Day}

	// get cached nonWeeks
	// if 0 are found get from gRPC stub
	if len(s.cachedNonWeeks) == 0 {
		resWeeks, err := s.dbClient.GetNonWeeks()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not get static weeks: %v", err)
		}
		// caching upon user execution if initial request fails
		s.cachedNonWeeks = resWeeks
	}
	nonWeeks := s.cachedNonWeeks

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

	nonWeeks, err := db.GetNonWeeks()
	if err != nil {
		// if couldn't get nonWeeks only return dbClient in order to
		// cache NonWeeks in future calls (in case dbServer is currently unavailable)
		return &GroupsServer{dbClient: db}
	}

	// if no error, return db client and cachedNonWeeks
	return &GroupsServer{
		dbClient:       db,
		cachedNonWeeks: nonWeeks,
	}
}

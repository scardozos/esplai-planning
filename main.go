package main

import (
	"context"
	"log"
	"net"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	wh "github.com/scardozos/ep-weekhandler/grpc/dates"
	pb "github.com/scardozos/esplai-planning/grpc/groups"
	"github.com/scardozos/esplai-planning/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	startYear  = 2022
	startMonth = 1
	startDay   = 17
)

// Day must always be a saturday
/*
var nonWeeks = []time.Time{
	(&models.DateTime{Year: 2022, Month: 1, Day: 22}).ToTime(),
	(&models.DateTime{Year: 2022, Month: 2, Day: 5}).ToTime(),
	(&models.DateTime{Year: 2022, Month: 3, Day: 19}).ToTime(),
}
*/

type GroupsServer struct {
	dbClient *models.GrpcClient
	pb.UnimplementedGroupsServer
}

func newGrpcClientContext(endpoint string) (*models.GrpcClientContext, error) {
	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(10 * time.Millisecond)),
		grpc_retry.WithPerRetryTimeout(20 * time.Millisecond),
	}
	datesConn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
	)
	if err != nil {
		return nil, err
	}

	return &models.GrpcClientContext{DatesClient: wh.NewDatesClient(datesConn)}, nil

}

// TODO: Clean up this mess
func (s *GroupsServer) GetGroupPlaces(ctx context.Context, dateRequest *pb.DateRequest) (*pb.GroupsPlacesResponse, error) {
	date := dateRequest.Date
	log.Printf("Got new request for %v-%v-%v\n", date.Year, date.Month, date.Day)
	startDate := models.DateTime{Year: startYear, Month: startMonth, Day: startDay}
	requestedDate := models.DateTime{Year: date.Year, Month: date.Month, Day: date.Day}

	nonWeeks, err := s.dbClient.GetNonWeeks()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get static weeks: %v", err)
	}

	iterNum := CalcWeekNumNoWeeks(startDate.ToTime(), requestedDate.ToTime(), nonWeeks)
	groups := InitialGroupState()

	reqDateUnmarshaled := requestedDate.ToTime()
	saturday := ChangeWeekDay(reqDateUnmarshaled, time.Saturday)

	futureGroups := IterateNextWeeks(iterNum, groups)
	futureGroupsApiModel := MarshalGroupModel(futureGroups)

	return &pb.GroupsPlacesResponse{
		Groups:            futureGroupsApiModel,
		RequestedSaturday: &pb.Date{Year: int32(saturday.Year()), Month: int32(saturday.Month()), Day: int32(saturday.Day())},
	}, nil
}

func newGroupServer() *GroupsServer {
	// Initialize database client config containing static weeks
	clientCtx, err := newGrpcClientContext("localhost:9001")
	if err != nil {
		log.Fatal(err)
	}
	db := &models.GrpcClient{Context: clientCtx}
	return &GroupsServer{dbClient: db}
}

func main() {

	// Server logic
	lis, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	groupServer := newGroupServer()
	pb.RegisterGroupsServer(grpcServer, groupServer)

	/* Testing:
	d := groupServer.dbClient
	go d.UnsetNonWeek(&models.DateTime{Year: 2022, Month: 1, Day: 23})
	go d.AddStaticDate(&models.DateTime{Year: 2022, Month: 1, Day: 23})
	go d.IsNonWeek(&models.DateTime{Year: 2022, Month: 1, Day: 22})
	go d.IsNonWeek(&models.DateTime{Year: 2022, Month: 1, Day: 23})
	go d.IsNonWeek(&models.DateTime{Year: 2022, Month: 1, Day: 24})
	go d.IsNonWeek(&models.DateTime{Year: 2022, Month: 1, Day: 25})
	go d.UnsetNonWeek(&models.DateTime{Year: 2022, Month: 1, Day: 23})
	*/

	grpcServer.Serve(lis)

}

// Returns a slice with a list of groups and their places
// Initial state:
// Teatre | Parc Central |  Pista    | Plaça | Passarel·la
//   <-          <-           <-        <-         <-
// Aneto  |  Pedraforca	 | Matagalls | Cadí  | Puigmal
func InitialGroupState() models.GroupsList {
	// Initialize places
	passarela := &models.Place{Name: "Pasarel·la", Next: nil}
	plaza := &models.Place{Name: "Plaça", Next: passarela}
	pista := &models.Place{Name: "Pista", Next: plaza}
	parcCentral := &models.Place{Name: "Parc Central", Next: pista}
	teatre := &models.Place{Name: "Teatre", Next: parcCentral}
	passarela.Next = teatre

	// Start groups and assign places
	return models.GroupsList{
		GroupsList: []*models.Group{
			{Name: "Aneto", Place: teatre},
			{Name: "Pedraforca", Place: parcCentral},
			{Name: "Matagalls", Place: pista},
			{Name: "Cadí", Place: plaza},
			{Name: "Puigmal", Place: passarela},
		},
	}
}

// Translates local group logic declarations to protobuf format
func MarshalGroupModel(groups models.GroupsList) []*pb.Group {
	var groupApiModel = make([]*pb.Group, len(groups.GroupsList))
	for index, group := range groups.GroupsList {
		groupApiModel[index] = &pb.Group{
			GroupName: group.Name,
			GroupPlace: &pb.Place{
				PlaceName: group.Place.Name,
			},
		}
	}
	//log.Println(groupApiModel)
	return groupApiModel
}

// Gets the state for x total number of weeks, taking into account the groups
func IterateNextWeeks(weeks int, groups models.GroupsList) models.GroupsList {
	for i := 0; i < weeks; i++ {
		groups.NextIteration()
	}
	return groups
}

// Calculates number of weeks for which to know their respective state in the future
// Takes into account startDate, the requested date and the list of days in which state won't change
func CalcWeekNumNoWeeks(startDate time.Time, requestedDate time.Time, nonWeeks []time.Time) int {
	// Change requestedDate to that week's Monday in order to preserve state
	requestedDate = ChangeWeekDay(requestedDate, time.Monday)

	// Compute the amount of days from startDate
	days := requestedDate.Sub(startDate).Hours() / 24
	// Compute the amount of weeks from startDate
	// taking into account the amount of days
	weeks := int(days / 7)

	// initialize sub var, which accounts for the total number
	// of static weeks that have occurred after startDate
	// and before the requestedDate
	var sub int
	for _, time := range nonWeeks {
		if time.After(startDate) && time.Before(requestedDate) {
			sub += 1
		}
	}

	// returns the number of weeks that passed from startDate
	// minus the amount of static weeks that passed in between
	return weeks - sub
}

// Change week day "from" `Time.Time` to a given weekday "to" `time.Weekday`
// Returns time.Time
func ChangeWeekDay(from time.Time, to time.Weekday) time.Time {
	if currentWeekDay := int(from.Weekday()); currentWeekDay != int(to) {
		sub := currentWeekDay
		if currentWeekDay == 0 {
			sub = 7
		}
		return from.AddDate(0, 0, int(to)-sub)
	}
	return from
}

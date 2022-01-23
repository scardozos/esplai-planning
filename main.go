package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/scardozos/esplai-planning/grpc/groups"
	"github.com/scardozos/esplai-planning/models"
	"google.golang.org/grpc"
)

const (
	startYear  = 2022
	startMonth = 1
	startDay   = 17
)

// TODO: Delegate the storage and retrieval of nonWeeks to another gRPC endpoint
// Day must always be a saturday
var nonWeeks = []time.Time{
	(&models.DateTime{Year: 2022, Month: 1, Day: 22}).ToTime(),
	(&models.DateTime{Year: 2022, Month: 2, Day: 5}).ToTime(),
	(&models.DateTime{Year: 2022, Month: 3, Day: 19}).ToTime(),
}

type GroupsServer struct {
	pb.UnimplementedGroupsServer
}

func (s *GroupsServer) GetGroupPlaces(ctx context.Context, dateRequest *pb.DateRequest) (*pb.GroupsPlacesResponse, error) {
	date := dateRequest.Date
	log.Printf("Got new request for %v-%v-%v\n", date.Year, date.Month, date.Day)
	startDate := models.DateTime{Year: startYear, Month: startMonth, Day: startDay}
	requestedDate := models.DateTime{Year: date.Year, Month: date.Month, Day: date.Day}

	iterNum := CalcWeekNumNoWeeks(startDate.ToTime(), requestedDate.ToTime(), nonWeeks)
	groups := InitialGroupState()

	reqDateUnmarshaled := requestedDate.ToTime()
	subSat := int(reqDateUnmarshaled.Weekday())
	if reqDateUnmarshaled.Weekday() == 0 {
		subSat = 7
	}
	saturday := reqDateUnmarshaled.AddDate(0, 0, 6-subSat)

	futureGroups := IterateNextWeeks(iterNum, groups)
	futureGroupsApiModel := MarshalGroupModel(futureGroups)

	return &pb.GroupsPlacesResponse{
		Groups:            futureGroupsApiModel,
		RequestedSaturday: &pb.Date{Year: int32(saturday.Year()), Month: int32(saturday.Month()), Day: int32(saturday.Day())},
	}, nil
}

func newGroupServer() *GroupsServer {
	s := &GroupsServer{}
	return s
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGroupsServer(grpcServer, newGroupServer())
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
	// Convert any non-Monday date to Monday
	// This preserves state for the entire week, as if it were Monday
	if reqWeekDay := int(requestedDate.Weekday()); reqWeekDay != 5 {
		sub := reqWeekDay
		if reqWeekDay == 0 {
			sub = 7
		}
		requestedDate = requestedDate.AddDate(0, 0, 5-sub)
	}
	// Calculate number of weeks since startDate
	fmt.Println(requestedDate.Weekday())
	days := requestedDate.Sub(startDate).Hours() / 24
	weeks := int(days / 7)

	// Compute total number of nonWeek occurrences since startDate, that happen before the requestedDate
	var sub int
	for _, time := range nonWeeks {
		if time.After(startDate) && time.Before(requestedDate) {
			sub += 1
		}
	}
	return weeks - sub
}

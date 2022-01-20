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

type GroupsServer struct {
	pb.UnimplementedGroupsServer
}

func (s *GroupsServer) GetGroupPlaces(ctx context.Context, date *pb.DateRequest) (*pb.GroupsPlacesResponse, error) {
	log.Printf("Got new request for %v-%v-%v\n", date.Year, date.Month, date.Day)
	startDate := models.DateTime{Year: startYear, Month: startMonth, Day: startDay}
	endDate := models.DateTime{Year: date.Year, Month: date.Month, Day: date.Day}

	weekNum := CalcWeekNum(startDate.ToTime(), endDate.ToTime())
	groups := GetInitialState()

	futureGroups := IterateNextWeeks(weekNum, groups)
	futureGroupsApiModel := LocalGroupModelToApi(futureGroups)
	return &pb.GroupsPlacesResponse{
		Groups: futureGroupsApiModel,
	}, nil
}

func newServer() *GroupsServer {
	s := &GroupsServer{}
	return s
}

func main() {

	lis, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGroupsServer(grpcServer, newServer())
	grpcServer.Serve(lis)

}

// Returns a slice with a list of groups and their places
// Initial state:
// Teatre | Parc Central |  Pista    | Plaça | Passarel·la
//   <-          <-           <-        <-         <-
// Aneto  |  Pedraforca	 | Matagalls | Cadí  | Puigmal
func GetInitialState() models.GroupsList {

	// Initialize places
	passarela := &models.Place{Name: "Pasarel·la", Next: nil}
	plaza := &models.Place{Name: "Plaça", Next: passarela}
	pista := &models.Place{Name: "Pista", Next: plaza}
	parcCentral := &models.Place{Name: "Parc Central", Next: pista}
	teatre := &models.Place{Name: "Teatre", Next: parcCentral}
	passarela.Next = teatre

	/**
	ant := &models.Group{Name: "Aneto", Place: teatre}
	pdf := &models.Group{Name: "Pedraforca", Place: parcCentral}
	mtg := &models.Group{Name: "Matagalls", Place: pista}
	cdi := &models.Group{Name: "Cadí", Place: plaza}
	pgm := &models.Group{Name: "Puigmal", Place: passarela}
	groups := models.GroupsList{GroupsList: []*models.Group{ant, pdf, mtg, cdi, pgm}}
	return groups
	**/
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
func LocalGroupModelToApi(groups models.GroupsList) []*pb.Group {
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

// Prints the groups and their respective place
func ShowGroupsPlaces(groups models.GroupsList) {
	for _, group := range groups.GroupsList {
		fmt.Printf("%v - %v\n", group.Name, group.Place.Name)
	}
	fmt.Println("--------------------------")
}

// Gets the state for x total number of weeks, taking into account the groups
func IterateNextWeeks(weeks int, groups models.GroupsList) models.GroupsList {
	for i := 0; i < weeks; i++ {
		groups.NextIteration()

		// For debugging
		/**
		fmt.Printf("Iteration number %02d\n", i+1)
		ShowGroupsPlaces(groups)
		**/
	}
	return groups
}

// Calculates number of weeks for which to know their respective state in the future
func CalcWeekNum(startDateTime time.Time, endDateTime time.Time) int {

	//t2 := time.Now().UTC()
	//t2 := time.Date(2022, time.Month(1), 31, 1, 0, 0, 0, time.UTC)
	days := endDateTime.Sub(startDateTime).Hours() / 24
	weeks := int(days / 7)

	return weeks
}

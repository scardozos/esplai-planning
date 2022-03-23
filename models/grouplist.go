package models

import (
	pb "github.com/scardozos/esplai-planning/api/grpc/groups"
)

type GroupsList struct {
	GroupsList []*Group
}

func (groups *GroupsList) NextIteration() {
	//fmt.Println(groups.GroupsList[0].Name, groups.GroupsList[1].Place.Name)
	for _, group := range groups.GroupsList {
		//group.NextPlace()
		group.Place = group.Place.Next
	}
	//fmt.Println(groups.GroupsList[0].Name, groups.GroupsList[1].Place.Name)
}

// Returns a slice with a list of groups and their places
// Initial state:
// Teatre | Parc Central |  Pista    | Plaça | Passarel·la
//   <-          <-           <-        <-         <-
// Aneto  |  Pedraforca	 | Matagalls | Cadí  | Puigmal
func InitialGroupState() GroupsList {
	// Initialize places
	passarela := &Place{Name: "Pasarel·la", Next: nil}
	plaza := &Place{Name: "Plaça", Next: passarela}
	pista := &Place{Name: "Pista", Next: plaza}
	parcCentral := &Place{Name: "Parc Central", Next: pista}
	teatre := &Place{Name: "Teatre", Next: parcCentral}
	passarela.Next = teatre

	// Start groups and assign places
	return GroupsList{
		GroupsList: []*Group{
			{Name: "Aneto", Place: teatre},
			{Name: "Pedraforca", Place: parcCentral},
			{Name: "Matagalls", Place: pista},
			{Name: "Cadí", Place: plaza},
			{Name: "Puigmal", Place: passarela},
		},
	}
}

// Translates local group logic declarations to protobuf format
func MarshalGroupModel(groups GroupsList) []*pb.Group {
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
func IterateNextWeeks(weeks int, groups GroupsList) GroupsList {
	for i := 0; i < weeks; i++ {
		groups.NextIteration()
	}
	return groups
}

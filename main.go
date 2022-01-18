package main

import (
	"fmt"

	"github.com/scardozos/esplai-planning/models"
)

func main() {
	// Teatre | Parc Central |  Pista    | Plaça | Passarel·la
	// Aneto  |  Pedraforca	 | Matagalls | Cadí  | Puigmal

	// Initialize places
	passarela := &models.Place{Name: "Pasarel·la", Next: nil}
	plaza := &models.Place{Name: "Plaça", Next: passarela}
	pista := &models.Place{Name: "Pista", Next: plaza}
	parcCentral := &models.Place{Name: "Parc Central", Next: pista}
	teatre := &models.Place{Name: "Teatre", Next: parcCentral}
	passarela.Next = teatre

	// Start groups and assign places
	ant := &models.Group{Name: "Aneto", Place: teatre}
	pdf := &models.Group{Name: "Pedraforca", Place: parcCentral}
	mtg := &models.Group{Name: "Matagalls", Place: pista}
	cdi := &models.Group{Name: "Cadí", Place: plaza}
	pgm := &models.Group{Name: "Puigmal", Place: passarela}

	groups := models.GroupsList{GroupsList: []*models.Group{ant, pdf, mtg, cdi, pgm}}

	/**
	ShowGroupsPlaces(groups) // Shows current config
	groups.NextIteration()   // Switches to next iteration
	ShowGroupsPlaces(groups) // Shows next iteration
	**/

	fmt.Println("Initial state:")
	ShowGroupsPlaces(groups)
	fmt.Println("Iterations:")
	IterateNextWeeks(20, groups)
	// ShowGroupsPlaces(newGroups)

}

func ShowGroupsPlaces(groups models.GroupsList) {
	for _, group := range groups.GroupsList {
		fmt.Printf("%v - %v\n", group.Name, group.Place.Name)
	}
	fmt.Println("--------------------------")
}

func IterateNextWeeks(weeks int, groups models.GroupsList) models.GroupsList {
	for i := 0; i < weeks; i++ {
		groups.NextIteration()

		// For debugging. TODO: Comment out
		fmt.Printf("Iteration number %02d\n", i+1)
		ShowGroupsPlaces(groups)
	}
	return groups
}

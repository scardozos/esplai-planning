package main

import (
	"fmt"
	"time"

	"github.com/scardozos/esplai-planning/models"
)

func main() {
	// Teatre | Parc Central |  Pista    | Plaça | Passarel·la
	// Aneto  |  Pedraforca	 | Matagalls | Cadí  | Puigmal
	const (
		startYear  = 2022
		startMonth = 1
		startDay   = 17
	)
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

	//res := CalcWeekNum([]int{2022, 1, 17})

	startDate := models.DateTime{Year: startYear, Month: startMonth, Day: startDay}
	endDate := models.DateTime{Year: 2022, Month: 1, Day: 24}
	weekNum := CalcWeekNum(startDate.ToTime(), endDate.ToTime())

	newGroups := IterateNextWeeks(weekNum, groups)
	fmt.Println("New state:")
	ShowGroupsPlaces(newGroups)
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

		// For debugging
		/**
		fmt.Printf("Iteration number %02d\n", i+1)
		ShowGroupsPlaces(groups)
		**/
	}
	return groups
}

func CalcWeekNum(startDateTime time.Time, endDateTime time.Time) int {

	//t2 := time.Now().UTC()
	//t2 := time.Date(2022, time.Month(1), 31, 1, 0, 0, 0, time.UTC)
	days := endDateTime.Sub(startDateTime).Hours() / 24
	weeks := int(days / 7)

	return weeks
}

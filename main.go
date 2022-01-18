package main

import (
	"fmt"

	"github.com/scardozos/esplai-planning/models"
)

func main() {
	// Teatre | Parc Central |  Pista    | Plaça | Passarel·la
	// Aneto  |  Pedraforca	 | Matagalls | Cadí  | Puigmal
	passarela := &models.Place{Name: "Pasarel·la", Next: nil}
	plaza := &models.Place{Name: "Plaça", Next: passarela}
	pista := &models.Place{Name: "Pista", Next: plaza}
	parcCentral := &models.Place{Name: "Parc Central", Next: pista}
	teatre := &models.Place{Name: "Teatre", Next: parcCentral}
	passarela.Next = teatre

	ant := &models.Group{Name: "Aneto", Place: teatre}
	pdf := &models.Group{Name: "Pedraforca", Place: parcCentral}
	mtg := &models.Group{Name: "Matagalls", Place: pista}
	cdi := &models.Group{Name: "Cadí", Place: plaza}
	pgm := &models.Group{Name: "Puigmal", Place: passarela}

	groups := models.GroupsList{GroupsList: []*models.Group{ant, pdf, mtg, cdi, pgm}}
	ShowGroupsPlaces(groups)
	groups.NextIteration()
	ShowGroupsPlaces(groups) //next iteration

}

func ShowGroupsPlaces(groups models.GroupsList) {
	for _, group := range groups.GroupsList {
		fmt.Printf("%v - %v\n", group.Name, group.Place.Name)
	}
	fmt.Println("--------------------------")
}

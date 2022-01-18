package models

type Place struct {
	Name string
	Next *Place
}

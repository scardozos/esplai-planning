package models

type Group struct {
	Name string
	*Place
}

/** func (group *Group) NextPlace() {
	group.Place = group.Place.Next
} **/

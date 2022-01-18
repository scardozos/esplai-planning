package models

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

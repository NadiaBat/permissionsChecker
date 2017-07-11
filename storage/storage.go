package storage

import (
	"sync"
)

type Assignment struct {
	userId   int
	itemName string
}

type Assignments []Assignment

var assignments Assignments

func GetAllAssignments(loadIfEmpty bool) Assignments {
	if assignments == nil && loadIfEmpty {
		Refresh()
	}

	return assignments
}

func Refresh() {
	mutex := sync.Mutex{}

	mutex.Lock()
	assignments = getFromDb()
	mutex.Unlock()
}

func getFromDb() Assignments {
	// implement loading from db

	a := Assignment{userId: 123, itemName: "123_name"}
	b := Assignment{userId: 321, itemName: "321_name"}

	return Assignments{a, b}
}

package storage

import (
	"sync"
)

type Assignment struct {
	userId   int
	itemName string
}

type Assignments *[]Assignment

var permissions Assignments

func GetAllPermissions(loadIfEmpty bool) Assignments {
	if permissions == nil && loadIfEmpty {
		Refresh()
	}

	return permissions
}

func Refresh() {
	mutex := sync.Mutex{}

	mutex.Lock()
	permissions = getFromDb()
	mutex.Unlock()
}

func getFromDb() Assignments {
	// implement loading from db

	a := Assignment{userId: 123, itemName: "123_name"}
	b := Assignment{userId: 321, itemName: "321_name"}

	return Assignments{a, b}
}

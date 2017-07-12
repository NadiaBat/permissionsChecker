package storage

import (
	"sync"
)

type Assignment struct {
	userId   int
	itemName string
}
// !!! одна копия!
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
	// get by user

	a := Assignment{userId: 123, itemName: "123_name"}
	b := Assignment{userId: 321, itemName: "321_name"}

	return Assignments{a, b}
}

/*
- user assignments: roles list
- check access recursively

- auth_assignments, auth_item, auth_item_child
All auth data:
1. From auth_item: items => item, ...
2. From auth_item_child: parents[child][] => parent, ...
3. From auth_assignments: assignments[user][item name] => assignment

Грузить все данные не нужно в storage, нужны данныые только для пользователя

 */

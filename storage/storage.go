package storage

import (
	"sync"
)

type Assignment struct {
	UserId   int
	ItemName string
	Rule     string
	Data     string
}

type PermissionItem struct {
	Name     string
	ItemType int
	Rule     string
	Data     string
}

// !!! одна копия!
type Assignments map[string]Assignment

type PermissionItems map[string]PermissionItem

var assignments Assignments
var permissionItems PermissionItems

func GetAllAssignments(loadIfEmpty bool) Assignments {
	if assignments == nil && loadIfEmpty {
		RefreshAssignments()
	}

	return assignments
}

func GetAllPermissionItems(loadIfEmpty bool) PermissionItems {
	if permissionItems == nil && loadIfEmpty {
		RefreshPermissionItems()
	}

	return permissionItems
}

func RefreshAssignments() {
	mutex := sync.Mutex{}

	mutex.Lock()
	assignments = getAssignmentsFromDb()
	mutex.Unlock()
}

func RefreshPermissionItems() {
	mutex := sync.Mutex{}

	mutex.Lock()
	permissionItems = getPermissionItemsFromDb()
	mutex.Unlock()
}

func getAssignmentsFromDb() Assignments {
	// implement loading from db
	// get by user

	a := Assignment{UserId: 123, ItemName: "123_name"}
	b := Assignment{UserId: 321, ItemName: "321_name"}

	return Assignments{"123_name": a, "321_name": b}
}

func getPermissionItemsFromDb() PermissionItems {
	a := make(PermissionItems)
	a["ncc.region.access"] = PermissionItem{
		Name:     "ncc.region.access",
		ItemType: 0,
		Rule:     "",
		Data:     ""}

	return a
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

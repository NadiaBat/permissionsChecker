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

type Assignments map[string]Assignment

type PermissionItems map[string]PermissionItem

var assignments struct {
	data  Assignments
	mutex sync.Mutex
}

var permissionItems struct {
	data  PermissionItems
	mutex sync.Mutex
}

var testMode = false

func SetTestMode(mode bool) {
	testMode = mode
}

func GetAllAssignments(loadIfEmpty bool) Assignments {
	if assignments.data == nil && loadIfEmpty {
		RefreshAssignments()
	}

	return assignments.data
}

func GetAllPermissionItems(loadIfEmpty bool) PermissionItems {
	if permissionItems.data == nil && loadIfEmpty {
		RefreshPermissionItems()
	}

	return permissionItems.data
}

func RefreshAssignments() {
	assignments.mutex.Lock()

	if !testMode {
		assignments.data = getAssignmentsFromDb()
	} else {
		assignments.data = getTestAssignmentsFromDb()
	}

	assignments.mutex.Unlock()
}

func RefreshPermissionItems() {
	permissionItems.mutex.Lock()

	if !testMode {
		permissionItems.data = getPermissionItemsFromDb()
	} else {
		permissionItems.data = getTestPermissionItemsFromDb()
	}

	permissionItems.mutex.Unlock()
}

func getAssignmentsFromDb() Assignments {
	// implement loading from db

	a := Assignment{UserId: 123, ItemName: "123_name"}
	b := Assignment{UserId: 321, ItemName: "321_name"}

	return Assignments{"123_name": a, "321_name": b}
}

func getPermissionItemsFromDb() PermissionItems {
	// implement loading from db
	a := make(PermissionItems)
	a["ncc.region.access"] = PermissionItem{
		Name:     "ncc.region.access",
		ItemType: 0,
		Rule:     "",
		Data:     ""}

	return a
}

func getTestAssignmentsFromDb() Assignments {
	return GetAssignmentsMock()
}

func getTestPermissionItemsFromDb() PermissionItems {
	return GetPermissionItemsMock()
}

// maybe not PermissionItems, only string list
func getTestPermissionParentsFromDb() PermissionItems {
	return GetPermissionParentsMock()
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

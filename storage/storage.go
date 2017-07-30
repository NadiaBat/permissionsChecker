package storage

import (
	"sync"
)

// @see table auth_assignment
type Assignment struct {
	ItemName string
	Rule     string
	Data     string
}

// grouped by item name auth assignments
type UserAssignment struct {
	UserId int
	Items  map[string]Assignment
}

// grouped by user id user auth assignments
type Assignments map[int]UserAssignment

type PermissionItem struct {
	Name     string
	ItemType int
	Rule     string
	Data     string
}

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
	connection := Connection{}
	db := connection.Init()

	res, err := db.Query("SELECT IFNULL(`item_name`, ''), " +
		"IFNULL(`user_id`, 0), " +
		"IFNULL(`biz_rule`, ''), " +
		"IFNULL(`data`, '') " +
		"FROM `auth_assignment`")

	if err != nil {
		panic(err)
	}

	currentUserId := 0
	currentItemName := ""
	currentRule := ""
	currentData := ""

	result := Assignments{}

	for res.Next() {
		err = res.Scan(&currentItemName, &currentUserId, &currentRule, &currentData)
		if err != nil {
			panic(err)
		}
		if currentUserId == 0 || len(currentItemName) == 0 {
			continue
		}

		if result[currentUserId].UserId == 0 {
			result[currentUserId] = UserAssignment{
				UserId: currentUserId,
				Items:  make(map[string]Assignment)}
		}

		result[currentUserId].Items[currentItemName] = Assignment{
			ItemName: currentItemName,
			Rule:     currentRule,
			Data:     currentData}
	}

	return result
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

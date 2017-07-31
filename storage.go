package main

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

type AssignmentRow struct {
	UserId   int
	ItemName string
	Rule     string
	Data     string
}

type Cache struct {
}

var assignments struct {
	data  Assignments
	mutex sync.Mutex
}

var permissionItems struct {
	data  PermissionItems
	mutex sync.Mutex
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

func RefreshAssignments() error {
	var err error
	assignments.mutex.Lock()
	assignments.data, err = getAssignmentsFromDb()
	assignments.mutex.Unlock()
	return err
}

func RefreshPermissionItems() {
	permissionItems.mutex.Lock()
	permissionItems.data = getPermissionItemsFromDb()
	permissionItems.mutex.Unlock()
}

func getAssignmentsFromDb() (Assignments, error) {
	// как узнать пустая ли у тебя выборка?
	rows, err := mysql.Query(
		"SELECT IFNULL(`item_name`, ''), " +
			"IFNULL(`user_id`, 0), " +
			"IFNULL(`biz_rule`, ''), " +
			"IFNULL(`data`, '') " +
			"FROM `auth_assignment`",
	)

	if err != nil {
		panic(err)
	}

	result := Assignments{}
	for rows.Next() {
		var aRow AssignmentRow
		err = rows.Scan(
			&aRow.ItemName,
			&aRow.UserId,
			&aRow.Rule,
			&aRow.Data,
		)
		if err != nil {
			panic(err)
		}
		if aRow.UserId == 0 || len(aRow.ItemName) == 0 {
			continue
		}

		_, exist := result[aRow.UserId]
		if !exist {
			result[aRow.UserId] = UserAssignment{
				UserId: aRow.UserId,
				Items:  make(map[string]Assignment),
			}
		}

		result[aRow.UserId].Items[aRow.ItemName] = Assignment{
			ItemName: aRow.ItemName,
			Rule:     aRow.Rule,
			Data:     aRow.Data,
		}
	}

	return result, nil
}

func getPermissionItemsFromDb() PermissionItems {
	res, err := mysql.Query(
		"SELECT IFNULL(`name`, ''), " +
			"IFNULL(`type`, 0), " +
			"IFNULL(`biz_rule`, ''), " +
			"IFNULL(`data`, '') " +
			"FROM `auth_item`",
	)

	if err != nil {
		panic(err)
	}

	currentName := ""
	currentType := 0
	currentRule := ""
	currentData := ""

	result := PermissionItems{}

	for res.Next() {
		err := res.Scan(&currentName, &currentType, &currentRule, &currentData)
		if err != nil {
			panic(err)
		}
		if len(currentName) == 0 {
			continue
		}

		result[currentName] = PermissionItem{
			Name:     currentName,
			ItemType: currentType,
			Rule:     currentRule,
			Data:     currentData}
	}

	return result
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

!!! auth_users (find usages)

*/

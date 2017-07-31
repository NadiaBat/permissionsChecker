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

var Cache struct {
	assignments struct {
		sync.Mutex
		data Assignments
	}
	permissionItems struct {
		sync.Mutex
		data PermissionItems
	}
}

func GetAllAssignments() Assignments {
	return Cache.assignments.data
}

func GetAllPermissionItems() PermissionItems {
	return Cache.permissionItems.data
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

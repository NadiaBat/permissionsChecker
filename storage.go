package main

import (
	"sync"
)

type Rule struct {
	ParamsKey string
	Data      []string
}

// @see table auth_assignment
type Assignment struct {
	ItemName string
	Rule     Rule
}

// grouped by item ParamsKey auth assignments
type UserAssignment struct {
	UserId int
	Items  map[string]Assignment
}

// grouped by user id user auth assignments
type Assignments map[int]UserAssignment

type PermissionItem struct {
	Name     string
	ItemType int
	Rule     Rule
}

type PermissionItems map[string]PermissionItem

type AssignmentRow struct {
	UserId   int
	ItemName string
	Rule     Rule
}

type ItemParents []string

type AllParents map[string]ItemParents

var Cache struct {
	assignments struct {
		sync.Mutex
		data Assignments
	}
	permissionItems struct {
		sync.Mutex
		data PermissionItems
	}
	parents struct {
		sync.Mutex
		data AllParents
	}
}

func GetAllAssignments() Assignments {
	return Cache.assignments.data
}

func GetAllPermissionItems() PermissionItems {
	return Cache.permissionItems.data
}

func GetAllParents() AllParents {
	return Cache.parents.data
}

/*
- user assignments: roles list
- check access recursively

- auth_assignments, auth_item, auth_item_child
All auth Data:
1. From auth_item: items => item, ...
2. From auth_item_child: parents[child][] => parent, ...
3. From auth_assignments: assignments[user][item ParamsKey] => assignment

Грузить все данные не нужно в storage, нужны данныые только для пользователя

!!! auth_users (find usages)

*/

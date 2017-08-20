package main

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"sync"
)

type Permission struct {
	UserId     int
	ActionName string
	HasAccess  bool
}

type checkingParams struct {
	userId       int
	region       int
	project      int
	isCommercial bool
}

type Permissions []*Permission

type Checker struct {
	sync.WaitGroup
	permissions Permissions
}

func BulkCheck(userId int, actions []string, additionalParams map[string]string) (Permissions, error) {
	checker := &Checker{
		permissions: make(Permissions, len(actions)),
	}

	params, err := getCheckingParams(userId, additionalParams)
	if err != nil {
		return nil, errors.Wrap(err, "Не удалось выполнить проверку.")
	}

	// @todo probably, should not have async checking
	// only for several users (unlikely case)
	for _, action := range actions {
		checker.Add(1)

		permission := &Permission{UserId: userId, ActionName: action}
		go func(permission *Permission) {
			permission.HasAccess = checkAccess(userId, permission.ActionName, params)
			checker.permissions = append(checker.permissions, permission)
			checker.Done()
		}(permission)
	}

	checker.Wait()

	return checker.permissions, nil
}

func getCheckingParams(userId int, additionalParams map[string]string) (*checkingParams, error) {
	params := checkingParams{userId: userId, region: 0, project: 0}
	var err error
	for name, value := range additionalParams {
		switch name {
		case "region":
			params.region, err = strconv.Atoi(value)
		case "project":
			params.project, err = strconv.Atoi(value)
		default:
			continue
		}
	}

	if err != nil {
		return nil, err
	}

	return &params, nil
}

func checkAccess(userId int, actionName string, params *checkingParams) bool {
	userAssignments, err := getUserAssignments(userId)
	if err != nil {
		return false
	}

	return checkAccessRecursive(userId, actionName, params, userAssignments)
}

func getUserAssignments(userId int) (map[string]Assignment, error) {
	allAssignments := GetAllAssignments()
	userAssignments, ok := allAssignments[userId]
	if !ok {
		return nil, errors.New("User assignments doesn`t exists.")
	}

	return userAssignments.Items, nil

}

func checkAccessRecursive(
	userId int, itemName string, params *checkingParams, assignments map[string]Assignment,
) bool {
	permissionItem, err := getPermissionItem(itemName)
	if err != nil {
		return false
	}

	if !executeRule(permissionItem.Rule, params, permissionItem.Data) {
		return false
	}

	itemAssignment, ok := assignments[itemName]
	if ok {
		if executeRule(itemAssignment.Rule, params, itemAssignment.Data) {
			return true
		}
	}

	parents, err := getParents(itemName)
	if err != nil {
		return false
	}

	for _, parentItem := range parents {
		if checkAccessRecursive(userId, parentItem, params, assignments) {
			return true
		}
	}

	return false
}

func getPermissionItem(name string) (*PermissionItem, error) {
	allPermissionItems := GetAllPermissionItems()
	permissionItem, ok := allPermissionItems[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("There is no permission item %s", name))
	}

	return &permissionItem, nil
}

func getParents(childName string) (ItemParents, error) {
	allParents := GetAllParents()
	itemParents, ok := allParents[childName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("There is no parents for item %s", childName))
	}

	return itemParents, nil
}

func executeRule(rule string, params *checkingParams, data string) bool {
	if len(rule) == 0 {
		return true
	}

	// @TODO there is only 1 rule (isCommercial = 1)
	if rule != "News_Permissions_Rules::inArray" || len(data) == 0 {
		return false
	}

	return params.isCommercial
}

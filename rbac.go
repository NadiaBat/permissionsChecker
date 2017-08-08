package main

import (
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
	userId  int
	region  int
	project int
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
		go func(permissions *Permission) {
			permission.HasAccess = checkAccess(permission.ActionName, params)
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

func checkAccess(actionName string, params *checkingParams) bool {
	// implement checking logic
	//allAssignments := GetAllAssignments(true)
	//return checkRecursively(actionName, allAssignments, params)
	return false
}

//func checkRecursively(itemName string, assignments Assignments, params checkingParams) bool {
//	item, err := getItem(itemName)
//	if err != nil {
//		return false
//	}
//	if !executeRule(item.Rule, params, item.Data) {
//	}
//
//	_, isExists := assignments[itemName]
//	if !isExists {
//		assignment := assignments[itemName]
//		if executeRule(assignment.Rule, params, assignment.Data) {
//			return true
//		}
//	}
//
//	parents := getParents(itemName)
//	for _, parent := range parents {
//		if checkRecursively(parent, assignments, params) {
//			return true
//		}
//	}
//
//	return false
//}
//
//func getItem(name string) (PermissionItem, error) {
//	allPermissionItems := GetAllPermissionItems(true)
//	_, isExists := allPermissionItems[name]
//	if isExists {
//		return allPermissionItems[name], nil
//	}
//	return PermissionItem{}, errors.New("Permission item doesn`t exist")
//}
//
//func getParents(itemName string) map[int]string {
//	// get parents
//	// see News_Permissions_Cache_AuthData::getParent
//	result := make(map[int]string)
//	result[0] = "123"
//	result[1] = "222"
//
//	return result
//}
//
//func executeRule(rule string, params checkingParams, data string) bool {
//	// execute rule by params and data
//	// see News_Permission_Checker::executeRule
//	return true
//}

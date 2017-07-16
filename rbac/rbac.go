package rbac

import (
	"github.com/NadiaBat/permissionsChecker/storage"
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

type Permissions []Permission

type Checker struct {
	wg          sync.WaitGroup
	permissions Permissions
}

func BulkCheck(userId int, actions []string, additionalParams map[string]string) Permissions {
	checker := new(Checker)

	checker.permissions = make(Permissions, len(actions))
	checker.wg = sync.WaitGroup{}

	params := getCheckingParams(userId, additionalParams)

	for _, action := range actions {
		checker.wg.Add(1)

		permission := Permission{UserId: userId, ActionName: action}
		go func(permissions *Permission) {
			permission.HasAccess = checkAccess(permission.ActionName, params)
			checker.permissions = append(checker.permissions, permission)
			checker.wg.Done()
		}(&permission)
	}

	checker.wg.Wait()

	return checker.permissions
}

func getCheckingParams(userId int, additionalParams map[string]string) checkingParams {
	params := checkingParams{userId: userId, region: nil, project: nil}
	for name, value := range additionalParams {
		switch name {
		case "region":
			params.region, _ = strconv.Atoi(value)
		case "project":
			params.project, _ = strconv.Atoi(value)
		default:
			continue
		}
	}

	return params
}

func checkAccess(actionName string, params checkingParams) bool {
	// implement checking logic
	allAssignments := storage.GetAllAssignments(true)
	return checkRecursively(actionName, allAssignments, params)
}

func checkRecursively(itemName string, assignments storage.Assignments, params checkingParams) bool {
	item := getItem(itemName)
	if !executeRule(item.Rule, params, item.Data) {
	}

	if assignments[itemName] {
		assignment := assignments[itemName]
		if executeRule(assignment.Rule, params, assignment.Data) {
			return true
		}
	}

	parents := getParents(itemName)
	for _, parent := range parents {
		if checkRecursively(parent, assignments, params) {
			return true
		}
	}

	return false
}

func getItem(name string) storage.PermissionItem {
	allPermissionItems := storage.GetAllPermissionItems(true)
	if allPermissionItems[name] {
		return allPermissionItems[name]
	}

	return nil
}

func getParents(itemName string) map[int]string {
	// get parents
	// see News_Permissions_Cache_AuthData::getParent
	result := make(map[int]string)
	result[0] = "123"
	result[1] = "222"

	return result
}

func executeRule(rule string, params checkingParams, data string) bool {
	// execute rule by params and data
	// see News_Permission_Checker::executeRule
	return true
}

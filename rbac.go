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
	userId       int
	region       int
	project      int
	isCommercial bool
	stringParams map[string]string
}

type Permissions []*Permission

type Checker struct {
	sync.WaitGroup
	permissions Permissions
}

// @TODO 3
func BulkCheck(userId int, actions []string, additionalParams map[string]string) (Permissions, error) {
	checker := &Checker{
		permissions: make(Permissions, len(actions)),
	}

	params, err := getCheckingParams(userId, additionalParams)
	if err != nil {
		return nil, errors.Wrap(err, "Не удалось выполнить проверку.")
	}

	// @TODO 4
	var errs []error
	for _, action := range actions {
		checker.Add(1)

		permission := &Permission{UserId: userId, ActionName: action}
		go func(permission *Permission, errs *[]error) {
			var checkingErr error
			permission.HasAccess, checkingErr = checkAccess(userId, permission.ActionName, params)

			if checkingErr != nil {
				checkingErr = errors.Wrapf(
					checkingErr,
					"Can`t execute checking for userId=%d, actionName=%s",
					permission.UserId,
					permission.ActionName,
				)

				*errs = append(*errs, checkingErr)
			}
			checker.permissions = append(checker.permissions, permission)

			checker.Done()
		}(permission, &errs)
	}

	checker.Wait()
	if len(errs) > 0 {
		return checker.permissions, errs[0]
	}

	return checker.permissions, nil
}

// return checking params from additional params
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

// check access for user and action name
func checkAccess(userId int, actionName string, params *checkingParams) (bool, error) {
	userAssignments, err := getUserAssignments(userId)
	if err != nil {
		return false, errors.Wrapf(
			err,
			"Can`t get user assignments. User id is %d. Action name is \"%s\".",
			userId,
			actionName,
		)
	}

	hasAccess, err := checkAccessRecursive(userId, actionName, params, userAssignments)
	if err != nil {
		return false, errors.Wrapf(
			err,
			"Recursive checking access error. User id is %d. Action name is \"%s\".",
			userId,
			actionName,
		)
	}

	return hasAccess, nil
}

// return all user assignments or error in case of user assignments doesn`t exist
func getUserAssignments(userId int) (map[string]Assignment, error) {
	allAssignments := GetAllAssignments()
	userAssignments, ok := allAssignments[userId]
	if !ok {
		return nil, errors.New("User assignments doesn`t exists.")
	}

	return userAssignments.Items, nil

}

// check access recursive
// get all parents items for checking item while there is no any parents items or access is permitted
func checkAccessRecursive(
	userId int, itemName string, params *checkingParams, assignments map[string]Assignment,
) (bool, error) {
	var err error
	var permissionItem *PermissionItem
	permissionItem = getPermissionItem(itemName)
	if permissionItem == nil {
		return false, err
	}

	var hasAccess bool
	hasAccess, err = executeRule(permissionItem.Rule, params)
	if err != nil {
		return false, errors.Wrapf(err, "Executing rule error. Item name is \"%s\".", itemName)
	}

	if !hasAccess {
		return hasAccess, nil
	}

	itemAssignment, ok := assignments[itemName]
	if ok {
		hasAccess, err = executeRule(itemAssignment.Rule, params)
		if err != nil {
			return false, errors.Wrapf(
				err,
				"Executing rule error for assignment. Item name is \"%s\".",
				itemName,
			)
		}

		if hasAccess {
			return true, nil
		}
	}

	parents := getParents(itemName)
	if parents == nil {
		return false, err
	}

	for _, parentItem := range parents {
		hasAccess, err = checkAccessRecursive(userId, parentItem, params, assignments)
		if err != nil {
			return false, errors.Wrapf(
				err,
				"Parent rule execution error. Item name is \"%s\". Parent name is \"%s\"",
				itemName,
				parentItem,
			)
		}

		if hasAccess {
			return true, nil
		}
	}

	return false, nil
}

// return permission item with rule and data or error in case of item doesn`t exist
func getPermissionItem(name string) *PermissionItem {
	allPermissionItems := GetAllPermissionItems()
	permissionItem, ok := allPermissionItems[name]
	if !ok {
		return nil
	}

	return &permissionItem
}

// return all auth item parents or error in case of parent doesn`t exist
func getParents(childName string) ItemParents {
	allParents := GetAllParents()
	itemParents, ok := allParents[childName]
	if !ok {
		return nil
	}

	return itemParents
}

// execute auth item rule with user or role parameters
// @TODO 6
func executeRule(rule Rule, params *checkingParams) (bool, error) {
	if len(rule.paramsKey) == 0 || len(rule.Data) == 0 {
		return true, nil
	}

	switch rule.paramsKey {
	case "pid":
		hasAccess, err := executeIntegerInArrayRule(rule.Data, params.userId)
		if err != nil {
			return false, errors.Wrap(err, "Param name is \"pid\"")
		}

		return hasAccess, nil
	case "region":
		hasAccess, err := executeIntegerInArrayRule(rule.Data, params.region)
		if err != nil {
			return false, errors.Wrap(err, "Param name is\"region\"")
		}

		return hasAccess, nil
	case "project":
		hasAccess, err := executeIntegerInArrayRule(rule.Data, params.project)
		if err != nil {
			return false, errors.Wrap(err, "Param name is \"project\"")
		}

		return hasAccess, nil
	case "isCommercial":
		hasAccess, err := executeBooleanInArrayRule(rule.Data, params.isCommercial)
		if err != nil {
			return false, errors.Wrap(err, "Param name is \"isCommercial\"")
		}

		return hasAccess, nil
	default:
		// @TODO 5
		if len(rule.Data) == 0 {
			_, ok := params.stringParams[rule.paramsKey]
			return ok, nil
		}

		hasAccess, err := executeStringInArrayRule(rule.Data, params.stringParams[rule.paramsKey])
		if err == nil {
			return false, errors.Wrapf(err, "String param execution error.")
		}

		return hasAccess, nil
	}

	return false, nil
}

// execute in array rule for an integer param
func executeIntegerInArrayRule(data []string, value int) (bool, error) {
	for _, item := range data {
		integerItem, err := strconv.Atoi(item)
		if err != nil {
			return false, errors.Wrapf(err, "Executing rule for integer value error.")
		}

		if integerItem == value {
			return true, nil
		}
	}

	return false, nil
}

// execute in array rule for a boolean param
func executeBooleanInArrayRule(data []string, value bool) (bool, error) {
	for _, item := range data {
		booleanItem := item == "1"

		if booleanItem == value {
			return true, nil
		}
	}

	return false, nil
}

// execute in array rule for a string param
func executeStringInArrayRule(data []string, value string) (bool, error) {
	for _, item := range data {
		if item == value {
			return true, nil
		}
	}

	return false, nil
}

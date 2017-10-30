package main

import (
	"github.com/pkg/errors"
	"strconv"
)

type Permission struct {
	UserId     int
	ActionName string
	HasAccess  bool
}

type AdditionalParams struct{
	userId       int
	region       int
	project      int
	isCommercial bool
	stringParams map[string]string
}

type Permissions []*Permission

type Checker struct {
	permissions Permissions
}

// @TODO 3
func BulkCheck(userId int, actions []string, additionalParams AdditionalParams) (Permissions, error) {
	additionalParams.userId = userId

	// @TODO 4
	checker := &Checker{
		permissions: make(Permissions, len(actions)),
	}

	var err error
	for _, action := range actions {
		permission := &Permission{UserId: userId, ActionName: action}
		permission.HasAccess, err = checkAccess(userId, permission.ActionName, &additionalParams)

		if err != nil {
			return nil, errors.Wrapf(
				err,
				"Can`t execute checking for userId=%d, actionName=%s",
				permission.UserId,
				permission.ActionName,
			)
		}
		checker.permissions = append(checker.permissions, permission)
	}

	return checker.permissions, nil
}

// check access for user and action name
func checkAccess(userId int, actionName string, params *AdditionalParams) (bool, error) {
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
	if userAssignments, ok := allAssignments[userId]; ok {
		return userAssignments.Items, nil
	}

	return nil, errors.New("User assignments doesn`t exists.")
}

// check access recursive
// get all parents items for checking item while there is no any parents items or access is permitted
func checkAccessRecursive(
	userId int, itemName string, params *AdditionalParams, assignments map[string]Assignment,
) (bool, error) {
	var err error
	var permissionItem *PermissionItem
	permissionItem = getPermissionItem(itemName)
	if permissionItem == nil {
		return false, err
	}

	hasAccess := false
	hasAccess, err = executeRule(permissionItem.Rule, params)
	if err != nil {
		return false, errors.Wrapf(err, "Executing rule error. Item name is \"%s\".", itemName)
	}

	if !hasAccess {
		return hasAccess, nil
	}

	if itemAssignment, ok := assignments[itemName]; ok {
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
	if permissionItem, ok := allPermissionItems[name]; ok {
		return &permissionItem
	}

	return nil
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
func executeRule(rule Rule, params *AdditionalParams) (bool, error) {
	if len(rule.ParamsKey) == 0 || len(rule.Data) == 0 {
		return true, nil
	}

	switch rule.ParamsKey {
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
			_, ok := params.stringParams[rule.ParamsKey]
			return ok, nil
		}

		hasAccess, err := executeStringInArrayRule(rule.Data, params.stringParams[rule.ParamsKey])
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

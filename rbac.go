package main

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type Permission struct {
	UserID     int
	ActionName string
	HasAccess  bool
}

type Permissions []*Permission

type PermissionsAll []Permissions

type Checker struct {
	permissions Permissions
}

type Rbac struct {
	dp *DataProvider
}

func NewRbac(dp *DataProvider) *Rbac {
	return &Rbac{
		dp: dp,
	}
}

// @TODO 3
func (r *Rbac) Check(UserID int, action string, additionalParams AdditionalParams) (Permission, error) {
	fmt.Println(UserID)
	additionalParams.UserID = UserID
	fmt.Println(additionalParams.UserID)

	var err error
	permission := Permission{UserID: UserID, ActionName: action}
	permission.HasAccess, err = r.checkAccess(UserID, permission.ActionName, &additionalParams)

	if err != nil {
		return Permission{}, errors.Wrapf(
			err,
			"Can`t execute checking for UserID=%d, actionName=%s",
			permission.UserID,
			permission.ActionName,
		)
	}

	return permission, nil
}

// check access for user and action name
func (r *Rbac) checkAccess(UserID int, actionName string, params *AdditionalParams) (bool, error) {
	userAssignments, err := r.getUserAssignments(UserID)
	if err != nil {
		return false, errors.Wrapf(
			err,
			"Can`t get user assignments. User id is %d. Action name is \"%s\".",
			UserID,
			actionName,
		)
	}

	hasAccess, err := r.checkAccessRecursive(UserID, actionName, params, userAssignments)
	if err != nil {
		return false, errors.Wrapf(
			err,
			"Recursive checking access error. User id is %d. Action name is \"%s\".",
			UserID,
			actionName,
		)
	}

	return hasAccess, nil
}

// return all user assignments or error in case of user assignments doesn`t exist
func (r *Rbac) getUserAssignments(UserID int) (map[string]Assignment, error) {
	allAssignments := r.dp.GetAllAssignments()
	if userAssignments, ok := allAssignments[UserID]; ok {
		return userAssignments.Items, nil
	}

	return nil, errors.New("User assignments doesn`t exists.")
}

// check access recursive
// get all parents items for checking item while there is no any parents items or access is permitted
func (r *Rbac) checkAccessRecursive(
	UserID int, itemName string, params *AdditionalParams, assignments map[string]Assignment,
) (bool, error) {
	var err error
	var permissionItem *PermissionItem
	permissionItem = r.getPermissionItem(itemName)
	if permissionItem == nil {
		return false, err
	}

	hasAccess := false
	hasAccess, err = r.executeRule(permissionItem.Rule, params)
	if err != nil {
		return false, errors.Wrapf(err, "Executing rule error. Item name is \"%s\".", itemName)
	}

	if !hasAccess {
		return hasAccess, nil
	}

	if itemAssignment, ok := assignments[itemName]; ok {
		hasAccess, err = r.executeRule(itemAssignment.Rule, params)
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

	parents := r.getParents(itemName)
	if parents == nil {
		return false, err
	}

	for _, parentItem := range parents {
		hasAccess, err = r.checkAccessRecursive(UserID, parentItem, params, assignments)
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
func (r *Rbac) getPermissionItem(name string) *PermissionItem {
	allPermissionItems := r.dp.GetAllPermissions()
	permissionItem, ok := allPermissionItems[name]
	if !ok {
		return nil
	}
	return &permissionItem
}

// return all auth item parents or error in case of parent doesn`t exist
func (r *Rbac) getParents(childName string) ItemParents {
	allParents := r.dp.GetAllParents()
	itemParents, ok := allParents[childName]
	if !ok {
		return nil
	}

	return itemParents
}

// execute auth item rule with user or role parameters
// @TODO 6
func (r *Rbac) executeRule(rule Rule, params *AdditionalParams) (bool, error) {
	if len(rule.ParamsKey) == 0 || len(rule.Data) == 0 {
		return true, nil
	}

	switch rule.ParamsKey {
	case "pid":
		hasAccess, err := r.executeIntegerInArrayRule(rule.Data, params.UserID)
		if err != nil {
			return false, errors.Wrap(err, "Param name is \"pid\"")
		}

		return hasAccess, nil
	case "region":
		hasAccess, err := r.executeIntegerInArrayRule(rule.Data, params.Region)
		if err != nil {
			return false, errors.Wrap(err, "Param name is\"region\"")
		}

		return hasAccess, nil
	case "project":
		hasAccess, err := r.executeIntegerInArrayRule(rule.Data, params.Project)
		if err != nil {
			return false, errors.Wrap(err, "Param name is \"project\"")
		}

		return hasAccess, nil
	case "isCommercial":
		hasAccess, err := r.executeBooleanInArrayRule(rule.Data, params.IsCommercial)
		if err != nil {
			return false, errors.Wrap(err, "Param name is \"isCommercial\"")
		}

		return hasAccess, nil
	default:
		// @TODO 5
		if len(rule.Data) == 0 {
			_, ok := params.StringParams[rule.ParamsKey]
			return ok, nil
		}

		hasAccess, err := r.executeStringInArrayRule(rule.Data, params.StringParams[rule.ParamsKey])
		if err == nil {
			return false, errors.Wrapf(err, "String param execution error.")
		}

		return hasAccess, nil
	}

	return false, nil
}

// execute in array rule for an integer param
func (r *Rbac) executeIntegerInArrayRule(data []string, value int) (bool, error) {
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
func (r *Rbac) executeBooleanInArrayRule(data []string, value bool) (bool, error) {
	for _, item := range data {
		booleanItem := item == "1"

		if booleanItem == value {
			return true, nil
		}
	}

	return false, nil
}

// execute in array rule for a string param
func (r *Rbac) executeStringInArrayRule(data []string, value string) (bool, error) {
	for _, item := range data {
		if item == value {
			return true, nil
		}
	}

	return false, nil
}

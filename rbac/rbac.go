package main

import (
	"github.com/pkg/errors"
	"fmt"
)

type InputParams map[string]string

type CheckingSet struct {
	UserId int
	AuthItemName string
	Params InputParams
	HasAccess bool
}

type CheckingSets []*CheckingSet

func BulkCheckAccess(checkingSets CheckingSets) (error) {
	var err error
	for _, checkingSet := range checkingSets {
		checkingSet.HasAccess, err = checkAccess(checkingSet.UserId, checkingSet.AuthItemName, checkingSet.Params)
		if err != nil {
			err = errors.Wrapf(
				err,
				"Couldn`t check access. User id: %d, Item name: %s, Input params: %s",
				checkingSet.UserId,
				checkingSet.AuthItemName,
				checkingSet.Params,
			)
		}
	}

	return err
}

// ? inputParams всегда скалярный тип или строка? может ли быть массивом? (входные параметры)
func checkAccess(userId int, authItemName string, inputParams InputParams) (bool, error) {
	userAssignments, ok := allUserAssignments[userId]
	if !ok {
		return false, nil
	}

	rule, ok := authItems[authItemName]

	if !ok {
		return false, errors.New(fmt.Sprintf("Unknown action authItemName %s", authItemName))
	}

	ruleIsOk, err := executeRule(inputParams, rule)
	if err != nil {
		return false, errors.Wrap(err, "Couldn`t execute rule.")
	}

	if !ruleIsOk {
		return false, nil
	}

	userAssignmentItem, ok := userAssignments[authItemName]
	if ok {
		userAuthItem := AuthItem{
			userAssignmentItem.Rule,
			userAssignmentItem.Params,
		}

		userHasAccess, err := executeRule(inputParams, userAuthItem)
		if err != nil {
			return false, errors.Wrap(err, "Couldn`t execute rule.")
		}

		if userHasAccess {
			return true, nil
		}
	}

	parents, ok := authParents[authItemName]
	if !ok {
		return false, nil
	}

	for _, parent := range parents {
		parentHasAccess, err := checkAccess(userId, parent, inputParams)
		if err != nil {
			return false, errors.Wrapf(err, "Couldn`t check parent access. Parent: %s, userId: %d", parent, userId)
		}

		if parentHasAccess {
			return true, nil
		}
	}

	return false, nil
}

func executeRule(inputParams InputParams, rule AuthItem) (bool, error) {
	if rule.Rule == "" {
		return true, nil
	}

	if rule.Rule == "News_Permissions_Rules::inArray" {
		return executeInArray(inputParams, rule.Params), nil
	}

	return false, errors.New(fmt.Sprintf("Unknown rule %s", rule.Rule))
}

func executeInArray(inputParams InputParams, data Data) bool {
	inputValue, ok := inputParams[data.Key]
	if !ok {
		return false
	}

	for _, allowedValue := range data.Data {
		if allowedValue == inputValue {
			return true
		}
	}

	return false
}

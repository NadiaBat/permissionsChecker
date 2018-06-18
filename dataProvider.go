package main

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/NadiaBat/permissionsChecker/phpserialize"
	"github.com/pkg/errors"
)

type Cache struct {
	assignments struct {
		sync.Map
		data Assignments
	}
	permissionItems struct {
		sync.Map
		data PermissionItems
	}
	parents struct {
		sync.Map
		data AllParents
	}
}

type DataProvider struct {
	cache Cache
	db    *sql.DB
}

func NewDataProvider(db *sql.DB) *DataProvider {
	return &DataProvider{
		cache: Cache{},
		db:    db,
	}
}

func (dp *DataProvider) GetAllAssignments() Assignments {
	if data, ok := dp.cache.assignments.Load("data"); ok {
		fmt.Println("LOAD ASSIGNMENTS")
		return data.(Assignments)
	}
	if data, err := dp.getAssignmentsFromDb(); err == nil {
		fmt.Println("STORE ASSIGNMENTS")
		dp.cache.assignments.LoadOrStore("data", data)
		return data
	}
	return Assignments{}
}

func (dp *DataProvider) GetAllPermissions() PermissionItems {
	if data, ok := dp.cache.permissionItems.Load("data"); ok {
		return data.(PermissionItems)
	}
	if data, err := dp.getPermissionItemsFromDb(); err == nil {
		dp.cache.permissionItems.LoadOrStore("data", data)
		return data
	}
	return PermissionItems{}
}

func (dp *DataProvider) GetAllParents() AllParents {
	if data, ok := dp.cache.parents.Load("data"); ok {
		return data.(AllParents)
	}
	if data, err := dp.getParentsFromDb(); err == nil {
		dp.cache.parents.LoadOrStore("data", data)
		return data
	}
	return AllParents{}
}

func (dp *DataProvider) getAssignmentsFromDb() (Assignments, error) {
	rows, err := dp.db.Query(
		"SELECT IFNULL(`item_name`, ''), IFNULL(`user_id`, 0), IFNULL(`Data`, '') FROM `auth_assignment`",
	)

	result := Assignments{}
	if err != nil {
		if err == sql.ErrNoRows {
			return result, nil
		}
		return nil, errors.Wrap(err, "Get all users assignments failed.")
	}

	var currentRule string
	rule := Rule{}

	for rows.Next() {
		aRow := AssignmentRow{}
		err = rows.Scan(
			&aRow.ItemName,
			&aRow.UserID,
			&currentRule,
		)

		if err != nil {
			errors.Wrapf(err, "Assignment row scanning error.")
		}

		if aRow.UserID == 0 || len(aRow.ItemName) == 0 {
			errors.Wrap(err, "Empty UserID or itemName for assignment row.")
			continue
		}

		_, exist := result[aRow.UserID]
		if !exist {
			result[aRow.UserID] = UserAssignment{
				UserID: aRow.UserID,
				Items:  make(map[string]Assignment),
			}
		}

		rule, err = dp.getRuleFromSerialized(currentRule)
		if err != nil {
			return nil, errors.Wrapf(err, "Unserialize currentRule failed. Rule was \"%s\"", currentRule)
		}

		result[aRow.UserID].Items[aRow.ItemName] = Assignment{
			ItemName: aRow.ItemName,
			Rule:     rule,
		}
	}

	return result, err
}

func (dp *DataProvider) getPermissionItemsFromDb() (PermissionItems, error) {
	res, err := dp.db.Query("SELECT IFNULL(`name`, ''),IFNULL(`type`, 0), IFNULL(`Data`, '') FROM `auth_item`")

	if err != nil {
		return nil, errors.Wrapf(err, "Auth items getting failed.")
	}

	currentName := ""
	currentType := 0
	currentRule := ""

	rule := Rule{}

	items := PermissionItems{}

	for res.Next() {
		var currentErr error
		currentErr = res.Scan(&currentName, &currentType, &currentRule)
		if currentErr != nil {
			err = errors.Wrap(err, "Auth item row scanning error.")
			continue
		}

		if len(currentName) == 0 {
			err = errors.Wrap(err, "Auth item ParamsKey is empty.")
			continue
		}

		rule, currentErr = dp.getRuleFromSerialized(currentRule)
		if currentErr != nil {
			err = errors.Wrapf(err, "Rule json decode error. Rule was \"%s\"", currentRule)
			continue
		}

		items[currentName] = PermissionItem{
			Name:     currentName,
			ItemType: currentType,
			Rule:     rule,
		}
	}

	return items, err
}

// @TODO 1
func (dp *DataProvider) getParentsFromDb() (AllParents, error) {
	res, err := dp.db.Query("SELECT `child`, `parent` FROM `auth_item_child`")

	if err != nil {
		return nil, errors.Wrapf(err, "Parents getting failed.")
	}

	currentChild := ""
	currentParent := ""
	parents := AllParents{}
	for res.Next() {
		err := res.Scan(&currentChild, &currentParent)
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"Auth item row scanning error with child %s and parent %s.",
				currentChild,
				currentParent,
			)
		}

		parents[currentChild] = append(parents[currentChild], currentParent)
	}

	return parents, nil
}

func (dp *DataProvider) getRuleFromSerialized(rule string) (Rule, error) {
	decoded, err := phpserialize.Decode(rule)
	if err != nil {
		return Rule{}, errors.Wrapf(err, "Rule decoding error. Rule was \"%s\"", rule)
	}

	result, ok := decoded.(map[interface{}]interface{})
	if !ok {
		return Rule{}, errors.Wrapf(err, "Rule decoding error. Rule was \"%s\"", rule)
	}

	key, ok := result["paramsKey"].(string)
	if !ok {
		return Rule{}, err
	}

	resultRule := Rule{}
	resultRule.ParamsKey = key

	data, ok := result["data"]
	if !ok || data == nil {
		return resultRule, err
	}

	decodedData := data.(map[interface{}]interface{})
	for _, value := range decodedData {
		stringValue := fmt.Sprintf("%v", value)
		resultRule.Data = append(resultRule.Data, stringValue)
	}

	return resultRule, nil
}

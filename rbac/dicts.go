package main

import (
	"database/sql"
	"github.com/pkg/errors"
	"encoding/json"
)

type Data struct {
	Key  string `json:"paramsKey"`
	Data []interface{} `json:"data"`
}

type AuthItem struct {
	Rule   string
	Params Data
}

type AuthItems map[string]AuthItem

var authItems AuthItems

type AuthParent []string

type AuthParents map[string]AuthParent

var authParents AuthParents

type Assignment struct {
	UserId int
	ItemName string
	Rule string
	Params Data
}

type UserAssignment map[string]Assignment

type UserAssignments map[int]UserAssignment

var allUserAssignments UserAssignments

func InitDicts() error {
	err := InitAuthItemsFromDb()
	if err != nil {
		return errors.Wrap(err, "Couldn`t init auth items.")
	}

	err = InitParentsFromDb()
	if err != nil {
		return errors.Wrap(err, "Couldn`t init auth parents items.")
	}

	err = InitUserAssignmentsFromDb()
	if (err != nil) {
		return errors.Wrap(err, "Couldn`t init users assignments.")
	}

	return nil
}

// auth_item (project, name, type, description, biz_rule, data)
func InitAuthItemsFromDb() (error) {
	authItems = make(AuthItems)
	rows, err := mysql.Query(
		"SELECT IFNULL(`name`, ''), " +
			"IFNULL(`biz_rule`, ''), " +
			"IFNULL(`data`, '') " +
			"FROM `auth_item`",
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		return errors.Wrap(err, "Get auth items failed.")
	}


	for rows.Next() {
		name := ""
		data := ""
		item := AuthItem{}
		err = rows.Scan(
			&name,
			&item.Rule,
			&data,
		)

		if err != nil {
			errors.Wrapf(
				err,
				"Auth item scanning error. name: %s, item: %s, data: %s",
				name,
				item.Rule,
				data,
			)
			continue
		}

		if len(data) > 0 {
			json.Unmarshal([]byte(data), &item.Params)
		}

		if len(name) > 0 {
			authItems[name] = item
		}
	}

	return err
}

// auth_item_child (child, parent)
func InitParentsFromDb() error {
	authParents = make(AuthParents)
	rows, err := mysql.Query(
		"SELECT IFNULL(`child`, ''), " +
			"IFNULL(`parent`, '') " +
			"FROM `auth_item_child`",
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		return errors.Wrap(err, "Get auth parent items failed.")
	}


	for rows.Next() {
		child := ""
		parent := ""
		err = rows.Scan(
			&child,
			&parent,
		)

		if err != nil {
			errors.Wrapf(
				err,
				"Auth parent item scanning error. child: %s, parent: %s",
				child,
				parent,
			)
			continue
		}

		authParents[child] = append(authParents[child], parent)
	}

	return err
}

// auth_assignment (user_id, item_name, biz_rule, data)
func InitUserAssignmentsFromDb() error {
	allUserAssignments = make(UserAssignments)
	rows, err := mysql.Query(
		"SELECT IFNULL(`user_id`, 0), " +
			"IFNULL(`item_name`, ''), " +
			"IFNULL(`biz_rule`, ''), " +
			"IFNULL(`data`, '') " +
			"FROM `auth_assignment`",
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		return errors.Wrap(err, "Get user assignments failed.")
	}


	for rows.Next() {
		item := Assignment{}
		data := ""
		err = rows.Scan(
			&item.UserId,
			&item.ItemName,
			&item.Rule,
			&data,
		)

		if err != nil {
			errors.Wrapf(
				err,
				"Auth item scanning error. userId: %d, name: %s, item: %s, data: %s",
				item.UserId,
				item.ItemName,
				item.Rule,
				data,
			)
			continue
		}

		if len(data) > 0 {
			json.Unmarshal([]byte(data), &item.Params)
		}

		if allUserAssignments[item.UserId] == nil {
			allUserAssignments[item.UserId] = make(UserAssignment)
		}
		// стоит ли строить полный список по дереву сразу же (найти все итемы для пользователя, включая дочерние)?
		allUserAssignments[item.UserId][item.ItemName] = item
	}

	return err
}

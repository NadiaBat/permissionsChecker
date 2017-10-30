package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/NadiaBat/permissionsChecker/phpserialize"
)

type MySQLConnectionConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Base       string `yaml:"base"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Parameters struct {
		MaxIdleConns int `yaml:"max_idle_conns"`
		MaxOpenConns int `yaml:"max_open_conns"`
	} `yaml:"parameters"`
}

func NewMySQL(config *MySQLConnectionConfig) (*sql.DB, error) {
	db, err := openDBConnection(config)
	if err != nil {
		return nil, errors.Wrap(err, "Fail create mysql client")
	}
	return db, nil
}

func openDBConnection(config *MySQLConnectionConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.GetDSN())
	if err != nil {
		return nil, errors.Wrapf(err, "can't open mysql connection"+
			" for dsn \"%s\"", config.GetDSN())
	}

	db.SetMaxIdleConns(config.Parameters.MaxIdleConns)
	db.SetMaxOpenConns(config.Parameters.MaxOpenConns)

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "can't ping mysql "+
			"after open connection")
	}
	return db, nil
}

func (config MySQLConnectionConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Base,
	)
}

func RefreshCache() error {
	err := RefreshAssignments()
	if err != nil {
		return errors.Wrap(err, "Assignments refreshing failed")
	}

	err = RefreshPermissionItems()
	if err != nil {
		return errors.Wrap(err, "Permission items refreshing failed")
	}

	err = RefreshParents()
	if err != nil {
		return errors.Wrap(err, "Parents refreshing failed")
	}

	return nil
}

func RefreshAssignments() error {
	var err error
	Cache.assignments.Lock()
	Cache.assignments.data, err = getAssignmentsFromDb()
	Cache.assignments.Unlock()
	return err
}

func RefreshPermissionItems() error {
	var err error
	Cache.permissionItems.Lock()
	Cache.permissionItems.data, err = getPermissionItemsFromDb()
	Cache.permissionItems.Unlock()
	return err
}

func RefreshParents() error {
	var err error
	Cache.parents.Lock()
	Cache.parents.data, err = getParentsFromDb()
	Cache.parents.Unlock()
	return err
}

func getAssignmentsFromDb() (Assignments, error) {
	// как узнать пустая ли у тебя выборка?
	rows, err := mysql.Query(
		"SELECT IFNULL(`item_name`, ''), " +
			"IFNULL(`user_id`, 0), " +
			"IFNULL(`Data`, '') " +
			"FROM `auth_assignment`",
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
			&aRow.UserId,
			&currentRule,
		)

		if err != nil {
			errors.Wrapf(err, "Assignment row scanning error.")
		}

		if aRow.UserId == 0 || len(aRow.ItemName) == 0 {
			errors.Wrap(err, "Empty userId or itemName for assignment row.")
			continue
		}

		_, exist := result[aRow.UserId]
		if !exist {
			result[aRow.UserId] = UserAssignment{
				UserId: aRow.UserId,
				Items:  make(map[string]Assignment),
			}
		}

		rule, err = getRuleFromSerialized(currentRule)
		if err != nil {
			return nil, errors.Wrapf(err, "Unserialize currentRule failed. Rule was \"%s\"", currentRule)
		}

		result[aRow.UserId].Items[aRow.ItemName] = Assignment{
			ItemName: aRow.ItemName,
			Rule:     rule,
		}
	}

	return result, err
}

// @TODO 1
func getPermissionItemsFromDb() (PermissionItems, error) {
	res, err := mysql.Query(
		"SELECT IFNULL(`name`, ''), " +
			"IFNULL(`type`, 0), " +
			"IFNULL(`Data`, '') " +
			"FROM `auth_item`",
	)

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

		rule, currentErr = getRuleFromSerialized(currentRule)
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
func getParentsFromDb() (AllParents, error) {
	res, err := mysql.Query("SELECT `child`, `parent` FROM `auth_item_child`")

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

func getRuleFromSerialized(rule string) (Rule, error) {
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

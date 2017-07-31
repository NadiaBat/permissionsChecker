package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type MySQLConnectionConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Base       string `yaml:"base"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Parameters struct {
		MaxIdleConns int
		MaxOpenConns int
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
	//return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
	//	config.User,
	//	config.Password,
	//	config.Host,
	//	config.Port,
	//	config.Base,
	//)
	return "ngs_regionnews:nae9be9eiW@tcp(192.168.134.144:3306)/ngs_regionnews"
}

func RefreshCache()  {
	RefreshAssignments()
	RefreshPermissionItems()
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

func getAssignmentsFromDb() (Assignments, error) {
	// как узнать пустая ли у тебя выборка?
	rows, err := mysql.Query(
		"SELECT IFNULL(`item_name`, ''), " +
			"IFNULL(`user_id`, 0), " +
			"IFNULL(`biz_rule`, ''), " +
			"IFNULL(`data`, '') " +
			"FROM `auth_assignment`",
	)

	if err != nil {
		return nil, errors.Wrap(err, "Get all users assignments failed.")
	}

	result := Assignments{}
	for rows.Next() {
		var aRow AssignmentRow
		err = rows.Scan(
			&aRow.ItemName,
			&aRow.UserId,
			&aRow.Rule,
			&aRow.Data,
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

		result[aRow.UserId].Items[aRow.ItemName] = Assignment{
			ItemName: aRow.ItemName,
			Rule:     aRow.Rule,
			Data:     aRow.Data,
		}
	}

	return result, err
}

func getPermissionItemsFromDb() (PermissionItems, error) {
	res, err := mysql.Query(
		"SELECT IFNULL(`name`, ''), " +
			"IFNULL(`type`, 0), " +
			"IFNULL(`biz_rule`, ''), " +
			"IFNULL(`data`, '') " +
			"FROM `auth_item`",
	)

	if err != nil {
		return nil, errors.Wrapf(err, "Auth items getting failed.")
	}

	currentName := ""
	currentType := 0
	currentRule := ""
	currentData := ""

	items := PermissionItems{}

	for res.Next() {
		err := res.Scan(&currentName, &currentType, &currentRule, &currentData)
		if err != nil {
			errors.Wrapf(err, "Auth item row scanning error.")
			continue
		}

		if len(currentName) == 0 {
			errors.Wrapf(err, "Auth item name is empty.")
			continue
		}

		items[currentName] = PermissionItem{
			Name:     currentName,
			ItemType: currentType,
			Rule:     currentRule,
			Data:     currentData,
		}
	}

	return items, err
}

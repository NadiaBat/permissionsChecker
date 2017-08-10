package main

import (
	"database/sql"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type message struct {
	isActive chan int
}

type user struct {
	id  int
	fio string
}

var (
	mysql *sql.DB
)

func main() {
	configFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Can`t read config.yml."))
	}

	config := MySQLConnectionConfig{}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Yaml config decoding error."))
	}

	mysql, err = NewMySQL(&config)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Mysql error."))
	}

	RefreshCache()

	if err != nil {
		log.Fatal(errors.Wrap(err, "Can`t find user assignments"))
	}

	additionalParams := make(map[string]string)
	additionalParams["region"] = "54"
	additionalParams["project"] = "1"

	actions := []string{"ncc.unblock.record.access"}
	res, err := BulkCheck(223814181, actions, additionalParams)
	println(res)

	//assignments := storage.GetAllAssignments(true)
	//userAssignments := assignments[200132743]
	//println(userAssignments.UserId)
	//
	//for _, assignment := range userAssignments.Items {
	//	println(assignment.ItemName)
	//}

	//permissionItems := GetAllPermissionItems()
	//for key, item := range permissionItems {
	//	println(key, item.Name, item.ItemType)
	//}

	//
	//h := http.Server{}
	//h.Serve()
}

// bulk checker by user and operations

// cli
// --host
// --port
// --dbDsn
// docopt.go (for cli generations) @see git.rn/projects/PORTAL/repos/paged/browse

// for http use standart golang http (/net/http) @see https://cryptic.io/go-http/

// use channels for updating data from db

// handlers for bulk operations

// use Notify for hardware interruptions
// use logger instead of panic, panic as an exclusion

// teeeestsss!!!!!!!!!!!!!!! don`t run? only tests

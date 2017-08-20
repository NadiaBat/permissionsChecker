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

	server := Server{}
	server.Serve()
}

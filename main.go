package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type message struct {
	isActive chan int
}

type user struct {
	id  int
	fio string
}

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

	fmt.Println("start")
	mysql, err := NewMySQL(&config)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Mysql error."))
	}

	rbacDataProvider := NewDataProvider(mysql)
	rbac := NewRbac(rbacDataProvider)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Cache refreshing failed"))
	}

	server := Server{}
	server.Serve(rbac)
}

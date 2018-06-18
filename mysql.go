package main

import (
	"database/sql"
	"fmt"

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

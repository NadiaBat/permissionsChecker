package storage

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type Connection struct {
}

// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
func (connection *Connection) Init() *sql.DB {
	db, err := sql.Open("mysql", "ngs_regionnews:nae9be9eiW@tcp(192.168.134.144:3306)/ngs_regionnews")
	if err != nil {
		panic(err)
	}

	return db
}

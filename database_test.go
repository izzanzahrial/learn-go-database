package learn_go_database

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql" // use only the innit method by using "_"
)

// https://pkg.go.dev/database/sql
// database driver : https://github.com/golang/go/wiki/SQLDrivers

func TestOpenConnection(t *testing.T) {
	db, err := sql.Open("mysql", "username:password@tcp(host:3306)/dbname")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

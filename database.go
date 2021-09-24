package learn_go_database

import (
	"database/sql"
	"time"
)

func GetConnection() *sql.DB {
	// "parseTime=true" parse date type from []uint8 to time.Time
	db, err := sql.Open("mysql", "username:password@tcp(host:3306)/dbname?parseTime=true")
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)                  // Maximum number of idle connection
	db.SetMaxOpenConns(100)                 // Maximum number of connection
	db.SetConnMaxIdleTime(10 * time.Minute) // Maximum time of connection being idle
	db.SetConnMaxLifetime(60 * time.Minute) // Maximum time of connection

	return db
}

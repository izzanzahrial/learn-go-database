package learn_go_database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

// ExecContext
func TestExecSql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	query := "INSERT INTO customer(id, name) VALUES('izzan', 'Izzan')"
	_, err := db.ExecContext(ctx, query) // ExecContext = query sql without return value
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new customer")
}

// QueryContext
func TestQuerySql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	query := "SELECT id, name FROM customer"
	rows, err := db.QueryContext(ctx, query) // QueryContext = query sql with return value
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}

		fmt.Println("id:", id)
		fmt.Println("name:", name)
	}
}

// Nullable value
func TestQuerySqlComplex(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	query := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx, query) // QueryContext = query sql with return value
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		var email sql.NullString
		var balance int32
		var rating float64
		var birthDate sql.NullTime
		var createdAt time.Time
		var married bool
		// if your table have nullable, rather than using standart golang data type
		// use nullable data type from package database/sql
		// because golang doesnt support nullable value other than empty interface
		// string = NullString
		// bool = NullBool
		// float64 = NullFloat64
		// int32 = NullInt32
		// int64 = NullInt64
		// time.Time = NullTime
		// https://pkg.go.dev/database/sql#NullBool
		err := rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}

		fmt.Println("id:", id, "name:", name, "balance:", balance, "rating:",
			rating, "married:", married, "created at:", createdAt)
		if email.Valid {
			fmt.Println("email:", email.String) // nullable
		}
		if birthDate.Valid {
			fmt.Println("birth date:", birthDate.Time) // nullable
		}
	}
}

// Sql Injection
func TestSqlInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	// sql injection
	username := "admin'; #" // the ";" will end the query, the "#" will make the rest of the query turn into comments
	password := "izzan"

	query := "SELECT username FROM user WHERE username = '" + username + "' AND passowrd = '" + password + "' LIMIT 1"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Login success")
	} else {
		fmt.Println("Login failed")
	}
}

// Handle Sql Injection
func TestQueryParameter(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	// sql injection
	username := "admin'; #"
	password := "izzan"

	query := "SELECT username FROM user WHERE username = ? AND passowrd = ? LIMIT 1" // "?" is the value for QueryContext variadic function
	rows, err := db.QueryContext(ctx, query, username, password)                     // By using variadic function in QueryContext
	// we can handle the sql injection
	// you can also do this in ExecContext
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Login success")
	} else {
		fmt.Println("Login failed")
	}
}

// Autoincrement and Last ID
func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	// sql injection
	email := "izzanzahrial@gmail.com"
	comment := "Test"

	query := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	result, err := db.ExecContext(ctx, query, email, comment)
	if err != nil {
		panic(err)
	}
	insertId, err := result.LastInsertId() // "LastInsertID" to get the last inserted ID
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert the new comments ID :", insertId)
}

// Prepare Statement
// let say you want to input multiple data with the same query but different parameter
// that's where prepare statement come in to play
// by using prepare statement you dont have to make a new connection to the DB pool, only using the same connection
// actually in ExecContext and QueryContext also using prepare statement, but
// they make a new connection for each individual query
// https://pkg.go.dev/database/sql#Stmt
func TestPrepareStatement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	query := "INSERT INTO comments(email, comments) VALUES(?, ?)" // using the same query
	statement, err := db.PrepareContext(ctx, query)               // Prepare statement
	if err != nil {
		panic(err)
	}
	defer statement.Close()

	for i := 0; i < 10; i++ { // Create multiple parameter using for loop
		email := "izzan" + strconv.Itoa(i) + "gmail.com"
		comment := strconv.Itoa(i) + "comment"

		result, err := statement.ExecContext(ctx, email, comment) // only pass the parameter, no query needed
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Println("ID:", id)
	}
}

// Database Transaction
// https://pkg.go.dev/database/sql#Tx
func TestTransaction(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	tx, err := db.Begin() // Begin the transaction
	if err != nil {
		panic(err)
	}

	query := "INSERT INTO comments(email, comment) VALUES(?, ?)"

	for i := 0; i < 10; i++ {
		email := "izzan" + strconv.Itoa(i) + "gmail.com"
		comment := strconv.Itoa(i) + "comment"

		result, err := tx.ExecContext(ctx, query, email, comment) // Transaction execcontext
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Println("ID:", id)
	}

	err = tx.Commit() // Commit the transaction
	if err != nil {
		panic(err)
	}
}

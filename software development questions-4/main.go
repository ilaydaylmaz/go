package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open a connection to an SQLite database
	db, err := sql.Open("sqlite3", "./user.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create a new table
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		panic(err)
	}

	// Insert a new user into the table
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	stmt, err := tx.Prepare("INSERT INTO user(name, email) VALUES(?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec("John Doe", "john.doe@example.com")
	if err != nil {
		panic(err)
	}
	tx.Commit()

	// Query all users from the table
	rows, err := db.Query("SELECT id, name, email FROM user")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name, email string
		err = rows.Scan(&id, &name, &email)
		if err != nil {
			panic(err)
		}
		fmt.Println(id, name, email)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
}

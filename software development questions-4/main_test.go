package main

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setup() {
	// Open a connection to an SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create a new table for testing
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

	// Set the global variable 'db' to the test database
	db = testDB
}

func tearDown() {
	// Clean up resources if needed
}

func TestInsertUser(t *testing.T) {
	// Insert a new user into the test table
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO user(name, email) VALUES(?, ?)")
	if err != nil {
		t.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec("John Doe", "john.doe@example.com")
	if err != nil {
		t.Fatal(err)
	}
	tx.Commit()

	// Query the test user from the table
	row := db.QueryRow("SELECT id, name, email FROM user WHERE name=?", "John Doe")

	var id int
	var name, email string
	err = row.Scan(&id, &name, &email)
	if err != nil {
		t.Fatal(err)
	}

	// Assert that the inserted user matches the expected values
	assert.Equal(t, "John Doe", name)
	assert.Equal(t, "john.doe@example.com", email)
}

func TestQueryAllUsers(t *testing.T) {
	// Query all users from the test table
	rows, err := db.Query("SELECT id, name, email FROM user")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	// Assert that there are no users in the table initially
	assert.False(t, rows.Next(), "No users should be initially present in the table")

	// Insert a new user into the test table
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO user(name, email) VALUES(?, ?)")
	if err != nil {
		t.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec("Jane Doe", "jane.doe@example.com")
	if err != nil {
		t.Fatal(err)
	}
	tx.Commit()

	// Query all users from the test table again
	rows, err = db.Query("SELECT id, name, email FROM user")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	// Assert that the inserted user is now present in the table
	assert.True(t, rows.Next(), "There should be at least one user in the table")

}

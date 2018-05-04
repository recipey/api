package main_test

import (
	"log"
	"os"
	"testing"

	"github.com/recipey/api"
)

var a main.App

func TestMain(m *testing.M) {
	a = main.App{}

	a.Initialize(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"),
	)

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM users")
	a.DB.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
}

const tableCreationQuery = `
  CREATE TABLE IF NOT EXISTS users(
		id SERIAL,
		name TEXT NOT NULL,
		CONSTRAINT users_pkey PRIMARY KEY (id)
	)
`

package main

import (
	"database/sql"
)

type user struct {
	ID                   int    `json:"id"`
	Username             string `json:"username"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

// (u *user) means this function can only be called on a
// pointer to a user type
func (u *user) getUser(db *sql.DB) error {
	return db.QueryRow("SELECT id, username, email FROM users WHERE id=$1",
		u.ID).Scan(&u.ID, &u.Username, &u.Email)
}

func (u *user) updateUser(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE users SET username=$1 WHERE id=$2", u.Username, u.ID)

	return err
}

func (u *user) deleteUser(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", u.ID)

	return err
}

func (u *user) createUser(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO users(username, email) VALUES($1, $2) RETURNING id",
		u.Username, u.Email).Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func getUsers(db *sql.DB, start, count int) ([]user, error) {
	rows, err := db.Query("SELECT id, username, email FROM users LIMIT $1 OFFSET $2", count, start)

	if err != nil {
		return nil, err
	}

	// deferred function isn't called until the surrounding function returns
	// if you have multiple defer functions they're executed in LIFO order
	// any arguments you have for the function are evaluated immediately
	defer rows.Close()

	users := []user{}

	for rows.Next() {
		var u user
		if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

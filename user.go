package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type user struct {
	ID                   int    `json:"id"`
	Username             string `json:"username"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

// NOTE: to customize errors and json marshalling for any `entity`
// * create a custom error struct
// * define a Error() method to customize error message
// * define a MarshalJSON() method to customize JSON marshalling
type userError struct {
	Errors map[string]string
}

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
	ue := &userError{Errors: map[string]string{}}
	if u.Username == "" {
		ue.Errors["username"] = "username is required"
	}

	if u.Email == "" {
		ue.Errors["email"] = "email is required"
	}

	if u.Password == "" {
		ue.Errors["password"] = "password is required"
	}

	if u.Password != u.PasswordConfirmation {
		ue.Errors["password_confirmation"] = "password confirmation does not match password"
	}

	if len(ue.Errors) > 0 {
		return ue
	}

	err := db.QueryRow(
		"INSERT INTO users(username, email, password) VALUES($1, $2, $3) RETURNING id",
		u.Username, u.Email, u.Password).Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func getUsers(db *sql.DB, start, count int) ([]user, error) {
	rows, err := db.Query(
		"SELECT id, username, email FROM users LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

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

func (ue *userError) Error() string {
	errStr, _ := json.Marshal(ue.Errors)

	return string(errStr)
}

func (ue *userError) MarshalJSON() ([]byte, error) {
	aux := ue.Errors

	jsonStr, err := json.Marshal(aux)

	if err != nil {
		return []byte{}, fmt.Errorf("decode userError: %v", err)
	}

	return jsonStr, nil
}

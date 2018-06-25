package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type recipe struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Author    string `json:"author"`
	ImageUrl  string `json:"image_url"`
	SourceUrl string `json:"source_url"`
}

type recipeError struct {
	Errors map[string]string
}

func getRecipes(db *sql.DB, start, count int) ([]recipe, error) {
	rows, err := db.Query("SELECT id, name, author, image_url, source_url FROM recipes LIMIT $1 OFFSET $2", count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	recipes := []recipe{}

	for rows.Next() {
		var r recipe
		if err := rows.Scan(&r.ID, &r.Name, &r.Author, &r.ImageUrl, &r.SourceUrl); err != nil {
			return nil, err
		}

		recipes = append(recipes, r)
	}

	return recipes, nil
}

func (r *recipe) getRecipe(db *sql.DB) error {
	return db.QueryRow("SELECT id, name, author, image_url, source_url FROM recipes WHERE id=$1", r.ID).Scan(&r.ID, &r.Name, &r.Author, &r.ImageUrl, &r.SourceUrl)
}

func (r *recipe) createRecipe(db *sql.DB) error {
	re := &recipeError{Errors: map[string]string{}}
	if r.Name == "" {
		re.Errors["name"] = "name is required"
	}

	if len(re.Errors) > 0 {
		return re
	}

	err := db.QueryRow(
		"INSERT INTO recipes(name, author, image_url, source_url) VALUES($1, $2, $3, $4) RETURNING id", r.Name, r.Author, r.ImageUrl, r.SourceUrl).Scan(&r.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *recipe) updateRecipe(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE recipes SET name=$1, author=$2, image_url=$3, source_url=$4 WHERE id=$5",
			r.Name, r.Author, r.ImageUrl, r.SourceUrl, r.ID)

	return err
}

func (r *recipe) deleteRecipe(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM recipes WHERE id=$1", r.ID)

	return err
}

func (re *recipeError) Error() string {
	errStr, _ := json.Marshal(re.Errors)

	return string(errStr)
}

func (re *recipeError) MarshalJSON() ([]byte, error) {
	aux := re.Errors

	jsonStr, err := json.Marshal(aux)

	if err != nil {
		return []byte{}, fmt.Errorf("decode recipeError: %v", err)
	}

	return jsonStr, nil
}

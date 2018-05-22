package main

import (
	"database/sql"
)

type recipe struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Author    string `json:"author"`
	ImageUrl  string `json:"image_url"`
	SourceUrl string `json:"source_url"`
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

func searchRecipes(db *sql.DB, term string) ([]recipe, error) {
	//postgresql query need typo correction (fuzzy matching?) and full text search
	//
	// if err != nil {
	// 	return nil, err
	// }
	//
	// defer rows.Close()

	recipes := []recipe{}

	// Test data
	recipes = append(recipes, recipe{
		ID: 1,
		Name: "Fried Rice",
		Author: "Betty",
		ImageUrl: "https://www.spendwithpennies.com/wp-content/uploads/2016/02/fried-rice-recipe-21.jpg",
		SourceUrl: "https://www.spendwithpennies.com/easy-fried-rice-recipe/",
	})

	recipes = append(recipes, recipe{
		ID: 2,
		Name: "Steamed Rice",
		Author: "Billy",
		ImageUrl: "https://assets.marthastewart.com/styles/wmax-300/d9/06edf19_e/06edf19_e_vert.jpg?itok=gGqkyPnq",
		SourceUrl: "https://www.marthastewart.com/355435/best-steamed-rice",
	})

	// for rows.Next() {
	// 	var r recipe
	// 	if err := rows.Scan(&r.ID, &r.Name, &r.Author, &r.ImageUrl, &r.SourceUrl); err != nil {
	// 		return nil, err
	// 	}
	//
	// 	recipes = append(recipes, r)
	// }

	return recipes, nil
}
package models

import (
	"database/sql"
)

type Food struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type FoodForm struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required,min=1"`
}

type FoodListResponse struct {
	Status string `json:"status"`
	Data   []Food `json:"data"`
}

func GetFoodByID(db *sql.DB, id int) (*Food, error) {
	food := &Food{}
	err := db.QueryRow("SELECT id, name, description, price FROM foods WHERE id = ?", id).Scan(
		&food.ID, &food.Name, &food.Description, &food.Price)
	if err != nil {
		return nil, err
	}
	return food, nil
}

func (f *Food) Insert(db *sql.DB) error {
	res, err := db.Exec("INSERT INTO foods (name, description, price) VALUES (?, ?, ?)", f.Name, f.Description, f.Price)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	f.ID = int(lastID)

	return nil
}

func (f *Food) Update(db *sql.DB) error {
	_, err := db.Exec("UPDATE foods SET name = ?, description = ?, price = ? WHERE id = ?", f.Name, f.Description, f.Price, f.ID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFood(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM foods WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func GetFoodById(db *sql.DB, id int) (*Food, error) {
	var food Food
	err := db.QueryRow("SELECT id, name, description, price FROM foods WHERE id = ?", id).Scan(&food.ID, &food.Name, &food.Description, &food.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Food not found, return nil food and no error
		}
		return nil, err // Return nil food and the error
	}

	return &food, nil // Return the food and no error
}

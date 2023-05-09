package models

import (
	"gorm.io/gorm"
)

type Food struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// todos GetFood
func GetFood(db *gorm.DB) ([]*Food, error) {
	foods := []*Food{}
	result := db.Find(&foods)
	if result.Error != nil {
		return nil, result.Error
	}
	return foods, nil
}

// todo GetfoodById
func GetFoodById(db *gorm.DB, id uint) (*Food, error) {
	food := &Food{}
	result := db.First(food, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return food, nil
}

// todo GetFoodInDB
func CreateFoodInDB(db *gorm.DB, food *Food) error {
	result := db.Create(food)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// todo UpdateFoodInDB
func UpdateFoodInDB(db *gorm.DB, food *Food) error {
	result := db.Save(food)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// todo DeleteFoodInDB
func DeleteFoodInDB(db *gorm.DB, id int) error {
	result := db.Delete(&Food{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

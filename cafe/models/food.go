package models

import "gorm.io/gorm"

type Food struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// GetFoods mengembalikan daftar makanan yang tersedia di database.
func GetFoods(db *gorm.DB) ([]Food, error) {
	var foods []Food
	result := db.Find(&foods)
	if result.Error != nil {
		return nil, result.Error
	}
	return foods, nil
}

// GetFoodByID mengembalikan makanan dengan ID yang diberikan dari database.
func GetFoodByID(db *gorm.DB, id uint) (*Food, error) {
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

// CreateFood membuat dan menyimpan makanan baru ke database.
func CreateFood(db *gorm.DB, food *Food) error {
	result := db.Create(food)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateFood mengubah makanan yang ada di database.
func UpdateFood(db *gorm.DB, food *Food) error {
	result := db.Save(food)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteFood menghapus makanan dengan ID yang diberikan dari database.
func DeleteFood(db *gorm.DB, id int) error {
	result := db.Delete(&Food{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

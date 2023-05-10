package models

import "gorm.io/gorm"

type User struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Userrole string `json:"userrole"`
}

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type LoginCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// GetUserById gets a user by ID
func GetUserById(db *gorm.DB, id uint64) (*User, error) {
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // User not found, return nil user and no error
		}
		return nil, result.Error // Return nil user and the error
	}
	return &user, nil // Return the user and no error
}

// Fungsi untuk menghapus data user berdasarkan ID
func DeleteUser(db *gorm.DB, id int) error {
	err := db.Delete(&User{}, id).Error
	return err
}

// Fungsi untuk mengambil semua data user
func GetUsers(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Find(&users).Error
	return users, err
}

// Fungsi untuk mengambil data user berdasarkan ID
func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var user User
	err := db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Fungsi untuk membuat user baru
func CreateUser(db *gorm.DB, user *User) error {
	err := db.Create(user).Error
	return err
}

// Fungsi untuk mengupdate data user
func UpdateUser(db *gorm.DB, user *User) error {
	err := db.Save(user).Error
	return err
}

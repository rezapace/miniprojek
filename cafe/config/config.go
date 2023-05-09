package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Config struct {
	DB_Username string
	DB_Password string
	DB_Port     string
	DB_Host     string
	DB_Name     string
}

func LoadConfig() *Config {
	return &Config{
		DB_Username: "root",
		DB_Password: "",
		DB_Port:     "3306",
		DB_Host:     "localhost",
		DB_Name:     "cafe",
	}
}

func ConnectToDB(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DB_Username, config.DB_Password, config.DB_Host, config.DB_Port, config.DB_Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitDB() {
	config := LoadConfig()
	db, err := ConnectToDB(config)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	DB = db
}

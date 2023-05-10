package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// struktur data untuk konfigurasi database
type Config struct {
	DBUsername string
	DBPassword string
	DBPort     string
	DBHost     string
	DBName     string
}

// LoadConfig akan memuat konfigurasi database dari variabel lingkungan atau nilai yang dikodekan.
func LoadConfig() *Config {
	return &Config{
		DBUsername: "root",
		DBPassword: "",
		DBPort:     "3306",
		DBHost:     "localhost",
		DBName:     "cafe",
	}
}

// ConnectToDB untuk mengkonekan data yang telah di configurasikan
func ConnectToDB(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// InitDB meng inialisasi database
func InitDB() {
	config := LoadConfig()
	db, err := ConnectToDB(config)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	DB = db
}

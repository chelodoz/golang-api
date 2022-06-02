package database

import (
	"fmt"
	"golang-api/repository"
	"golang-api/util"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBConnection(config util.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)
	fmt.Print(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&repository.UserGorm{})

	return db, nil
}

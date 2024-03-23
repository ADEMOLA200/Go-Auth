package database

import (
	"github.com/ADEMOLA200/Go-Auth/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() error {
    dsn := "root:rootroot@tcp(127.0.0.1:3306)/go_auth?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return err
    }
    
    DB = db

    // AutoMigrate models
    if err := DB.AutoMigrate(&models.User{}); err != nil {
        return err
    }

    return nil
}
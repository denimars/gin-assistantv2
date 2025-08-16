package code

func Connection() string {
	return `
package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connection() *gorm.DB {
	godotenv.Load()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", os.Getenv("EXAMPLE_USER"), os.Getenv("EXAMPLE_PASSWORD"), os.Getenv("EXAMPLE_HOST"), os.Getenv("EXAMPLE_PORT"), os.Getenv("EXAMPLE_DB_NAME"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}

func init() {
	db := Connection()
	err := db.AutoMigrate()
	if err != nil {
		panic(err)
	}
}
`
}

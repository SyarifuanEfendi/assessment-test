package database

import (
	"assesment-test/api/logger"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() *gorm.DB {
	Db = connectDB()
	return Db
}

func connectDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		logger.LoggerFatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.LoggerFatal("Error Connect to Database")
		return nil
	}
	log.WithFields(log.Fields{}).Info("Successfully connected to database")
	return db
}

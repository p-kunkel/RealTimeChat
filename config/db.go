package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectToDB() error {
	var (
		err error
	)
	if DB, err = gorm.Open(postgres.Open(getDBAddres()), getDBConfig()); err != nil {
		return err
	}
	return nil
}

func getDBAddres() string {
	var (
		host       = os.Getenv("DB_ADDRESS")
		user       = os.Getenv("DB_LOGIN")
		password   = os.Getenv("DB_PASSWORD")
		dbName     = os.Getenv("DB_NAME")
		port       = os.Getenv("DB_PORT")
		tz         = os.Getenv("DB_TIMEZONE")
		searchPath = os.Getenv("DB_SEARCH_PATH")
	)

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s pg_trgm.similarity_threshold=0.02 search_path=%s", host, user, password, dbName, port, tz, searchPath)
}

func getDBConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
}

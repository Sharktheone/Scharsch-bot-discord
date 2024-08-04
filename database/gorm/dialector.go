package gorm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/url"
)

func GetDialector() gorm.Dialector {
	switch config.Provider {
	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			url.QueryEscape(config.User),
			url.QueryEscape(config.Pass),
			url.QueryEscape(config.Host),
			config.Port,
			url.QueryEscape(config.DatabaseName),
		)
		return mysql.Open(dsn)

	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%v",
			config.Host,
			config.User,
			config.Pass,
			config.DatabaseName,
			config.Port,
			config.TimeZone,
		)

		return postgres.Open(dsn)

	case "sqlite":
		return sqlite.Open(config.SqLiteFile)

	case "mongodb":
		log.Panicf("We probably messed up internally..., cannot connect to MongoDB using GORM, please open an issue on GitHub")

	default:
		log.Panicf("Unknown database name: %s\n Supported options are: `mysql`, `postgres`, `mysql`, `mongodb`", config.DatabaseName)
	}

	log.Panicf("We should never reach this point, please open an issue on GitHub")
	return nil
}

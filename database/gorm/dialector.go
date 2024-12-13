package gorm

import (
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/url"
)

func GetDialector() gorm.Dialector {
	switch conf.Config.Whitelist.Database.Provider {
	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			url.QueryEscape(conf.Config.Whitelist.Database.User),
			url.QueryEscape(conf.Config.Whitelist.Database.Pass),
			url.QueryEscape(conf.Config.Whitelist.Database.Host),
			conf.Config.Whitelist.Database.Port,
			url.QueryEscape(conf.Config.Whitelist.Database.DatabaseName),
		)
		return mysql.Open(dsn)

	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%v",
			conf.Config.Whitelist.Database.Host,
			conf.Config.Whitelist.Database.User,
			conf.Config.Whitelist.Database.Pass,
			conf.Config.Whitelist.Database.DatabaseName,
			conf.Config.Whitelist.Database.Port,
			conf.Config.Whitelist.Database.TimeZone,
		)

		return postgres.Open(dsn)

	case "sqlite":
		return sqlite.Open(conf.Config.Whitelist.Database.SqLiteFile)

	case "mongodb":
		log.Panicf("We probably messed up internally..., cannot connect to MongoDB using GORM, please open an issue on GitHub")

	default:
		log.Panicf("Unknown database name: %s\n Supported options are: `mysql`, `postgres`, `mysql`, `mongodb`", conf.Config.Whitelist.Database.DatabaseName)
	}

	log.Panicf("We should never reach this point, please open an issue on GitHub")
	return nil
}

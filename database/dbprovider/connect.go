package dbprovider

import (
	"context"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/database/gorm"
	"github.com/Sharktheone/ScharschBot/database/mongodb"
	"log"
)

func Connect() {

	switch conf.Config.Whitelist.Database.Provider {
	case "mongodb":
		database.DB = &mongodb.MongoConnection{
			Ctx: context.Background(),
		}

	case "mysql":
	case "sqlite":
	case "postgres":
		database.DB = &gorm.GormConnection{}

	default:
		log.Panicf("Unknown database provider: %s\n Supported options are: `mongodb`, `mysql`, `postgres`, `sqlite`", conf.Config.Whitelist.Database.Provider)
	}

	database.DB.Connect()
}

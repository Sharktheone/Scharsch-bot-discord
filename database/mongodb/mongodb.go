package mongodb

import (
	"context"
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/url"
	"time"
)

var (
	config             = conf.Config
	uri                = fmt.Sprintf("mongodb://%v:%v@%v:%v", url.QueryEscape(config.Whitelist.Database.User), url.QueryEscape(config.Whitelist.Database.Pass), config.Whitelist.Database.Host, config.Whitelist.Database.Port)
	Client, connectErr = mongo.NewClient(options.Client().ApplyURI(uri))
	Cancel             context.CancelFunc
	connected          = false
	db                 *mongo.Database
	Ready              = false
)

func Connect() *mongo.Database {
	if !connected {
		if connectErr != nil {
			log.Fatalf("Failed to apply mongo URI: %v", connectErr)
		}
		var ctx context.Context
		ctx, Cancel = context.WithTimeout(context.Background(), 10*time.Second)
		err := Client.Connect(ctx)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		} else {
			connected = true
		}
	}
	err := Client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB")
	db = Client.Database(config.Whitelist.Database.DatabaseName)
	Ready = true
	log.Printf("Using datbase %v", db.Name())

	return db

}

func Disconnect() {
	var ctx context.Context
	ctx, Cancel = context.WithTimeout(context.Background(), 10*time.Second)
	err := Client.Disconnect(ctx)
	if err != nil {
		log.Printf("Failed to disconnect: %v \n", err)
	}
}

func Write(collection string, data bson.D) {
	writeColl := db.Collection(collection)
	var ctx context.Context
	ctx, Cancel = context.WithTimeout(context.Background(), 10*time.Second)
	_, err := writeColl.InsertOne(ctx, data)
	if err != nil {
		log.Printf("Failed to write data: %v \n", err)
	}
}
func Read(collection string, filter bson.M) (data []bson.M, found bool) {
	readColl := db.Collection(collection)
	var (
		ctx       context.Context
		dataFound bool
	)
	ctx, Cancel = context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := readColl.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to find data: %v \n", err)
	}
	var read []bson.M
	if err = cursor.All(ctx, &read); err != nil {
		log.Printf("Failed to read mongo Cursor: %v \n", err)
	}
	dataFound = len(read) > 0
	return read, dataFound
}

func Remove(collection string, filter bson.M) *mongo.DeleteResult {
	removeColl := db.Collection(collection)
	var ctx context.Context
	ctx, Cancel = context.WithTimeout(context.Background(), 10*time.Second)
	result, err := removeColl.DeleteMany(ctx, filter)
	if err != nil {
		log.Printf("Failed to remove data: %v \n", err)
	}
	return result
}

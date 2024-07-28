package mongodb

import (
	"context"
	"fmt"
	"github.com/Sharktheone/ScharschBot/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/url"
	"time"
)

type MongoConnection struct {
	ctx       context.Context
	db        *mongo.Database
	connected bool
}

func (m *MongoConnection) Connect() {

	uri := fmt.Sprintf(
		"mongodb://%v:%v@%v:%v",
		url.QueryEscape(config.Whitelist.Database.User),
		url.QueryEscape(config.Whitelist.Database.Pass),
		config.Whitelist.Database.Host,
		config.Whitelist.Database.Port,
	)

	client, err := mongo.Connect(m.ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to apply mongo URI: %v", err)
	}

	err = client.Ping(m.ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB")

	m.db = client.Database(config.Whitelist.Database.DatabaseName)
	m.connected = true
}

func (m *MongoConnection) Disconnect() {
	ctx, cancel := context.WithTimeout(m.ctx, 10*time.Second)
	defer cancel()
	m.connected = false
	err := m.db.Client().Disconnect(ctx)
	if err != nil {
		log.Fatalf("Failed to disconnect: %v \n", err)
	}
}

func (m *MongoConnection) Read(collection string, filter bson.M) (*mongo.Cursor, error) {
	coll := db.Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return coll.Find(ctx, filter)
}

func (m *MongoConnection) Write(collection string, data interface{}) {
	writeColl := db.Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := writeColl.InsertOne(ctx, data)
	if err != nil {
		log.Printf("Failed to write to MongoDB: %v", err)
	}
}

func (m *MongoConnection) Remove(collection string, filter bson.M) {
	writeColl := db.Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := writeColl.DeleteMany(ctx, filter)
	if err != nil {
		log.Printf("Failed to remove from MongoDB: %v", err)
	}
}

func (m *MongoConnection) WhitelistPlayer(user database.UserID, player database.Player, roles []database.Role) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) UnWhitelistAccount(user database.UserID) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) UnWhitelistPlayer(player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) MoveToReWhitelist(user database.UserID) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) ReWhitelist(user database.UserID, roles []database.Role) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) RemoveAll() {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) RemoveAllFrom(user database.UserID) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) Owner(player database.Player) database.UserID {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) Players(user database.UserID) []database.Player {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) BanUser(user database.UserID, reason string) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) BanPlayer(user database.UserID, player database.Player, reason string) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) UnBanUser(user database.UserID) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) UnBanPlayer(player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) UnBanPlayerFrom(user database.UserID, player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) IsUserBanned(user database.UserID) bool {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) IsPlayerBanned(player database.Player) bool {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) BannedPlayers(user database.UserID) []database.Player {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) RemoveAccounts(user database.UserID) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) AddWaitlist(user database.UserID, player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) RemoveWaitlist(user database.UserID, player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) Report(reporter database.UserID, reported database.Player, reason string) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) AllWhitelists() []database.PlayerData {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) AllReWhitelists() []database.PlayerData {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) DeleteReport(reported database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) RemoveAccount(player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) IsWhitelisted(player database.Player) bool {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) IsWhitelistedBy(user database.UserID, player database.Player) bool {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) GetReports() []database.ReportData {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) IsAlreadyReported(reported database.Player) bool {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) GetReportedPlayer(reported database.Player) (database.ReportData, bool) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) NumberWhitelistedPlayers(user database.UserID) int {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) GetWhitelistedPlayer(player database.Player) (database.WhitelistedPlayerData, bool) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) GetAllWhitelistedPlayers() []database.WhitelistedPlayerData {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) GetAccountsOf(user database.UserID) []database.Player {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) GetBan(user database.UserID) (string, bool) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoConnection) GetPlayerBan(player database.Player) (database.PlayerBan, bool) {
	//TODO implement me
	panic("implement me")
}

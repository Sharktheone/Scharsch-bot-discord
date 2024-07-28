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
	"slices"
	"time"
)

var (
	whitelistCollection   = config.Whitelist.Database.WhitelistTable
	banCollection         = config.Whitelist.Database.BanTable
	reWhitelistCollection = config.Whitelist.Database.ReWhitelistTable
	reportCollection      = config.Whitelist.Database.ReportTable
	waitlistCollection    = config.Whitelist.Database.WaitListTable
	rolesCollection       = config.Whitelist.Database.RolesTable
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

func (m *MongoConnection) WhitelistPlayer(user database.UserID, player database.Player) {
	entry := database.WhitelistEntry{
		UserID: user,
		Player: player,
	}

	m.Write(whitelistCollection, entry)
}

func (m *MongoConnection) UnWhitelistAccount(user database.UserID) {
	m.Remove(whitelistCollection, bson.M{"userID": user})
}

func (m *MongoConnection) UnWhitelistPlayer(player database.Player) {
	m.Remove(whitelistCollection, bson.M{"player": player})
}

func (m *MongoConnection) MoveToReWhitelist(user database.UserID, missingRole database.Role) {
	whitelisted := m.Players(user)

	entry := database.ReWhitelsitEntry{
		UserID:      user,
		Players:     whitelisted,
		MissingRole: missingRole,
	}

	m.Write(reWhitelistCollection, entry)
}

func (m *MongoConnection) ReWhitelist(user database.UserID, roles []database.Role) {
	cursor, err := m.Read(reWhitelistCollection, bson.M{"userID": user})
	if err != nil {
		log.Printf("Failed to rewhitelist: %v", err)
	}

	var reEntry database.ReWhitelsitEntry
	if err := cursor.Decode(&reEntry); err != nil {
		log.Printf("Failed to rewhitelist: %v", err)
	}

	if slices.Contains(roles, reEntry.MissingRole) {
		for _, player := range reEntry.Players {
			m.WhitelistPlayer(user, player)
		}
		m.Remove(reWhitelistCollection, bson.M{"userID": user})
	}

}

func (m *MongoConnection) RemoveAll() {
	m.Remove(whitelistCollection, bson.M{})
}

func (m *MongoConnection) RemoveAllFrom(user database.UserID) {
	m.Remove(whitelistCollection, bson.M{"userID": user})
}

func (m *MongoConnection) Owner(player database.Player) database.UserID {
	cursor, err := m.Read(whitelistCollection, bson.M{"player": player})
	if err != nil {
		log.Printf("Failed to get owner: %v", err)
	}

	var entry database.WhitelistEntry
	if err := cursor.Decode(&entry); err != nil {
		log.Printf("Failed to get owner: %v", err)
	}

	return entry.UserID
}

func (m *MongoConnection) Players(user database.UserID) []database.Player {
	cursor, err := m.Read(whitelistCollection, bson.M{"userID": user})
	if err != nil {
		log.Printf("Failed to get players: %v", err)
	}

	var entries []database.Player
	for cursor.Next(m.ctx) {
		var entry database.WhitelistEntry
		if err := cursor.Decode(&entry); err != nil {
			log.Printf("Failed to get players: %v", err)
		}
		entries = append(entries, entry.Player)
	}

	return entries
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

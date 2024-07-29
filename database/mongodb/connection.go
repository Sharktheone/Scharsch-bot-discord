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

	entry := database.ReWhitelistEntry{
		UserID:      user,
		Players:     whitelisted,
		MissingRole: missingRole,
	}

	m.Write(reWhitelistCollection, entry)

	m.Remove(whitelistCollection, bson.M{"userID": user})
}

func (m *MongoConnection) ReWhitelist(user database.UserID, roles []database.Role) {
	cursor, err := m.Read(reWhitelistCollection, bson.M{"userID": user})
	if err != nil {
		log.Printf("Failed to rewhitelist: %v", err)
	}

	var reEntry database.ReWhitelistEntry
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

	players := m.Players(user)

	playerReasons := make([]database.PlayerBanData, len(players))
	for i, p := range players {
		playerReasons[i] = database.PlayerBanData{Player: p, Reason: reason}
	}

	entry := database.BanEntry{
		UserID:     user,
		UserBan:    true,
		UserReason: reason,
		Players:    playerReasons,
	}

	m.Write(banCollection, entry)
	m.Remove(whitelistCollection, bson.M{"userID": user})
}

func (m *MongoConnection) BanPlayer(player database.Player, reason string) {
	owner := m.Owner(player)
	entry := database.BanEntry{
		UserID: owner,
		Players: []database.PlayerBanData{
			{Player: player, Reason: reason},
		},
	}

	m.Write(banCollection, entry)
	m.Remove(whitelistCollection, bson.M{"player": player})

}

func (m *MongoConnection) UnBanUser(user database.UserID) {
	m.Remove(banCollection, bson.M{"userID": user})
}

func (m *MongoConnection) UnBanPlayer(player database.Player) {
	cursor, err := m.Read(banCollection, bson.M{"players": player})
	if err != nil {
		log.Printf("Failed to unban player: %v", err)
		return
	}

	var entry database.BanEntry
	if err := cursor.Decode(&entry); err != nil {
		log.Printf("Failed to unban player: %v", err)
		return
	}

	if len(entry.Players) == 1 {
		m.Remove(banCollection, bson.M{"userID": entry.UserID})
	}

	m.Remove(banCollection, bson.M{"players": player})
	// remove player from ban array
	for i, p := range entry.Players {
		if p.Player == player {
			entry.Players = append(entry.Players[:i], entry.Players[i+1:]...)
			break
		}
	}

	m.Write(banCollection, entry)
}

func (m *MongoConnection) UnBanPlayerFrom(user database.UserID, player database.Player) {
	cursor, err := m.Read(banCollection, bson.M{"userID": user})
	if err != nil {
		log.Printf("Failed to unban player: %v", err)
		return
	}

	var entry database.BanEntry
	if err := cursor.Decode(&entry); err != nil {
		log.Printf("Failed to unban player: %v", err)
		return
	}

	if len(entry.Players) == 1 {
		m.Remove(banCollection, bson.M{"userID": entry.UserID})
	}

	for i, p := range entry.Players {
		if p.Player == player {
			entry.Players = append(entry.Players[:i], entry.Players[i+1:]...)
			break
		}
	}

	m.Write(banCollection, entry)
}

func (m *MongoConnection) IsUserBanned(user database.UserID) bool {
	cursor, err := m.Read(banCollection, bson.M{"userID": user})
	if err != nil {
		log.Printf("Failed to check if user is banned: %v", err)
		return false
	}

	if cursor.RemainingBatchLength() == 0 {
		return false
	}

	var entry database.BanEntry
	if err := cursor.Decode(&entry); err != nil {
		log.Printf("Failed to check if user is banned: %v", err)
		return false
	}

	return entry.UserBan

}

func (m *MongoConnection) IsPlayerBanned(player database.Player) bool {
	cursor, err := m.Read(banCollection, bson.M{"players": player})
	if err != nil {
		log.Printf("Failed to check if player is banned: %v", err)
		return false
	}

	return cursor.RemainingBatchLength() != 0
}

func (m *MongoConnection) BannedPlayers(user database.UserID) []database.PlayerBanData {
	cursor, err := m.Read(banCollection, bson.M{"userID": user})
	if err != nil {
		log.Printf("Failed to get banned players: %v", err)
	}

	if cursor.RemainingBatchLength() == 0 {
		return []database.PlayerBanData{}
	}

	var entry database.BanEntry
	if err := cursor.Decode(&entry); err != nil {
		log.Printf("Failed to get banned players: %v", err)
	}

	return entry.Players
}

func (m *MongoConnection) RemoveAccounts(user database.UserID) {
	m.Remove(whitelistCollection, bson.M{"userID": user})
}

func (m *MongoConnection) AddWaitlist(user database.UserID, player database.Player) {
	entry := database.WaitlistEntry{
		UserID: user,
		Player: player,
	}

	m.Write(waitlistCollection, entry)
}

func (m *MongoConnection) RemoveWaitlist(user database.UserID, player database.Player) {
	m.Remove(waitlistCollection, bson.M{"userID": user, "player": player})
}

func (m *MongoConnection) Report(reporter database.UserID, reported database.Player, reason string) {
	entry := database.ReportEntry{
		ReporterID:     reporter,
		ReportedPlayer: reported,
		Reason:         reason,
	}

	m.Write(reportCollection, entry)
}

func (m *MongoConnection) AllWhitelists() []database.PlayerData {
	cursor, err := m.Read(whitelistCollection, bson.M{})
	if err != nil {
		log.Printf("Failed to get all whitelists: %v", err)
	}

	var entries []database.PlayerData
	if err := cursor.Decode(&entries); err != nil {
		log.Printf("Failed to get all whitelists: %v", err)
		return entries
	}

	return entries
}

func (m *MongoConnection) AllReWhitelists() []database.PlayerData {
	cursor, err := m.Read(reWhitelistCollection, bson.M{})
	if err != nil {
		log.Printf("Failed to get all rewhitelists: %v", err)
	}

	var entries []database.PlayerData
	if err := cursor.Decode(&entries); err != nil {
		log.Printf("Failed to get all rewhitelists: %v", err)
		return entries
	}

	return entries
}

func (m *MongoConnection) DeleteReport(reported database.Player) {
	m.Remove(reportCollection, bson.M{"reportedPlayer": reported})
}

func (m *MongoConnection) RemoveAccount(player database.Player) {
	m.Remove(whitelistCollection, bson.M{"player": player})
}

func (m *MongoConnection) IsWhitelisted(player database.Player) bool {
	cursor, err := m.Read(whitelistCollection, bson.M{"player": player})
	if err != nil {
		log.Printf("Failed to check if player is whitelisted: %v", err)
		return false
	}

	return cursor.RemainingBatchLength() != 0
}

func (m *MongoConnection) IsWhitelistedBy(user database.UserID, player database.Player) bool {
	cursor, err := m.Read(whitelistCollection, bson.M{"userID": user, "player": player})
	if err != nil {
		log.Printf("Failed to check if player is whitelisted by user: %v", err)
		return false
	}

	return cursor.RemainingBatchLength() != 0
}

func (m *MongoConnection) GetReports() []database.ReportData {
	cursor, err := m.Read(reportCollection, bson.M{})
	if err != nil {
		log.Printf("Failed to get reports: %v", err)
	}

	var entries []database.ReportData
	if err := cursor.Decode(&entries); err != nil {
		log.Printf("Failed to get reports: %v", err)
		return entries
	}

	return entries
}

func (m *MongoConnection) IsAlreadyReported(reported database.Player) bool {
	cursor, err := m.Read(reportCollection, bson.M{"reportedPlayer": reported})
	if err != nil {
		log.Printf("Failed to check if player is already reported: %v", err)
		return false
	}

	return cursor.RemainingBatchLength() != 0
}

func (m *MongoConnection) GetReportedPlayer(reported database.Player) (database.ReportData, bool) {
	cursor, err := m.Read(reportCollection, bson.M{"reportedPlayer": reported})
	if err != nil {
		log.Printf("Failed to get reported player: %v", err)
		return database.ReportData{}, false
	}

	var entry database.ReportData
	if err := cursor.Decode(&entry); err != nil {
		log.Printf("Failed to get reported player: %v", err)
		return database.ReportData{}, false
	}

	return entry, true
}

func (m *MongoConnection) NumberWhitelistedPlayers(user database.UserID) int {
	cursor, err := m.Read(whitelistCollection, bson.M{"userID": user})
	if err != nil {
		log.Printf("Failed to get number of whitelisted players: %v", err)
		return 0
	}

	return cursor.RemainingBatchLength()
}

func (m *MongoConnection) GetWhitelistedPlayer(player database.Player) (database.WhitelistedPlayerData, bool) {
	cursor, err := m.Read(whitelistCollection, bson.M{"player": player})
	if err != nil {
		log.Printf("Failed to get whitelisted player: %v", err)
		return database.WhitelistedPlayerData{}, false
	}

	var entry database.WhitelistedPlayerData
	if err := cursor.Decode(&entry); err != nil {
		log.Printf("Failed to get whitelisted player: %v", err)
		return database.WhitelistedPlayerData{}, false
	}

	return entry, true
}

func (m *MongoConnection) GetAllWhitelistedPlayers() []database.WhitelistedPlayerData {
	cursor, err := m.Read(whitelistCollection, bson.M{})
	if err != nil {
		log.Printf("Failed to get all whitelisted players: %v", err)
	}

	var entries []database.WhitelistedPlayerData
	if err := cursor.Decode(&entries); err != nil {
		log.Printf("Failed to get all whitelisted players: %v", err)
		return entries
	}

	return entries
}

func (m *MongoConnection) GetAccountsOf(user database.UserID) []database.Player {
	cursor, err := m.Read(whitelistCollection, bson.M{"userID": user})
	if err != nil {
		log.Printf("Failed to get accounts of user: %v", err)
	}

	var entries []database.Player
	if err := cursor.Decode(&entries); err != nil {
		log.Printf("Failed to get accounts of user: %v", err)
		return entries
	}

	return entries
}

func (m *MongoConnection) GetBan(user database.UserID) (string, bool) {
	cursor, err := m.Read(banCollection, bson.M{"userID": user})
	if err != nil {
		log.Printf("Failed to get ban: %v", err)
		return "", false
	}

	var entry database.BanEntry
	if err := cursor.Decode(&entry); err != nil {
		log.Printf("Failed to get ban: %v", err)
		return "", false
	}

	if entry.UserBan {
		return entry.UserReason, true
	}

	return "", false
}

func (m *MongoConnection) GetPlayerBan(player database.Player) (database.PlayerBan, bool) {
	cursor, err := m.Read(banCollection, bson.M{"players": player})
	if err != nil {
		log.Printf("Failed to get player ban: %v", err)
		return database.PlayerBan{}, false
	}

	var entry database.BanEntry
	if err := cursor.Decode(&entry); err != nil {
		log.Printf("Failed to get player ban: %v", err)
		return database.PlayerBan{}, false
	}

	for _, p := range entry.Players {
		if p.Player == player {
			return database.PlayerBan{ID: entry.UserID, Reason: entry.UserReason}, true
		}
	}

	return database.PlayerBan{}, false
}

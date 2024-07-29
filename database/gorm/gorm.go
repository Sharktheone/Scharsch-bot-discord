package gorm

import (
	"github.com/Sharktheone/ScharschBot/database"
	"gorm.io/gorm"
	"log"
)

type GormConnection struct {
	DB *gorm.DB
}

func (m *GormConnection) Connect() {
	var err error
	m.DB, err = gorm.Open(GetDialector())
	if err != nil {
		log.Panicf("Failed to connect to database: %v", err)
	}

	if err := m.DB.AutoMigrate(&WhitelistEntry{}, &BanEntry{}, &ReWhitelistEntry{}, &ReportEntry{}, &WaitlistEntry{}); err != nil {
		log.Panicf("Failed to migrate database: %v", err)
	}
}

func (m *GormConnection) Disconnect() {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) WhitelistPlayer(user database.UserID, player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) UnWhitelistAccount(user database.UserID) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) UnWhitelistPlayer(player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) MoveToReWhitelist(user database.UserID, missingRole database.Role) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) ReWhitelist(user database.UserID, roles []database.Role) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) RemoveAll() {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) RemoveAllFrom(user database.UserID) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) Owner(player database.Player) database.UserID {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) Players(user database.UserID) []database.Player {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) BanUser(user database.UserID, reason string) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) BanPlayer(player database.Player, reason string) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) UnBanUser(user database.UserID) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) UnBanPlayer(player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) UnBanPlayerFrom(user database.UserID, player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) IsUserBanned(user database.UserID) bool {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) IsPlayerBanned(player database.Player) bool {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) BannedPlayers(user database.UserID) []database.PlayerBanData {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) RemoveAccounts(user database.UserID) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) AddWaitlist(user database.UserID, player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) RemoveWaitlist(user database.UserID, player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) Report(reporter database.UserID, reported database.Player, reason string) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) AllWhitelists() []database.PlayerData {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) AllReWhitelists() []database.PlayerData {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) DeleteReport(reported database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) RemoveAccount(player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) IsWhitelisted(player database.Player) bool {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) IsWhitelistedBy(user database.UserID, player database.Player) bool {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) GetReports() []database.ReportData {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) IsAlreadyReported(reported database.Player) bool {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) GetReportedPlayer(reported database.Player) (database.ReportData, bool) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) NumberWhitelistedPlayers(user database.UserID) int {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) GetWhitelistedPlayer(player database.Player) (database.WhitelistedPlayerData, bool) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) GetAllWhitelistedPlayers() []database.WhitelistedPlayerData {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) GetAccountsOf(user database.UserID) []database.Player {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) GetBan(user database.UserID) (string, bool) {
	//TODO implement me
	panic("implement me")
}

func (m *GormConnection) GetPlayerBan(player database.Player) (database.PlayerBan, bool) {
	//TODO implement me
	panic("implement me")
}

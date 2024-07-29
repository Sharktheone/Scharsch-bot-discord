package mariadb

import "github.com/Sharktheone/ScharschBot/database"

type MariaDBConnection struct {
}

func (m *MariaDBConnection) Connect() {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) Disconnect() {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) WhitelistPlayer(user database.UserID, player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) UnWhitelistAccount(user database.UserID) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) UnWhitelistPlayer(player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) MoveToReWhitelist(user database.UserID, missingRole database.Role) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) ReWhitelist(user database.UserID, roles []database.Role) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) RemoveAll() {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) RemoveAllFrom(user database.UserID) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) Owner(player database.Player) database.UserID {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) Players(user database.UserID) []database.Player {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) BanUser(user database.UserID, reason string) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) BanPlayer(player database.Player, reason string) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) UnBanUser(user database.UserID) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) UnBanPlayer(player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) UnBanPlayerFrom(user database.UserID, player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) IsUserBanned(user database.UserID) bool {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) IsPlayerBanned(player database.Player) bool {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) BannedPlayers(user database.UserID) []database.PlayerBanData {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) RemoveAccounts(user database.UserID) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) AddWaitlist(user database.UserID, player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) RemoveWaitlist(user database.UserID, player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) Report(reporter database.UserID, reported database.Player, reason string) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) AllWhitelists() []database.PlayerData {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) AllReWhitelists() []database.PlayerData {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) DeleteReport(reported database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) RemoveAccount(player database.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) IsWhitelisted(player database.Player) bool {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) IsWhitelistedBy(user database.UserID, player database.Player) bool {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) GetReports() []database.ReportData {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) IsAlreadyReported(reported database.Player) bool {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) GetReportedPlayer(reported database.Player) (database.ReportData, bool) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) NumberWhitelistedPlayers(user database.UserID) int {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) GetWhitelistedPlayer(player database.Player) (database.WhitelistedPlayerData, bool) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) GetAllWhitelistedPlayers() []database.WhitelistedPlayerData {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) GetAccountsOf(user database.UserID) []database.Player {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) GetBan(user database.UserID) (string, bool) {
	//TODO implement me
	panic("implement me")
}

func (m *MariaDBConnection) GetPlayerBan(player database.Player) (database.PlayerBan, bool) {
	//TODO implement me
	panic("implement me")
}

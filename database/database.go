package database

type UserID string
type Player string

type DatabaseConnection interface {
	Connect()
	Disconnect()
	WhitelistAccount(user UserID, player Player)
	UnWhitelistAccount(id string)
	RemoveAll()
	Owner(player Player) UserID
	AccountsOf(user UserID) []Player
	BanUser(user UserID, reason string)
	BanPlayer(player Player, reason string)
	UnBanUser(user UserID)
	UnBanPlayer(player Player)
	IsUserBanned(user UserID) bool
	IsPlayerBanned(player Player) bool
	BannedPlayers(user UserID) []Player
	RemoveAccounts(user UserID)
	AddWaitlist(user UserID, player Player)
	RemoveWaitlist(user UserID, player Player)
}

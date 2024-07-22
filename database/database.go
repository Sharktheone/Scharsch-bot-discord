package database

type UserID string
type Player string
type Role string

var (
	DB Connection
)

type PlayerData struct {
	ID      UserID
	Players []Player
	Roles   []Role
}

type Connection interface {
	Connect()
	Disconnect()
	WhitelistPlayer(user UserID, player Player, roles []Role)
	UnWhitelistAccount(user UserID)
	UnWhitelistPlayer(player Player)
	MoveToReWhitelist(user UserID)
	ReWhitelist(user UserID, roles []Role)
	RemoveAll()
	RemoveAllFrom(user UserID)
	Owner(player Player) UserID
	AccountsOf(user UserID) []Player
	BanUser(user UserID, reason string)
	BanPlayer(user UserID, player Player, reason string)
	UnBanUser(user UserID)
	UnBanPlayer(player Player)
	IsUserBanned(user UserID) bool
	IsPlayerBanned(player Player) bool
	BannedPlayers(user UserID) []Player
	RemoveAccounts(user UserID)
	AddWaitlist(user UserID, player Player)
	RemoveWaitlist(user UserID, player Player)
	Report(reporter UserID, reported Player, reason string)
	AllWhitelists() []PlayerData
	AllReWhitelists() []PlayerData
	DeleteReport(reported Player)
}

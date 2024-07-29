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

type WhitelistedPlayerData struct {
	ID   UserID
	Name Player
}

type ReportData struct {
	ReporterID     UserID
	ReportedPlayer Player
	Reason         string
}

type PlayerBan struct {
	ID     UserID
	Reason string
}

type Connection interface {
	Connect()
	Disconnect()
	WhitelistPlayer(user UserID, player Player)
	UnWhitelistAccount(user UserID)
	UnWhitelistPlayer(player Player)
	MoveToReWhitelist(user UserID, missingRole Role)
	ReWhitelist(user UserID, roles []Role)
	RemoveAll()
	RemoveAllFrom(user UserID)
	Owner(player Player) UserID
	Players(user UserID) []Player
	BanUser(user UserID, reason string)
	BanPlayer(player Player, reason string)
	UnBanUser(user UserID)
	UnBanPlayer(player Player)
	UnBanPlayerFrom(user UserID, player Player)
	IsUserBanned(user UserID) bool
	IsPlayerBanned(player Player) bool
	BannedPlayers(user UserID) []PlayerBanData
	RemoveAccounts(user UserID)
	AddWaitlist(user UserID, player Player)
	RemoveWaitlist(user UserID, player Player)
	Report(reporter UserID, reported Player, reason string)
	AllWhitelists() []WhitelistEntry
	AllReWhitelists() []ReWhitelistEntry
	DeleteReport(reported Player)
	RemoveAccount(player Player)
	IsWhitelisted(player Player) bool
	IsWhitelistedBy(user UserID, player Player) bool
	GetReports() []ReportData
	IsAlreadyReported(reported Player) bool
	GetReportedPlayer(reported Player) (ReportData, bool)
	NumberWhitelistedPlayers(user UserID) int
	GetWhitelistedPlayer(player Player) (WhitelistedPlayerData, bool)
	GetAllWhitelistedPlayers() []WhitelistedPlayerData
	GetAccountsOf(user UserID) []Player
	GetBan(user UserID) (string, bool)
	GetPlayerBan(player Player) (PlayerBan, bool)
}

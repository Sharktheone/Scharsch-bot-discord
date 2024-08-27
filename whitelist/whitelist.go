package whitelist

import (
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/types"
)

type WhitelistProvider interface {
	AddToWhitelist(player database.Player, member *types.Member)
	RemoveFromWhitelist(user database.UserID, player database.Player)
	MoveToReWhitelist(user database.UserID, missingRole database.Role)
	UnWhitelistAccount(user database.UserID)
	UnWhitelistPlayer(player database.Player)
	BanUser(user database.UserID, reason string)
	BanPlayer(user database.UserID, player database.Player, reason string)
	UnBanUser(user database.UserID)
	UnBanPlayer(player database.Player)
	UnBanPlayerFrom(user database.UserID, player database.Player)
	RemoveAccounts(user database.UserID)
	RemoveAccount(player database.Player)
}

var Provider WhitelistProvider = nil

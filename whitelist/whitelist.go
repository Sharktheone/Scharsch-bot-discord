package whitelist

import (
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/types"
)

type WhitelistProvider interface {
	AddToWhitelist(player database.Player, member *types.Member)
	UnWhitelistPlayer(player database.Player, member *types.Member)
	MoveToReWhitelist(missingRole database.Role, member *types.Member)
	UnWhitelistAccount(member *types.Member)
	UnWhitelistAccounts(members []*types.Member)
	BanUser(member *types.Member, reason string)
	BanPlayer(player database.Player, member *types.Member, reason string)
	UnBanUser(user database.UserID)
	UnBanPlayer(player database.Player)
	UnBanPlayerFrom(user database.UserID, player database.Player)
}

var Provider WhitelistProvider

func SetupProvider() {
	Provider = GetDefaultProvider()
}

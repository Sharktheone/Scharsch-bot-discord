package whitelist

import (
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/types"
	"github.com/Sharktheone/ScharschBot/whitelist/server"
)

var (
	config        = conf.Config.Whitelist
	AddCommand    = config.AddCommand
	RemoveCommand = config.RemoveCommand
	BanCommand    = config.BanCommand
	Roles         = config.RolesConfig
)

type DefaultProvider struct {
	ServerProvider server.ServerProvider
}

func (p *DefaultProvider) AddToWhitelist(player database.Player, member *types.Member) {
	database.DB.WhitelistPlayer(member.ID, player)

	if p.ServerProvider != nil {

		for _, serverID := range p.ServerProvider.GetServers() {
			getWhitelistCommand(member, serverID)
		}

	}
}
func (p *DefaultProvider) RemoveFromWhitelist(user database.UserID, player database.Player) {
	database.DB.UnWhitelistPlayer(player)

	if p.ServerProvider != nil {
		for _, serverID := range p.ServerProvider.GetServers() {
			getUnWhitelistCommand(player, serverID)
		}
	}

}
func (p *DefaultProvider) MoveToReWhitelist(user database.UserID, missingRole database.Role) {

}
func (p *DefaultProvider) UnWhitelistAccount(user database.UserID) {

}
func (p *DefaultProvider) UnWhitelistPlayer(player database.Player) {

}
func (p *DefaultProvider) BanUser(user database.UserID, reason string) {

}
func (p *DefaultProvider) BanPlayer(user database.UserID, player database.Player, reason string) {

}
func (p *DefaultProvider) UnBanUser(user database.UserID) {

}
func (p *DefaultProvider) UnBanPlayer(player database.Player) {

}
func (p *DefaultProvider) UnBanPlayerFrom(user database.UserID, player database.Player) {

}
func (p *DefaultProvider) RemoveAccounts(user database.UserID) *[]database.Player {

}

func (p *DefaultProvider) RemoveAccount(player database.Player) {

}

func getWhitelistCommand(member *types.Member, serverID server.ServerID) {

}

func getUnWhitelistCommand(member *types.Member, serverID server.ServerID) {

}

func getBanCommand(member *types.Member, serverID server.ServerID) {

}

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
			command := getWhitelistCommand(member, serverID)

			if command == "" || command == "<default>" {
				p.ServerProvider.Whitelist(player, serverID)
			}

			if command == "<none>" {
				continue
			}

			p.ServerProvider.SendCommand(command, serverID)
		}

	}
}
func (p *DefaultProvider) UnWhitelistPlayer(player database.Player, member *types.Member) {
	database.DB.UnWhitelistPlayer(player)

	if p.ServerProvider != nil {
		for _, serverID := range p.ServerProvider.GetServers() {
			command := getUnWhitelistCommand(member, serverID)

			if command == "" || command == "<default>" {
				p.ServerProvider.UnWhitelist(player, serverID)
			}

			if command == "<none>" {
				continue
			}

			p.ServerProvider.SendCommand(command, serverID)
		}
	}

}
func (p *DefaultProvider) MoveToReWhitelist(missingRole database.Role, member *types.Member) {

}
func (p *DefaultProvider) UnWhitelistAccount(member *types.Member) {
	players := database.DB.Players(member.ID)

	for _, player := range players {
		p.UnWhitelistPlayer(player, member)
	}
}
func (p *DefaultProvider) BanUser(member *types.Member, reason string) {
	players := database.DB.Players(member.ID)

	for _, player := range players {
		p.BanPlayer(player, member, reason)
	}
}
func (p *DefaultProvider) BanPlayer(player database.Player, member *types.Member, reason string) {
	database.DB.BanPlayer(player, reason)

	if p.ServerProvider != nil {
		for _, serverID := range p.ServerProvider.GetServers() {
			command := getBanCommand(member, serverID)

			if command == "" || command == "<default>" {
				p.ServerProvider.Ban(player, reason, serverID)
			}

			if command == "<none>" {
				continue
			}

			p.ServerProvider.SendCommand(command, serverID)
		}
	}
}
func (p *DefaultProvider) UnBanUser(user database.UserID) {
	database.DB.UnBanUser(user)

}
func (p *DefaultProvider) UnBanPlayer(player database.Player) {
	database.DB.UnBanPlayer(player)
}
func (p *DefaultProvider) UnBanPlayerFrom(user database.UserID, player database.Player) {
	database.DB.UnBanPlayerFrom(user, player)
}
func (p *DefaultProvider) RemoveAccounts(user database.UserID) *[]database.Player {
	players := database.DB.Players(user)

	for _, player := range players {
		database.DB.UnWhitelistPlayer(player)
	}

	database.DB.RemoveAccounts(user)
	return &players
}

func getWhitelistCommand(member *types.Member, serverID server.ServerID) string {

}

func getUnWhitelistCommand(member *types.Member, serverID server.ServerID) string {

}

func getBanCommand(member *types.Member, serverID server.ServerID) string {

}

func GetDefaultProvider() WhitelistProvider {
	return &DefaultProvider{}
}

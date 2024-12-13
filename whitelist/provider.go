package whitelist

import (
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/types"
	"github.com/Sharktheone/ScharschBot/whitelist/server"
	"github.com/Sharktheone/ScharschBot/whitelist/whitelist/utils"
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

func (p *DefaultProvider) RemoveAccount(player database.Player) {
	database.DB.UnWhitelistPlayer(player)
}

func getWhitelistCommand(member *types.Member, serverID server.ServerID) string {
	command := conf.Config.Whitelist.AddCommand

	for _, rc := range conf.Config.Whitelist.RolesConfig {
		if utils.CheckRole(member, rc.RoleID) {
			if rc.WhitelistCommand != "" {
				command = rc.WhitelistCommand
			}

			for _, s := range rc.PerServer {
				if s.ServerID == serverID {
					if s.WhitelistCommand != "" {
						command = s.WhitelistCommand
					}
				}
			}
		}
	}

	return command
}

func getUnWhitelistCommand(member *types.Member, serverID server.ServerID) string {
	command := conf.Config.Whitelist.RemoveCommand

	for _, rc := range conf.Config.Whitelist.RolesConfig {
		if utils.CheckRole(member, rc.RoleID) {
			if rc.UnWhitelistCommand != "" {
				command = rc.UnWhitelistCommand
			}

			for _, s := range rc.PerServer {
				if s.ServerID == serverID {
					if s.UnWhitelistCommand != "" {
						command = s.UnWhitelistCommand
					}
				}
			}
		}
	}

	return command
}

func getBanCommand(member *types.Member, serverID server.ServerID) string {
	command := conf.Config.Whitelist.BanCommand

	for _, rc := range conf.Config.Whitelist.RolesConfig {
		if utils.CheckRole(member, rc.RoleID) {
			if rc.BanCommand != "" {
				command = rc.BanCommand
			}

			for _, s := range rc.PerServer {
				if s.ServerID == serverID {
					if s.BanCommand != "" {
						command = s.BanCommand
					}
				}
			}
		}
	}

	return command
}

type WhitelistInfo struct {
	PlayerName database.Player
	PlayerID   string
	DiscordID  database.UserID
}

func formatCommand(command string, info WhitelistInfo) string {
	return command
}

func GetDefaultProvider() WhitelistProvider {
	return &DefaultProvider{}
}

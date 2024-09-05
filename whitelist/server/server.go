package server

import "github.com/Sharktheone/ScharschBot/database"

type ServerID string

type ServerProvider interface {
	Whitelist(player database.Player, id ServerID)
	UnWhitelist(player database.Player, id ServerID)
	Ban(player database.Player, reason string, id ServerID)
	UnBan(player database.Player, id ServerID)
	SendCommand(command string, id ServerID)

	GetServers() []ServerID
}

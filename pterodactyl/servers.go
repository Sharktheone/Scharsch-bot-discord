package pterodactyl

import (
	"errors"
	"github.com/Sharktheone/ScharschBot/database"
)

var (
	ServerNotFoundErr = errors.New("server not found")
)

func GetServer(serverID string) (*Server, error) {
	var (
		server *Server
	)
	for _, s := range Servers {
		if s.Config.ServerID == serverID {
			server = s
			return server, nil
		}
	}
	return nil, ServerNotFoundErr
}

func GetServerByName(serverName string) (*Server, error) {
	var (
		server *Server
	)
	for _, s := range Servers {
		if s.Config.ServerName == serverName {
			server = s
			return server, nil
		}
	}
	return nil, ServerNotFoundErr
}

func GetAllPlayers() []*database.Player {
	var (
		players []*database.Player
	)
	for _, server := range Servers {
		if server.Config.SRV.Events.PlayerJoinLeft {
			for _, player := range server.OnlinePlayers.Players {
				players = append(players, player)
			}
		}
	}
	return players
}

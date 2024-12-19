package pterodactyl

import (
	"errors"
	server2 "github.com/Sharktheone/ScharschBot/whitelist/server"
)

var (
	ServerNotFoundErr = errors.New("server not found")
)

func GetServer(serverID string) (*Server, error) {
	var (
		server *Server
	)
	for _, s := range Servers {
		if s.Config.ServerID == server2.ServerID(serverID) {
			server = s
			return server, nil
		}
	}
	return nil, ServerNotFoundErr
}

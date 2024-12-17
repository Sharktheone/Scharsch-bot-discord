package pterodactyl

import (
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/whitelist/server"
	"log"
)

type Provider struct {
	servers map[server.ServerID]*Server
}

func (p Provider) Whitelist(player database.Player, id server.ServerID) {
	//TODO implement me
	panic("implement me")
}

func (p Provider) UnWhitelist(player database.Player, id server.ServerID) {
	//TODO implement me
	panic("implement me")
}

func (p Provider) Ban(player database.Player, reason string, id server.ServerID) {
	//TODO implement me
	panic("implement me")
}

func (p Provider) UnBan(player database.Player, id server.ServerID) {
	//TODO implement me
	panic("implement me")
}

func (p Provider) SendCommand(command string, id server.ServerID) {
	s, ok := p.servers[id]
	if !ok {
		log.Printf("Server %v not found", id)
		return
	}

	if err := s.SendCommand(command); err != nil {
		log.Printf("Failed to send command to server %v: %v", id, err)
	}
}

func (p Provider) GetServers() []server.ServerID {
	ids := make([]server.ServerID, 0, len(p.servers))

	for _, s := range p.servers {
		ids = append(ids, s.Config.ServerID)
	}

	return ids
}

func GetProvider() server.ServerProvider {
	return &Provider{}
}

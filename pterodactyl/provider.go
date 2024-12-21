package pterodactyl

import (
	"context"
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/pterodactyl/listeners"
	"github.com/Sharktheone/ScharschBot/pterodactyl/types"
	"github.com/Sharktheone/ScharschBot/whitelist/server"
	"log"
	"sync"
)

type Provider struct {
	servers map[server.ServerID]*Server
}

func (p *Provider) Whitelist(player database.Player, id server.ServerID) {
	command := conf.Config.Pterodactyl.WhitelistAddCommand

	whitelistCommand := fmt.Sprintf(command, player)

	p.SendCommand(whitelistCommand, id)
}

func (p *Provider) UnWhitelist(player database.Player, id server.ServerID) {
	command := conf.Config.Pterodactyl.WhitelistRemoveCommand

	unWhitelistCommand := fmt.Sprintf(command, player)

	p.SendCommand(unWhitelistCommand, id)
}

func (p *Provider) Ban(player database.Player, reason string, id server.ServerID) {
	//TODO implement me
	panic("implement me")
}

func (p *Provider) UnBan(player database.Player, id server.ServerID) {
	//TODO implement me
	panic("implement me")
}

func (p *Provider) SendCommand(command string, id server.ServerID) {
	log.Printf("Sending command %v to server %v", command, id)
	s, ok := p.servers[id]
	if !ok {
		log.Printf("Server %v not found", id)
		return
	}

	if err := s.SendCommand(command); err != nil {
		log.Printf("Failed to send command to server %v: %v", id, err)
	}
}

func (p *Provider) GetServers() []server.ServerID {
	ids := make([]server.ServerID, 0, len(p.servers))

	for _, s := range p.servers {
		ids = append(ids, s.Config.ServerID)
	}

	return ids
}

func GetProvider() server.ServerProvider {
	servers := make(map[server.ServerID]*Server, len(conf.Config.Pterodactyl.Servers))

	mu := &sync.Mutex{}

	wg := &sync.WaitGroup{}

	for _, cnf := range conf.Config.Pterodactyl.Servers {
		wg.Add(1)
		go func(server conf.Server) {
			ctx := context.Background()
			s := New(&ctx, &server) //TODO: this probably should not be in the srv package

			if server.Console.Enabled {
				s.AddConsoleListener(func(server *conf.Server, console chan string) {
					listeners.ConsoleListener(ctx, server, console, nil)
				})
			}
			if server.StateMessages.Enabled {
				s.AddListener(func(ctx *context.Context, server *conf.Server, data chan *types.ChanData) {
					listeners.StatusListener(*ctx, server, data)
				}, string(server.ServerID+"_stateMessages"))
			}
			if server.ChannelInfo.Enabled {
				s.AddListener(func(ctx *context.Context, server *conf.Server, data chan *types.ChanData) {
					listeners.StatsListener(*ctx, server, data)
				}, string(server.ServerID+"_channelInfo"))
			}
			mu.Lock()
			servers[server.ServerID] = s
			mu.Unlock()
			wg.Done()

			if err := s.Listen(); err != nil {
				log.Printf("Error while listening to server %v: %v", server.ServerID, err)
			}

		}(cnf)
	}

	wg.Wait()

	return &Provider{servers: servers}
}

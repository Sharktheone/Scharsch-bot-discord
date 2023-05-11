package srv

import (
	"context"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/pterodactyl"
	"github.com/Sharktheone/ScharschBot/pterodactyl/listeners"
	"github.com/Sharktheone/ScharschBot/srv/api"
	"log"
)

var config = conf.GetConf()

func Start() {
	go api.Start()
	for _, server := range config.Pterodactyl.Servers {
		go func(server conf.Server) {
			ctx := context.Background()
			s := pterodactyl.New(&ctx, &server)
			if server.Console.Enabled {
				s.AddConsoleListener(func(server *conf.Server, console chan string) {
					listeners.ConsoleListener(ctx, server, console, nil)
				})
			}
			if server.StateMessages.Enabled {
				s.AddListener(func(ctx *context.Context, server *conf.Server, data chan *pterodactyl.ChanData) {
					listeners.StatusListener(*ctx, server, data)
				}, server.ServerID+"_stateMessages")
			}
			if server.ChannelInfo.Enabled {
				s.AddListener(func(ctx *context.Context, server *conf.Server, data chan *pterodactyl.ChanData) {
					listeners.StatsListener(*ctx, server, data)
				}, server.ServerID+"_channelInfo")
			}
			if err := s.Listen(); err != nil {
				log.Printf("Error while listening to server %v: %v", server.ServerID, err)
			}

		}(server)
	}
}

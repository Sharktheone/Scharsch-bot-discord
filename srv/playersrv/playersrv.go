package playersrv

import (
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/pterodactyl"
	"github.com/Sharktheone/ScharschBot/whitelist/whitelist"
	"log"
)

func CheckAccount(name database.Player) ([]database.Player, []database.Player) {
	owner := whitelist.GetOwner(name, nil)
	if conf.Config.Whitelist.KickUnWhitelisted {
		if !owner.Whitelisted {
			command := fmt.Sprintf(conf.Config.Whitelist.KickCommand, name)
			for _, listedServer := range conf.Config.Whitelist.Servers {
				for _, server := range conf.Config.Pterodactyl.Servers {
					if server.ServerName == listedServer {
						if err := pterodactyl.SendCommand(command, string(server.ServerID)); err != nil {
							log.Printf("Failed to send command to server %v: %v", server.ServerID, err)
						}
					}
				}
			}
		}
	}
	return owner.Players, owner.BannedPlayers
}

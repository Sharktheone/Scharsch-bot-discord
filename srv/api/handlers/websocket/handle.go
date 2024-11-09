package websocket

import (
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/discord/embed/srvEmbed"
	"github.com/Sharktheone/ScharschBot/types"
	"strings"
)

var (
	config = conf.Config
)

func (p *PSRVEvent) processEvent() {
	if p.h.authenticated == false && p.e.Event != Auth {
		return
	}
	switch p.e.Event {
	case SendPlayers:
		p.sendPlayers()
	case KickPlayer:
		p.kickPlayer()
	case ReportPlayer:
		p.reportPlayer()
	case BanPlayer:
		p.banPlayer()
	case UnbanPlayer:
		p.unbanPlayer()
	case PlayerJoined:
		p.playerJoined()
	case PlayerLeft:
		p.playerLeft()
	case Players:
		p.players()
	case ChatMessage:
		p.chatMessage()
	case PlayerDeath:
		p.playerDeath()
	case PlayerAdvancement:
		p.playerAdvancement()
	case Auth:
		p.auth()
	}
}

// sendPlayers send total online players to server
func (p *PSRVEvent) sendPlayers() {
	var players []database.Player
	for _, player := range p.h.server.OnlinePlayers.Players {
		players = append(players, database.Player(*player))
	}

	p.h.send <- &types.WebsocketEvent{
		Event: Players,
		Data: types.WebsocketEventData{
			Players: players,
		},
	}
}

// kickPlayer kick player on all servers
func (p *PSRVEvent) kickPlayer() {

}

// reportPlayer report player
func (p *PSRVEvent) reportPlayer() {

}

// banPlayer ban player on all servers
func (p *PSRVEvent) banPlayer() {

}

// unbanPlayer unban player on all servers
func (p *PSRVEvent) unbanPlayer() {

}

func (p *PSRVEvent) playerJoined() {
	p.h.server.OnlinePlayers.Players = append(p.h.server.OnlinePlayers.Players, &p.e.Data.Player)
	if p.h.server.Config.SRV.Events.PlayerJoinLeft {
		messageEmbed := srvEmbed.PlayerJoin(p.e, p.h.server.Config, p.footerIcon, p.username, p.session)
		p.session.SendEmbeds(p.h.server.Config.SRV.ChannelID, messageEmbed, "Join")
	}
}

func (p *PSRVEvent) playerLeft() {
	p.h.server.OnlinePlayers.Players = append(p.h.server.OnlinePlayers.Players, &p.e.Data.Player)
	if p.h.server.Config.SRV.Events.PlayerJoinLeft {
		messageEmbed := srvEmbed.PlayerQuit(p.e, p.h.server.Config, p.footerIcon, p.username, p.session)
		p.session.SendEmbeds(p.h.server.Config.SRV.ChannelID, messageEmbed, "Left")
	}
}

func (p *PSRVEvent) players() {
	if p.h.server.Config.SRV.Events.PlayerJoinLeft {
		p.h.server.OnlinePlayers.Mu.Lock()
		defer p.h.server.OnlinePlayers.Mu.Unlock()
		var players []*database.Player
		for _, player := range p.e.Data.Players {
			players = append(players, &player)
		}
		p.h.server.OnlinePlayers.Players = players
	}
}

func (p *PSRVEvent) chatMessage() {
	if p.h.server.Config.Chat.Enabled {
		if p.h.server.Config.Chat.Embed {
			messageEmbed := srvEmbed.Chat(p.e, p.h.server.Config, p.footerIcon, p.username, p.session)
			p.session.SendEmbeds(p.h.server.Config.SRV.ChannelID, messageEmbed, "Chat")
		} else {
			p.session.SendMessages(p.h.server.Config.SRV.ChannelID, fmt.Sprintf("%v%v %v", p.e.Data.Player, p.h.server.Config.Chat.Prefix, p.e.Data.Message), "Chat")
		}
	}
}

func (p *PSRVEvent) playerDeath() {
	if p.h.server.Config.SRV.Events.Death {
		messageEmbed := srvEmbed.PlayerDeath(p.e, p.h.server.Config, p.footerIcon, p.username, p.session)
		p.session.SendEmbeds(p.h.server.Config.SRV.ChannelID, messageEmbed, "Death")
	}

}

func (p *PSRVEvent) playerAdvancement() {
	if p.h.server.Config.SRV.Events.Advancement {
		if strings.Contains(p.e.Data.Advancement, "root") && !p.h.server.Config.SRV.Events.RootAdvancements {
			return
		}
		messageEmbed := srvEmbed.PlayerAdvancement(p.e, p.h.server.Config, p.footerIcon, p.username, p.session)
		p.session.SendEmbeds(p.h.server.Config.SRV.ChannelID, &messageEmbed, "Advancement")
	}
}

func (p *PSRVEvent) auth() {
	if p.e.Data.Password == config.SRV.API.Password && p.e.Data.User == config.SRV.API.User {
		p.h.authenticated = true
		p.h.send <- &types.WebsocketEvent{
			Event: AuthSuccess,
		}
	} else {
		p.h.authenticated = false
		p.h.send <- &types.WebsocketEvent{
			Event: AuthFailed,
		}
	}
}

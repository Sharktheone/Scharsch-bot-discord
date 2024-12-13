package pEmbed

import (
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/pterodactyl"
	"github.com/bwmarrin/discordgo"
)

func Power(action string) *discordgo.MessageEmbed {
	var (
		color  int
		Fields = getServerFields()
	)
	switch action {
	case pterodactyl.PowerSignalStart:
		color = 0x00FF00
	case pterodactyl.PowerSignalStop:
		color = 0xFF0000
	case pterodactyl.PowerSignalRestart:
		color = 0xFFFF00
	case "status":
		color = 0x00AAFF
	default:
		color = 0x000000
	}
	return &discordgo.MessageEmbed{
		Title:  fmt.Sprintf("Select a server to %v", action),
		Color:  color,
		Fields: Fields,
	}
}

func getServerFields() []*discordgo.MessageEmbedField {
	var (
		Fields []*discordgo.MessageEmbedField
	)
	for _, server := range pterodactyl.Servers {
		var StateMsg string
		switch server.Status.State {
		case pterodactyl.PowerStatusStarting:
			StateMsg = conf.Config.SRV.States.Starting
		case pterodactyl.PowerStatusStopping:
			StateMsg = conf.Config.SRV.States.Stopping
		case pterodactyl.PowerStatusRunning:
			StateMsg = conf.Config.SRV.States.Online
		case pterodactyl.PowerStatusOffline:
			StateMsg = conf.Config.SRV.States.Offline
		}
		Fields = append(Fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("%v:", server.Config.ServerName),
			Value: StateMsg,
		})
	}
	return Fields
}
func PowerNotAllowed(avatarURL string, name string, action string, serverName string) discordgo.MessageEmbed {
	var (
		Title string
	)
	if serverName != "" {
		Title = fmt.Sprintf("You have no permission to %v %v", action, serverName)
	} else {
		Title = fmt.Sprintf("You have no permission to %v servers", action)
	}

	Embed := discordgo.MessageEmbed{
		Title: Title,
		Color: 0xFF0000,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    name,
			IconURL: avatarURL,
			URL:     conf.Config.Discord.EmbedErrorAuthorURL,
		},
	}
	return Embed
}
func PowerAction(action string, serverName string, avatarURL string, name string) discordgo.MessageEmbed {
	var (
		Title string
		color int
	)
	Title = fmt.Sprintf("Server %v is getting %ved", serverName, action)
	switch action {
	case "start":
		color = 0x00FF00
	case "stop":
		color = 0xFF0000
	case "restart":
		color = 0xFFFF00
	}

	Embed := discordgo.MessageEmbed{
		Title: Title,
		Color: color,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    name,
			IconURL: avatarURL,
		},
	}
	return Embed
}

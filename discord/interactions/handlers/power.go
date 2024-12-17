package handlers

import (
	"fmt"
	"github.com/Sharktheone/ScharschBot/discord/embed/pEmbed"
	"github.com/Sharktheone/ScharschBot/discord/session"
	"github.com/Sharktheone/ScharschBot/pterodactyl"
	"github.com/Sharktheone/ScharschBot/pterodactyl/types"
	"github.com/bwmarrin/discordgo"
	"log"
)

func PowerMain(s *session.Session, i *discordgo.InteractionCreate) {
	power(s, i, i.ApplicationCommandData().Options[0].Name)
}
func PowerStart(s *session.Session, i *discordgo.InteractionCreate) {
	power(s, i, types.PowerSignalStart)
}
func PowerRestart(s *session.Session, i *discordgo.InteractionCreate) {
	power(s, i, types.PowerSignalRestart)
}
func PowerStop(s *session.Session, i *discordgo.InteractionCreate) {
	power(s, i, types.PowerSignalStop)
}

func power(s *session.Session, i *discordgo.InteractionCreate, option string) {
	var (
		restartDisabled = false
		stopDisabled    = false
		startDisabled   = false
	)

	serverSelect := discordgo.SelectMenu{
		Placeholder: fmt.Sprintf("Select a server to %s", option),
		CustomID:    fmt.Sprintf("power_%s_select", option),
		Options:     getServerOptions(option),
	}

	switch option {
	case types.PowerSignalStart:
		startDisabled = true
	case types.PowerSignalStop:
		stopDisabled = true
	case types.PowerSignalRestart:
		restartDisabled = true
	}

	buttonRow := []discordgo.MessageComponent{
		&discordgo.Button{
			Label:    "Start",
			Style:    discordgo.SuccessButton,
			CustomID: "power_start",
			Disabled: startDisabled,
		},
		&discordgo.Button{
			Label:    "Restart",
			Style:    discordgo.PrimaryButton,
			CustomID: "power_restart",
			Disabled: restartDisabled,
		},
		&discordgo.Button{
			Label:    "Stop",
			Style:    discordgo.DangerButton,
			CustomID: "power_stop",
			Disabled: stopDisabled,
		},
	}

	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				pEmbed.Power(option),
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: buttonRow,
				},
				func(option string) discordgo.ActionsRow {
					if option != "status" {
						return discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								&serverSelect,
							},
						}
					}
					return discordgo.ActionsRow{}
				}(option),
			},
		},
	}); err != nil {
		log.Printf("Failed to respond to command power: %v", err)
	}
}

func getServerOptions(option string) []discordgo.SelectMenuOption {
	var fields []discordgo.SelectMenuOption
	for _, server := range pterodactyl.Servers {
		if option == types.PowerSignalStart && server.Status.State == types.PowerStatusRunning ||
			option == types.PowerSignalStop && server.Status.State == types.PowerStatusOffline {
			continue
		}
		fields = append(fields, discordgo.SelectMenuOption{
			Label: server.Config.ServerName,
			Value: string(server.Config.ServerID),
		})
	}
	return fields
}

func powerSelect(s *session.Session, i *discordgo.InteractionCreate, action string) {
	var (
		allowed         = false
		options         = i.MessageComponentData()
		requiredRole    []string
		serverConf, err = pterodactyl.GetServer(options.Values[0])
	)
	if err != nil {
		log.Printf("Failed to get server %s: %v", options.Values[0], err)
		return
	}
	switch action {
	case "start":
		requiredRole = serverConf.Config.PowerActionsRoleIDs.Start
	case "stop":
		requiredRole = serverConf.Config.PowerActionsRoleIDs.Stop
	case "restart":
		requiredRole = serverConf.Config.PowerActionsRoleIDs.Restart

	}
	for _, role := range i.Member.Roles {
		for _, required := range requiredRole {
			if required == role {
				allowed = true
				break
			}
		}
	}
	var messageEmbed discordgo.MessageEmbed
	if !allowed {
		messageEmbed = pEmbed.PowerNotAllowed(i.Member.User.AvatarURL("40"), i.Member.User.String(), action, serverConf.Config.ServerName)
	} else {
		messageEmbed = pEmbed.PowerAction(action, serverConf.Config.ServerName, i.Member.User.AvatarURL("40"), i.Member.User.Username)
		s, err := pterodactyl.GetServer(string(serverConf.Config.ServerID))
		if err != nil {
			log.Printf("Failed to get server %s: %v", serverConf.Config.ServerName, err)
			return
		}
		switch action {
		case "start":
			if err := s.Start(); err != nil {
				log.Printf("Failed to start server %s: %v", serverConf.Config.ServerName, err)
			}
		case "stop":
			if err := s.Stop(); err != nil {
				log.Printf("Failed to stop server %s: %v", serverConf.Config.ServerName, err)
			}
		case "restart":
			if err := s.Restart(); err != nil {
				log.Printf("Failed to restart server %s: %v", serverConf.Config.ServerName, err)
			}

		}
	}
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				&messageEmbed,
			},
		},
	}); err != nil {
		log.Printf("Failed send power_%v Embed: %v", action, err)
	}
}

func PowerStartSelect(s *session.Session, i *discordgo.InteractionCreate) {
	powerSelect(s, i, "start")
}
func PowerRestartSelect(s *session.Session, i *discordgo.InteractionCreate) {
	powerSelect(s, i, "restart")

}
func PowerStopSelect(s *session.Session, i *discordgo.InteractionCreate) {
	powerSelect(s, i, "stop")

}

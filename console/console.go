package console

import (
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/discord/session"
	"github.com/Sharktheone/ScharschBot/pterodactyl"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

var config = conf.Config

func Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, server := range conf.Config.Pterodactyl.Servers {
		if server.Console.Reverse {
			for _, neededChannelID := range server.Console.ChannelID {
				if neededChannelID == m.ChannelID {
					command := strings.SplitAfter(m.Message.Content, server.Console.ReversePrefix)
					if command[0] == server.Console.ReversePrefix {
						if session.HasRole(m.Member, config.Whitelist.Roles.ServerRoleID) {
							var commandString string
							for _, element := range command[1:] {
								commandString += element
							}
							log.Printf("%v is sending command to server %v: %v", m.Author.ID, server.ServerName, commandString)
							if err := pterodactyl.SendCommand(commandString, server.ServerID); err != nil {
								log.Printf("Failed to send command to server %v: %v", server.ServerID, err)
							}
						}
					}
					break
				}
			}
		}

	}
}

func ChatHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, server := range config.Pterodactyl.Servers {
		if server.Chat.Reverse {
			for _, neededChannelID := range server.Chat.ChannelID {
				if neededChannelID == m.ChannelID && m.Author.ID != s.State.User.ID {
					if session.HasRole(m.Member, config.Whitelist.Roles.ServerRoleID) {
						message := fmt.Sprintf(" %v: %v", m.Author.Username, m.Message.Content)
						command := fmt.Sprintf(config.Pterodactyl.ChatCommand, message)
						if err := pterodactyl.SendCommand(command, server.ServerID); err != nil {
							log.Printf("Failed to send chat message to server %v: %v", server.ServerID, err)
						}
					}
				}

			}
		}
	}
}

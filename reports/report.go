package reports

import (
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/discord/session"
	"github.com/bwmarrin/discordgo"
	"log"
)

var (
	config = conf.Config
)

type ReportData struct {
	ReporterID     string `bson:"reporterID"`
	ReportedPlayer string `bson:"reportedPlayer"`
	Reason         string `bson:"reason"`
}

func Report(name database.Player, reason string, i *discordgo.InteractionCreate, s *session.Session, messageEmbed discordgo.MessageEmbed) (reportAllowed bool, alreadyReported bool, enabled bool) {
	var (
		allowed   = false
		dataFound bool
	)
	if config.Whitelist.Report.Enabled {
		for _, role := range i.Member.Roles {
			for _, requiredRole := range config.Whitelist.Report.Roles {
				if role == requiredRole {
					allowed = true
					break
				}
			}
		}
		if allowed {
			dataFound = database.DB.IsAlreadyReported(database.Player(name)) //TODO: Probably we should still accept the report and allow persons to be reported multiple times

			if !dataFound {
				database.DB.Report(database.UserID(i.Member.User.ID), database.Player(name), reason)

				var roleMessage string
				for _, role := range config.Whitelist.Report.PingRoleID {
					ping := fmt.Sprintf("<@&%s> ", role)
					roleMessage += ping
				}
				for _, channel := range config.Whitelist.Report.ChannelID {
					_, err := s.ChannelMessageSendComplex(channel, &discordgo.MessageSend{
						Content: roleMessage,
						Embed:   &messageEmbed,
					})
					if err != nil {
						log.Printf("Failed to send report embed: %v", err)
					}
				}
			}
		}
	}
	return allowed, dataFound, config.Whitelist.Report.Enabled
}

func Reject(name string, i *discordgo.InteractionCreate, s *session.Session, notifyReporter bool, messageEmbed *discordgo.MessageEmbed, messageEmbedDMFailed *discordgo.MessageEmbed) (rejectAllowed bool, enabled bool) {
	var (
		allowed  = false
		notifyDM = config.Whitelist.Report.PlayerNotifyDM
	)
	if config.Whitelist.Report.Enabled {
		for _, role := range i.Member.Roles {
			for _, requiredRole := range config.Whitelist.Report.Roles {
				if role == requiredRole {
					allowed = true
					break
				}
			}
		}
	}

	if allowed {
		report, reportFound := database.DB.GetReportedPlayer(database.Player(name))
		if reportFound {
			if notifyDM {
				if notifyReporter {
					channel, err := s.UserChannelCreate(string(report.ReporterID))
					if err != nil {
						log.Printf("Failed to create DM with reporter: %v", err)

					}
					_, err = s.ChannelMessageSendEmbed(channel.ID, messageEmbed)
					if err != nil {
						log.Printf("Failed to send DM for reporter: %v, sending Message in normal Channels", err)
						for _, channelID := range config.Whitelist.Report.ChannelID {
							_, err = s.ChannelMessageSendEmbed(channelID, messageEmbedDMFailed)
							if err != nil {
								log.Printf("Failed to send Report message in normal Channel: %v", err)
							}
						}
					}

				}
			} else {
				if notifyReporter {
					for _, channelID := range config.Whitelist.Report.ChannelID {
						_, err := s.ChannelMessageSendEmbed(channelID, messageEmbed)
						if err != nil {
							log.Printf("Failed to send Report message : %v", err)
						}
					}

				}
			}
			DeleteReport(name)
		}
	}

	return allowed, config.Whitelist.Report.Enabled
}
func Accept(name string, i *discordgo.InteractionCreate, s *session.Session, notifyreporter bool, messageEmbed *discordgo.MessageEmbed, messageEmbedDMFailed *discordgo.MessageEmbed) (acceptAllowed bool, enabled bool) {
	var (
		allowed  = false
		notifyDM = config.Whitelist.Report.PlayerNotifyDM
	)
	if config.Whitelist.Report.Enabled {
		for _, role := range i.Member.Roles {
			for _, requiredRole := range config.Whitelist.Report.Roles {
				if role == requiredRole {
					allowed = true
					break
				}
			}
		}
	}
	if allowed {
		report, reportFound := database.DB.GetReportedPlayer(database.Player(name))
		if reportFound {
			if notifyDM {
				if notifyreporter {
					if err := s.SendDM(report.ReporterID, &discordgo.MessageSend{
						Embed: messageEmbed,
					},
						&discordgo.MessageSend{
							Content: fmt.Sprintf("<@%v>", report.ReporterID),
							Embed:   messageEmbedDMFailed,
						}); err != nil {
						log.Printf("Failed to send DM for reporter: %v, sending Message in normal Channels", err)
					}

				}
			} else {
				if notifyreporter {
					for _, channelID := range config.Whitelist.Report.ChannelID {
						_, err := s.ChannelMessageSendEmbed(channelID, messageEmbed)
						if err != nil {
							log.Printf("Failed to send Report message : %v", err)
						}
					}

				}
			}
			DeleteReport(name)
		}
	}
	return allowed, config.Whitelist.Report.Enabled
}

func DeleteReport(name string) {
	database.DB.DeleteReport(database.Player(name))
}

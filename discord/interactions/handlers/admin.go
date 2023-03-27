package handlers

import (
	"Scharsch-Bot/conf"
	"Scharsch-Bot/database/mongodb"
	"Scharsch-Bot/discord/embed/wEmbed"
	"Scharsch-Bot/discord/session"
	"Scharsch-Bot/reports"
	"Scharsch-Bot/whitelist/whitelist"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

var (
	config = conf.GetConf()
)

func Admin(s *session.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options[0].Options {
		optionMap[opt.Name] = opt
	}
	switch options[0].Name {
	case "whois":
		name := strings.ToLower(optionMap["name"].StringValue())
		var messageEmbed discordgo.MessageEmbed
		if mongodb.Ready {

			playerID, allowed, found := whitelist.Whois(name, i.Member.User.ID, i.Member.Roles)
			if allowed {
				if found {
					messageEmbed = wEmbed.WhitelistIsListedBy(name, playerID, i, s)
				} else {
					messageEmbed = wEmbed.WhitelistNotListed(name, i)
				}
			} else {
				messageEmbed = wEmbed.WhitelistWhoisNotAllowed(name, i)
			}
		} else {
			messageEmbed = wEmbed.DatabaseNotReady
		}
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					&messageEmbed,
				},
			},
		})
		if err != nil {
			log.Printf("Failed execute command whitelist: %v", err)
		}
	case "user":
		user := optionMap["user"].UserValue(s.Session)
		playerID := user.ID
		var messageEmbed discordgo.MessageEmbed
		if mongodb.Ready {
			accounts, allowed, found, bannedPlayers := whitelist.HasListed(playerID, i.Member.User.ID, i.Member.Roles)
			if allowed {
				if found || len(bannedPlayers) > 0 {
					messageEmbed = wEmbed.WhitelistHasListed(accounts, playerID, bannedPlayers, i, s)
				} else {
					messageEmbed = wEmbed.WhitelistNoAccounts(i, playerID)
				}
			} else {
				messageEmbed = wEmbed.WhitelistUserNotAllowed(accounts, playerID, bannedPlayers, i)
			}
		} else {
			messageEmbed = wEmbed.DatabaseNotReady
		}
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					&messageEmbed,
				},
			},
		})
		if err != nil {
			log.Printf("Failed execute command whitelistuser: %v", err)
		}
	case "banuser":
		user := optionMap["user"].UserValue(s.Session)
		var reason = "No reason provided"
		if optionMap["reason"] != nil {
			reason = optionMap["reason"].StringValue()
		}
		playerID := user.ID
		banAccounts := true
		if optionMap["removeaccounts"] != nil {
			banAccounts = optionMap["removeaccounts"].BoolValue()
		}
		var messageEmbed discordgo.MessageEmbed
		if mongodb.Ready {
			allowed, alreadyBanned := whitelist.BanUserID(i.Member.User.ID, i.Member.Roles, playerID, banAccounts, reason, s)
			if allowed {
				if alreadyBanned {
					messageEmbed = wEmbed.AlreadyBanned(user.Username)
				} else {
					messageEmbed = wEmbed.WhitelistBanUserID(playerID, reason, i, s)
				}
			} else {
				messageEmbed = wEmbed.WhitelistBanUserIDNotAllowed(playerID, i)
			}
		} else {
			messageEmbed = wEmbed.DatabaseNotReady
		}
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					&messageEmbed,
				},
			},
		})
		if err != nil {
			log.Printf("Failed execute command whitelistbanuserid: %v", err)
		}
	case "banaccount":
		name := strings.ToLower(optionMap["name"].StringValue())
		var reason = "No reason provided"
		if optionMap["reason"] != nil {
			reason = optionMap["reason"].StringValue()
		}
		var messageEmbed discordgo.MessageEmbed
		if mongodb.Ready {

			allowed, playerID := whitelist.BanAccount(i.Member.User.ID, i.Member.Roles, name, reason, s)
			if allowed {
				messageEmbed = wEmbed.WhitelistBanAccount(name, playerID, reason, i, s)
			} else {
				messageEmbed = wEmbed.WhitelistBanAccountNotAllowed(name, i)
			}
		} else {
			messageEmbed = wEmbed.DatabaseNotReady
		}
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					&messageEmbed,
				},
			},
		})
		if err != nil {
			log.Printf("Failed execute command whitelistbanaccount: %v", err)
		}
	case "unbanuser":
		user := optionMap["user"].UserValue(s.Session)
		playerID := user.ID
		unbanAccounts := false
		if optionMap["unbanaccounts"] != nil {
			unbanAccounts = optionMap["unbanaccounts"].BoolValue()
		}
		var messageEmbed discordgo.MessageEmbed
		if mongodb.Ready {
			allowed := whitelist.UnBanUserID(i.Member.User.ID, i.Member.Roles, playerID, unbanAccounts, s)
			if allowed {
				messageEmbed = wEmbed.WhitelistUnBanUserID(playerID, i, s)
			} else {
				messageEmbed = wEmbed.WhitelistBanUserIDNotAllowed(playerID, i)
			}
		} else {
			messageEmbed = wEmbed.DatabaseNotReady
		}
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					&messageEmbed,
				},
			},
		})
		if err != nil {
			log.Printf("Failed execute command whitelistunbanuserid: %v", err)
		}
	case "unbanaccount":
		name := strings.ToLower(optionMap["name"].StringValue())
		var messageEmbed discordgo.MessageEmbed
		if mongodb.Ready {
			allowed := whitelist.UnBanAccount(i.Member.User.ID, i.Member.Roles, name, s)
			if allowed {
				messageEmbed = wEmbed.WhitelistUnBanAccount(name, i, s)
			} else {
				messageEmbed = wEmbed.WhitelistBanAccountNotAllowed(name, i)
			}
		} else {
			messageEmbed = wEmbed.DatabaseNotReady
		}
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					&messageEmbed,
				},
			},
		})
		if err != nil {
			log.Printf("Failed execute command whitelistunbanaccount: %v", err)
		}
	case "removeall":
		var (
			messageEmbed discordgo.MessageEmbed
			err          error
		)
		if mongodb.Ready {
			allowed := whitelist.RemoveAllAllowed(i.Member.Roles)
			if allowed {
				var button discordgo.Button
				messageEmbed, button = wEmbed.WhitelistRemoveAllSure(i)
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							&messageEmbed,
						},
						Components: []discordgo.MessageComponent{
							discordgo.ActionsRow{
								Components: []discordgo.MessageComponent{
									&button,
								},
							},
						},
					},
				})

			} else {
				messageEmbed = wEmbed.WhitelistRemoveAllNotAllowed(i)
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							&messageEmbed,
						},
					},
				})
			}
		} else {
			messageEmbed = wEmbed.DatabaseNotReady
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						&messageEmbed,
					},
				},
			})
		}

		if err != nil {
			log.Printf("Failed execute command whitelistremoveall: %v", err)
		}
	case "listreports":
		var (
			messageEmbed discordgo.MessageEmbed
			allowed      bool
			enabled      = config.Whitelist.Report.Enabled
		)
		if mongodb.Ready {
			if config.Whitelist.Report.Enabled {
				for _, role := range i.Member.Roles {
					for _, requiredRole := range config.Discord.WhitelistBanRoleID { // TODO: Add Report Admin Role
						if role == requiredRole {
							allowed = true
							break
						}
					}
				}
				if allowed {
					if enabled {
						messageEmbed = wEmbed.ListReports(i)
					} else {
						messageEmbed = wEmbed.ReportDisabled(i)
					}
				} else {
					messageEmbed = wEmbed.ReportNotALlowed(i)
				}
			}
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						&messageEmbed,
					},
				},
			})
			if err != nil {
				log.Printf("Failed execute command whitelistlistreports: %v", err)
			}
		}
	case "rejectreport":
		var (
			messageEmbed   discordgo.MessageEmbed
			name           string
			notifyreporter = true
		)
		if optionMap["name"] != nil {
			name = strings.ToLower(optionMap["name"].StringValue())
		}
		if optionMap["notifyreporter"] != nil {
			notifyreporter = optionMap["notifyreporter"].BoolValue()
		}

		if mongodb.Ready {
			report, _ := reports.GetReport(name)
			reportMessageEmbed := wEmbed.ReportUserAction(name, false, report.ReporterID, s, "rejected")
			reportMessageEmbedDMFailed := wEmbed.ReportUserAction(name, true, report.ReporterID, s, "rejected")

			allowed, enabled := reports.Reject(name, i, s, notifyreporter, &reportMessageEmbed, &reportMessageEmbedDMFailed)
			if allowed {
				if enabled {
					messageEmbed = wEmbed.ReportAction(name, "rejected", notifyreporter)
				} else {
					messageEmbed = wEmbed.ReportDisabled(i)
				}
			} else {
				messageEmbed = wEmbed.ReportNotALlowed(i)
			}
		}
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					&messageEmbed,
				},
			},
		})
		if err != nil {
			log.Printf("Failed execute command whitelistrejectreport: %v", err)
		}
	case "acceptreport":
		var (
			messageEmbed   discordgo.MessageEmbed
			name           string
			notifyreporter = true
		)
		if optionMap["name"] != nil {
			name = strings.ToLower(optionMap["name"].StringValue())
		}
		if optionMap["notifyreporter"] != nil {
			notifyreporter = optionMap["notifyreporter"].BoolValue()
		}

		if mongodb.Ready {
			report, _ := reports.GetReport(name)
			reportMessageEmbed := wEmbed.ReportUserAction(name, false, report.ReporterID, s, "accepted")
			reportMessageEmbedDMFailed := wEmbed.ReportUserAction(name, true, report.ReporterID, s, "accepted")

			allowed, enabled := reports.Accept(name, i, s, notifyreporter, &reportMessageEmbed, &reportMessageEmbedDMFailed)
			if allowed {
				if enabled {
					messageEmbed = wEmbed.ReportAction(name, "accepted", notifyreporter)
				} else {
					messageEmbed = wEmbed.ReportDisabled(i)
				}
			} else {
				messageEmbed = wEmbed.ReportNotALlowed(i)
			}
		}
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					&messageEmbed,
				},
			},
		})
		if err != nil {
			log.Printf("Failed execute command whitelistrejectreport: %v", err)
		}
	}
}

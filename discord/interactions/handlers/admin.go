package handlers

import (
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/discord/embed/wEmbed"
	"github.com/Sharktheone/ScharschBot/discord/session"
	"github.com/Sharktheone/ScharschBot/reports"
	"github.com/Sharktheone/ScharschBot/whitelist/whitelist"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

var (
	config = conf.Config
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

		accounts, allowed, found, bannedPlayers := whitelist.HasListed(playerID, i.Member.User.ID, i.Member.Roles, false)
		if allowed {
			if found || len(bannedPlayers) > 0 {
				messageEmbed = wEmbed.WhitelistHasListed(accounts, playerID, bannedPlayers, i, s)
			} else {
				messageEmbed = wEmbed.WhitelistNoAccounts(i, playerID)
			}
		} else {
			messageEmbed = wEmbed.WhitelistUserNotAllowed(accounts, playerID, bannedPlayers, i)
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

		allowed, owner := whitelist.BanAccount(i.Member.User.ID, i.Member.Roles, name, reason, s)
		if allowed {
			messageEmbed = wEmbed.WhitelistBanAccount(name, owner.ID, reason, i, s)
		} else {
			messageEmbed = wEmbed.WhitelistBanAccountNotAllowed(name, i)
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

		allowed := whitelist.UnBanUserID(i.Member.User.ID, i.Member.Roles, playerID, unbanAccounts, s)
		if allowed {
			messageEmbed = wEmbed.WhitelistUnBanUserID(playerID, i, s)
		} else {
			messageEmbed = wEmbed.WhitelistBanUserIDNotAllowed(playerID, i)
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
		allowed := whitelist.UnBanAccount(i.Member.User.ID, i.Member.Roles, name, s)
		if allowed {
			messageEmbed = wEmbed.WhitelistUnBanAccount(name, i, s)
		} else {
			messageEmbed = wEmbed.WhitelistBanAccountNotAllowed(name, i)
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

		if err != nil {
			log.Printf("Failed execute command whitelistremoveall: %v", err)
		}
	case "listreports":
		var (
			messageEmbed discordgo.MessageEmbed
			allowed      bool
			enabled      = config.Whitelist.Report.Enabled
		)
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

		report, _ := database.DB.GetReportedPlayer(database.Player(name))
		reportMessageEmbed := wEmbed.ReportUserAction(name, false, string(report.ReporterID), s, "rejected")
		reportMessageEmbedDMFailed := wEmbed.ReportUserAction(name, true, string(report.ReporterID), s, "rejected")

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

		report, _ := database.DB.GetReportedPlayer(database.Player(name))
		reportMessageEmbed := wEmbed.ReportUserAction(name, false, string(report.ReporterID), s, "accepted")
		reportMessageEmbedDMFailed := wEmbed.ReportUserAction(name, true, string(report.ReporterID), s, "accepted")

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

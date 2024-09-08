package handlers

import (
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/discord/embed/wEmbed"
	"github.com/Sharktheone/ScharschBot/discord/session"
	"github.com/Sharktheone/ScharschBot/reports"
	"github.com/Sharktheone/ScharschBot/types"
	"github.com/Sharktheone/ScharschBot/whitelist/whitelist"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func Whitelist(s *session.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options[0].Options {
		optionMap[opt.Name] = opt
	}
	switch options[0].Name {
	case "add":
		name := strings.ToLower(optionMap["name"].StringValue())
		var (
			messageEmbed discordgo.MessageEmbed
		)

		member := types.MemberFromDG(i.Member)

		result, banReason := whitelist.Add(database.Player(name), member)

		switch result {
		case whitelist.NoFreeAccount:
			listedAccounts := whitelist.ListedAccountsOf(member.ID, false)
			var (
				removeOptions []discordgo.SelectMenuOption
			)
			for _, acc := range listedAccounts {
				removeOptions = append(removeOptions, discordgo.SelectMenuOption{
					Label: string(acc),
					Value: string(acc),
				})
			}
			removeSelect := discordgo.SelectMenu{
				Placeholder: "Remove accounts",
				CustomID:    "remove_select",
				Options:     removeOptions,
			}
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						&messageEmbed,
					},
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								&removeSelect,
							},
						},
					},
				}})

			if err != nil {
				log.Printf("Failed execute command whitelistadd: %v", err)
			}

			return

		case whitelist.AlreadyListed:
			messageEmbed = wEmbed.WhitelistAlreadyListed(name, i)
		case whitelist.NotExisting:
			messageEmbed = wEmbed.WhitelistNotExisting(name, i)
		case whitelist.NotAllowed:
			messageEmbed = wEmbed.WhitelistAddNotAllowed(name, i)
		case whitelist.McBanned:
			messageEmbed = wEmbed.WhitelistBanned(name, false, banReason, i)
		case whitelist.DcBanned:
			messageEmbed = wEmbed.WhitelistBanned(name, true, banReason, i)
		case whitelist.BothBanned:
			messageEmbed = wEmbed.WhitelistBanned(name, true, banReason, i)
		case whitelist.Ok:
			messageEmbed = wEmbed.WhitelistAdding(name, i)
		}

		var err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					&messageEmbed,
				},
			},
		})
		if err != nil {
			log.Printf("Failed execute command whitelistadd: %v", err)
		}
	case "remove":
		name := strings.ToLower(optionMap["name"].StringValue())
		var messageEmbed discordgo.MessageEmbed

		allowed, onWhitelist := whitelist.Remove(database.Player(name), types.MemberFromDG(i.Member))

		if allowed {
			if onWhitelist {
				messageEmbed = wEmbed.WhitelistRemoving(name, i)
			} else {
				messageEmbed = wEmbed.WhitelistNotListed(name, i)
			}
		} else {
			messageEmbed = wEmbed.WhitelistRemoveNotAllowed(name, i)
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
			log.Printf("Failed execute command whitelistremove: %v", err)
		}
	case "myaccounts":
		var messageEmbed discordgo.MessageEmbed

		accounts, allowed, found, bannedPlayers := whitelist.HasListed(i.Member.User.ID, i.Member.User.ID, i.Member.Roles, true)
		if allowed {
			if found || len(bannedPlayers) > 0 {
				messageEmbed = wEmbed.WhitelistHasListed(accounts, i.Member.User.ID, bannedPlayers, i, s)
			} else {
				messageEmbed = wEmbed.WhitelistNoAccounts(i, i.Member.User.ID)
			}
		} else {
			messageEmbed = wEmbed.WhitelistUserNotAllowed(accounts, i.Member.User.ID, bannedPlayers, i)
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
			log.Printf("Failed execute command whitelistmyaccounts: %v", err)
		}
	case "removemyaccounts":
		var messageEmbed discordgo.MessageEmbed

		hasListedAccounts, listedAccounts := whitelist.RemoveMyAccounts(i.Member.User.ID)
		mcBans := whitelist.CheckBans(i.Member.User.ID)

		if hasListedAccounts {
			messageEmbed = wEmbed.WhitelistRemoveMyAccounts(listedAccounts, mcBans, i)
		} else {
			messageEmbed = wEmbed.WhitelistNoAccounts(i, i.Member.User.ID)
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
			log.Printf("Failed execute command whitelistremove: %v", err)
		}
	case "report":
		var (
			messageEmbed discordgo.MessageEmbed
			name         string
			reason       = "No reason provided"
		)
		if optionMap["name"] != nil {
			name = strings.ToLower(optionMap["name"].StringValue())
		}
		if optionMap["reason"] != nil {
			reason = optionMap["reason"].StringValue()
		}

		reportEmbed := wEmbed.NewReport(name, reason, i)
		allowed, alreadyReported, enabled := reports.Report(name, reason, i, s, reportEmbed)
		if allowed {
			if enabled {
				if alreadyReported {
					messageEmbed = wEmbed.AlreadyReported(name)
				} else {
					messageEmbed = wEmbed.ReportPlayer(name, reason, i)
				}
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
			log.Printf("Failed execute command report %v", err)
		}
		log.Printf("%v reported %v for %v", i.Member.User.Username, name, reason)

	}

}

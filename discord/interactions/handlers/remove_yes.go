package handlers

import (
	"github.com/Sharktheone/ScharschBot/discord/embed/wEmbed"
	"github.com/Sharktheone/ScharschBot/discord/session"
	"github.com/Sharktheone/ScharschBot/types"
	"github.com/Sharktheone/ScharschBot/whitelist/whitelist"
	"github.com/bwmarrin/discordgo"
	"log"
)

func RemoveYes(s *session.Session, i *discordgo.InteractionCreate) {
	var messageEmbed discordgo.MessageEmbed

	member := types.MemberFromDG(i.Member)

	allowed, onWhitelist := whitelist.RemoveAll(member)
	if allowed {
		if onWhitelist {
			messageEmbed = wEmbed.WhitelistRemoveAll(i)
		} else {
			messageEmbed = wEmbed.WhitelistRemoveAllNoWhitelistEntries(i)
		}
	} else {
		messageEmbed = wEmbed.WhitelistRemoveAllNotAllowed(i)
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
		log.Printf("Failed execute component remove_yes: %v", err)
	}
}

package session

import (
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/bwmarrin/discordgo"
	"log"
)

var (
	config = conf.Config
)

type Session struct {
	*discordgo.Session
}

func (s *Session) GetUserProfile(userID string) (*discordgo.Member, error) {
	if user, err := s.GuildMember(config.Discord.ServerID, userID); err != nil {
		return &discordgo.Member{}, fmt.Errorf("failed to get user profile: %v", err)
	} else {
		return user, nil
	}
}

func (s *Session) GetRoles(userID string) ([]string, error) {
	if user, err := s.GuildMember(config.Discord.ServerID, userID); err != nil {
		return nil, err
	} else {
		return user.Roles, nil
	}
}

func (s *Session) SendDM(userID string, messageComplexDM *discordgo.MessageSend, messageComplexDMFailed *discordgo.MessageSend) error {
	channel, err := s.UserChannelCreate(userID)
	if err != nil {
		log.Printf("Failed to create DM with reporter: %v", err)

	}
	_, err = s.ChannelMessageSendComplex(channel.ID, messageComplexDM)
	if err != nil {
		log.Printf("Failed to send DM: %v, sending Message in normal Channels", err)
		for _, channelID := range config.Whitelist.Report.ChannelID {
			_, err = s.ChannelMessageSendComplex(channelID, messageComplexDMFailed)
			if err != nil {
				return fmt.Errorf("failed to send message in dm alternative channel on server: %v", err)
			}
		}
	}
	return nil
}

func (s *Session) SendEmbeds(channelID []string, embed *discordgo.MessageEmbed, embedType string) {
	for _, channel := range channelID {
		if _, err := s.ChannelMessageSendEmbed(channel, embed); err != nil {
			log.Printf("Failed to send %v embed: %v (channelID: %v)", embedType, err, channel)
		}
	}
}

func (s *Session) SendMessages(channelID []string, message string, messageType string) {
	for _, channel := range channelID {
		if _, err := s.ChannelMessageSend(channel, message); err != nil {
			log.Printf("Failed to send %v message: %v (channelID: %v)", messageType, err, channel)
		}
	}
}

func HasRole(member *discordgo.Member, roleIDs []database.Role) bool {
	return HasRoleID(member.Roles, roleIDs)
}

func HasRoleID(hasRoles []string, neededRoles []database.Role) bool {
	for _, role := range hasRoles {
		for _, neededRole := range neededRoles {
			if role == string(neededRole) {
				return true
			}
		}
	}
	return false
}

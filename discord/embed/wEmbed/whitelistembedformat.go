package wEmbed

import (
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/discord/session"
	"github.com/Sharktheone/ScharschBot/types"
	"github.com/Sharktheone/ScharschBot/whitelist/whitelist"
	"github.com/Sharktheone/ScharschBot/whitelist/whitelist/utils"
	"github.com/bwmarrin/discordgo"
	"log"
)

var (
	BotAvatarURL string
)

func WhitelistAdding(PlayerName string, i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		member        = types.MemberFromDG(i.Member)
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		maxAccounts   = utils.GetMaxAccounts(member)
		Title         = fmt.Sprintf("%v is now on the whitelist", PlayerName)
		AuthorIconUrl = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", PlayerName)
		AuthorUrl     = fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
		Players       = whitelist.ListedAccountsOf(member.ID, true)
		bannedPlayers = utils.CheckBans(member.ID)
	)

	if !conf.Config.Whitelist.BannedUsersToMaxAccounts {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts (max %v)", username, len(Players), maxAccounts)
	} else {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts and %v are banned (max %v)", username, len(Players), len(bannedPlayers), maxAccounts)
	}
	if conf.Config.Discord.FooterIcon {
		Footer = &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		}
	} else {
		Footer = &discordgo.MessageEmbedFooter{
			Text: FooterText,
		}
	}

	Embed := discordgo.MessageEmbed{
		Title: Title,
		Color: 0x00FF00,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    PlayerName,
			IconURL: AuthorIconUrl,
			URL:     AuthorUrl,
		},
		Footer: Footer,
	}
	return Embed
}

func WhitelistAlreadyListed(PlayerName string, i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		member        = types.MemberFromDG(i.Member)
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		maxAccounts   = utils.GetMaxAccounts(member)
		Title         = fmt.Sprintf("%v is already on the whitelist", PlayerName)
		AuthorIconUrl = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", PlayerName)
		AuthorUrl     = fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
		Players       = whitelist.ListedAccountsOf(member.ID, true)
		bannedPlayers = utils.CheckBans(member.ID)
	)

	if !conf.Config.Whitelist.BannedUsersToMaxAccounts {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts (max %v)", username, len(Players), maxAccounts)
	} else {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts and %v are banned (max %v)", username, len(Players), len(bannedPlayers), maxAccounts)
	}
	if conf.Config.Discord.FooterIcon {
		Footer = &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		}
	} else {
		Footer = &discordgo.MessageEmbedFooter{
			Text: FooterText,
		}
	}
	Embed := discordgo.MessageEmbed{
		Title: Title,
		Color: 0xFF9900,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    PlayerName,
			IconURL: AuthorIconUrl,
			URL:     AuthorUrl,
		},
		Footer: Footer,
	}
	return Embed

}
func WhitelistNotExisting(PlayerName string, i *discordgo.InteractionCreate) discordgo.MessageEmbed {

	var (
		member        = types.MemberFromDG(i.Member)
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		maxAccounts   = utils.GetMaxAccounts(member)
		Title         = fmt.Sprintf("%v is not existing", PlayerName)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
		Players       = whitelist.ListedAccountsOf(member.ID, true)
		bannedPlayers = utils.CheckBans(member.ID)
	)

	if !conf.Config.Whitelist.BannedUsersToMaxAccounts {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts (max %v)", username, len(Players), maxAccounts)
	} else {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts and %v are banned (max %v)", username, len(Players), len(bannedPlayers), maxAccounts)
	}
	if conf.Config.Discord.FooterIcon {
		Footer = &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		}
	} else {
		Footer = &discordgo.MessageEmbedFooter{
			Text: FooterText,
		}
	}

	Embed := discordgo.MessageEmbed{
		Title: Title,
		Color: 0xFF0000,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    PlayerName,
			IconURL: conf.Config.Discord.EmbedErrorIcon,
			URL:     conf.Config.Discord.EmbedErrorAuthorURL,
		},
		Footer: Footer,
	}
	return Embed
}

func WhitelistNoFreeAccounts(PlayerName string, i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		member        = types.MemberFromDG(i.Member)
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		maxAccounts   = utils.GetMaxAccounts(member)
		Title         = "You have no free Accounts"
		AuthorIconUrl = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", PlayerName)
		AuthorUrl     = fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
		Players       = whitelist.ListedAccountsOf(member.ID, true)
		bannedPlayers = utils.CheckBans(member.ID)
	)
	if !conf.Config.Whitelist.BannedUsersToMaxAccounts {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts (max %v)", username, len(Players), maxAccounts)
	} else {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts and %v are banned (max %v)", username, len(Players), len(bannedPlayers), maxAccounts)
	}
	if conf.Config.Discord.FooterIcon {
		Footer = &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		}
	} else {
		Footer = &discordgo.MessageEmbedFooter{
			Text: FooterText,
		}
	}

	Embed := discordgo.MessageEmbed{
		Title: Title,
		Color: 0xFF0000,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    PlayerName,
			IconURL: AuthorIconUrl,
			URL:     AuthorUrl,
		},
		Footer: Footer,
	}
	return Embed
}

func WhitelistAddNotAllowed(PlayerName string, i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		member        = types.MemberFromDG(i.Member)
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		maxAccounts   = utils.GetMaxAccounts(member)
		Title         = fmt.Sprintf("You have no permission to add %v to the whitelist", PlayerName)
		AuthorIconUrl = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", PlayerName)
		AuthorUrl     = fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
		Players       = whitelist.ListedAccountsOf(member.ID, true)
		bannedPlayers = utils.CheckBans(member.ID)
	)
	if !conf.Config.Whitelist.BannedUsersToMaxAccounts {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts (max %v)", username, len(Players), maxAccounts)
	} else {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts and %v are banned (max %v)", username, len(Players), len(bannedPlayers), maxAccounts)
	}
	if conf.Config.Discord.FooterIcon {
		Footer = &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		}
	} else {
		Footer = &discordgo.MessageEmbedFooter{
			Text: FooterText,
		}
	}

	Embed := discordgo.MessageEmbed{
		Title: Title,
		Color: 0xFF0000,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    PlayerName,
			IconURL: AuthorIconUrl,
			URL:     AuthorUrl,
		},
		Footer: Footer,
	}
	return Embed
}

func WhitelistRemoving(PlayerName string, i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		member        = types.MemberFromDG(i.Member)
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		maxAccounts   = utils.GetMaxAccounts(member)
		Title         = fmt.Sprintf("%v is no longer on the whitelist", PlayerName)
		AuthorIconUrl = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", PlayerName)
		AuthorUrl     = fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
		Players       = whitelist.ListedAccountsOf(member.ID, true)
		bannedPlayers = utils.CheckBans(member.ID)
	)
	if !conf.Config.Whitelist.BannedUsersToMaxAccounts {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts (max %v)", username, len(Players), maxAccounts)
	} else {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts and %v are banned (max %v)", username, len(Players), len(bannedPlayers), maxAccounts)
	}
	if conf.Config.Discord.FooterIcon {
		Footer = &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		}
	} else {
		Footer = &discordgo.MessageEmbedFooter{
			Text: FooterText,
		}
	}

	Embed := discordgo.MessageEmbed{
		Title: Title,
		Color: 0x00FF00,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    PlayerName,
			IconURL: AuthorIconUrl,
			URL:     AuthorUrl,
		},
		Footer: Footer,
	}
	return Embed
}
func WhitelistRemoveNotAllowed(PlayerName string, i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		member        = types.MemberFromDG(i.Member)
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		maxAccounts   = utils.GetMaxAccounts(member)
		Title         = fmt.Sprintf("You have no permission to remove %v from the whitelist", PlayerName)
		AuthorIconUrl = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", PlayerName)
		AuthorUrl     = fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
		Players       = whitelist.ListedAccountsOf(member.ID, true)
		bannedPlayers = utils.CheckBans(member.ID)
	)
	if !conf.Config.Whitelist.BannedUsersToMaxAccounts {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts (max %v)", username, len(Players), maxAccounts)
	} else {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts and %v are banned (max %v)", username, len(Players), len(bannedPlayers), maxAccounts)
	}
	if conf.Config.Discord.FooterIcon {
		Footer = &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		}
	} else {
		Footer = &discordgo.MessageEmbedFooter{
			Text: FooterText,
		}
	}

	Embed := discordgo.MessageEmbed{
		Title: Title,
		Color: 0xFF0000,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    PlayerName,
			IconURL: AuthorIconUrl,
			URL:     AuthorUrl,
		},
		Footer: Footer,
	}
	return Embed
}

func WhitelistNotListed(PlayerName string, i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		member        = types.MemberFromDG(i.Member)
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		maxAccounts   = utils.GetMaxAccounts(member)
		Title         = fmt.Sprintf("%v is not on the whitelist", PlayerName)
		AuthorIconUrl = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", PlayerName)
		AuthorUrl     = fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
		Players       = whitelist.ListedAccountsOf(member.ID, true)
		bannedPlayers = utils.CheckBans(member.ID)
	)

	if !conf.Config.Whitelist.BannedUsersToMaxAccounts {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts (max %v)", username, len(Players), maxAccounts)
	} else {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts and %v are banned (max %v)", username, len(Players), len(bannedPlayers), maxAccounts)
	}
	if conf.Config.Discord.FooterIcon {
		Footer = &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		}
	} else {
		Footer = &discordgo.MessageEmbedFooter{
			Text: FooterText,
		}
	}

	Embed := discordgo.MessageEmbed{
		Title: Title,
		Color: 0xFF9900,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    PlayerName,
			IconURL: AuthorIconUrl,
			URL:     AuthorUrl,
		},
		Footer: Footer,
	}
	return Embed
}

func WhitelistIsListedBy(PlayerName string, playerID database.UserID, i *discordgo.InteractionCreate, s *session.Session) discordgo.MessageEmbed {
	var (
		byMember, err = types.MemberFromID(playerID, s)
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		maxAccounts   = utils.GetMaxAccounts(byMember)
		Title         = fmt.Sprintf("%v was whitelisted by", PlayerName)
		Description   = fmt.Sprintf("<@%v>", playerID)
		AuthorIconUrl = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", PlayerName)
		AuthorUrl     = fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
		Players       = whitelist.ListedAccountsOf(playerID, true)
		bannedPlayers = utils.CheckBans(playerID)
	)
	if err != nil {
		log.Printf("Error getting roles: %v", err)
	}

	if !conf.Config.Whitelist.BannedUsersToMaxAccounts {
		FooterText = fmt.Sprintf("%v has whitelisted %v accounts (max %v)", username, len(Players), maxAccounts)
	} else {
		FooterText = fmt.Sprintf("%v has whitelisted %v accounts and %v are banned (max %v)", username, len(Players), len(bannedPlayers), maxAccounts)
	}
	if conf.Config.Discord.FooterIcon {
		Footer = &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		}
	} else {
		Footer = &discordgo.MessageEmbedFooter{
			Text: FooterText,
		}
	}

	Embed := discordgo.MessageEmbed{
		Title:       Title,
		Description: Description,
		Color:       0xFF00,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    PlayerName,
			IconURL: AuthorIconUrl,
			URL:     AuthorUrl,
		},
		Footer: Footer,
	}
	return Embed
}

func WhitelistWhoisNotAllowed(PlayerName string, i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		member        = types.MemberFromDG(i.Member)
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		maxAccounts   = utils.GetMaxAccounts(member)
		Title         = fmt.Sprintf("You have no permission to lookup the owner of %v", PlayerName)
		AuthorIconUrl = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", PlayerName)
		AuthorUrl     = fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
		Players       = whitelist.ListedAccountsOf(member.ID, true)
		bannedPlayers = utils.CheckBans(member.ID)
	)
	if !conf.Config.Whitelist.BannedUsersToMaxAccounts {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts (max %v)", username, len(Players), maxAccounts)
	} else {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts and %v are banned (max %v)", username, len(Players), len(bannedPlayers), maxAccounts)
	}
	if conf.Config.Discord.FooterIcon {
		Footer = &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		}
	} else {
		Footer = &discordgo.MessageEmbedFooter{
			Text: FooterText,
		}
	}

	Embed := discordgo.MessageEmbed{
		Title: Title,
		Color: 0xFF0000,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    PlayerName,
			IconURL: AuthorIconUrl,
			URL:     AuthorUrl,
		},
		Footer: Footer,
	}
	return Embed
}

func WhitelistHasListed(PlayerNames []database.Player, playerID database.UserID, bannedPlayers []database.Player, i *discordgo.InteractionCreate, s *session.Session) discordgo.MessageEmbed {
	var (
		byMember, err = types.MemberFromID(playerID, s)
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		maxAccounts   = utils.GetMaxAccounts(byMember)
		Title         = "Whitelisted accounts of"
		Description   = fmt.Sprintf("<@%v>", playerID)
		embedAccounts []*discordgo.MessageEmbedField
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
	)
	if err != nil {
		log.Printf("Error getting roles: %v", err)
	}

	for _, PlayerName := range PlayerNames {
		userURL := fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		embedAccounts = append(embedAccounts, &discordgo.MessageEmbedField{
			Name:   string(PlayerName),
			Value:  userURL,
			Inline: false,
		})
	}
	for _, PlayerName := range bannedPlayers {
		userURL := fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		embedAccounts = append(embedAccounts, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%v (banned)", PlayerName),
			Value:  userURL,
			Inline: false,
		})
	}

	if !conf.Config.Whitelist.BannedUsersToMaxAccounts {
		FooterText = fmt.Sprintf("%v has whitelisted %v accounts (max %v)", username, len(PlayerNames), maxAccounts)
	} else {
		FooterText = fmt.Sprintf("%v has whitelisted %v accounts and %v are banned (max %v)", username, len(PlayerNames), len(bannedPlayers), maxAccounts)
	}
	if conf.Config.Discord.FooterIcon {
		Footer = &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		}
	} else {
		Footer = &discordgo.MessageEmbedFooter{
			Text: FooterText,
		}
	}

	Embed := discordgo.MessageEmbed{
		Title:       Title,
		Description: Description,
		Color:       0x00FF00,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    username,
			IconURL: avatarURL,
		},
		Fields: embedAccounts,
		Footer: Footer,
	}

	return Embed
}

func WhitelistNoAccounts(i *discordgo.InteractionCreate, playerID string) discordgo.MessageEmbed {
	var (
		username    = i.Member.User.String()
		avatarURL   = i.Member.User.AvatarURL("40")
		Title       = "The following user has no whitelisted accounts:"
		Description = fmt.Sprintf("<@%v>", playerID)
		Embed       = discordgo.MessageEmbed{
			Title:       Title,
			Description: Description,
			Color:       0xFF0000,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    username,
				IconURL: avatarURL,
			},
		}
	)
	return Embed
}

func WhitelistUserNotAllowed(Players []database.Player, playerID database.UserID, bannedPlayers []database.Player, i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		username    = i.Member.User.String()
		avatarURL   = i.Member.User.AvatarURL("40")
		maxAccounts = utils.GetMaxAccounts(types.MemberFromDG(i.Member))
		Title       = "You have no permission to lookup the whitelisted accounts of"
		Description = fmt.Sprintf("<@%v>", playerID)
		FooterText  string
		Footer      *discordgo.MessageEmbedFooter
	)
	if !conf.Config.Whitelist.BannedUsersToMaxAccounts {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts (max %v)", username, len(Players), maxAccounts)
	} else {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts and %v are banned (max %v)", username, len(Players), len(bannedPlayers), maxAccounts)
	}
	if conf.Config.Discord.FooterIcon {
		Footer = &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		}
	} else {
		Footer = &discordgo.MessageEmbedFooter{
			Text: FooterText,
		}
	}

	Embed := discordgo.MessageEmbed{
		Title:       Title,
		Description: Description,
		Color:       0xFF0000,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    username,
			IconURL: avatarURL,
		},
		Footer: Footer,
	}
	return Embed
}

func WhitelistRemoveAllNotAllowed(i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		member        = types.MemberFromDG(i.Member)
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		maxAccounts   = utils.GetMaxAccounts(member)
		Title         = "You have no permission to remove all accounts from the whitelist"
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
		Players       = whitelist.ListedAccountsOf(member.ID, true)
		bannedPlayers = utils.CheckBans(member.ID)
	)
	if !conf.Config.Whitelist.BannedUsersToMaxAccounts {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts (max %v)", username, len(Players), maxAccounts)
	} else {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts and %v are banned (max %v)", username, len(Players), len(bannedPlayers), maxAccounts)
	}
	if conf.Config.Discord.FooterIcon {
		Footer = &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		}
	} else {
		Footer = &discordgo.MessageEmbedFooter{
			Text: FooterText,
		}
	}

	Embed := discordgo.MessageEmbed{
		Title: Title,
		Color: 0xFF0000,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    username,
			IconURL: avatarURL,
			URL:     conf.Config.Discord.EmbedErrorAuthorURL,
		},
		Footer: Footer,
	}
	return Embed
}

func WhitelistRemoveAllNoWhitelistEntries(i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		username  = i.Member.User.String()
		avatarURL = i.Member.User.AvatarURL("40")
		Title     = "There is no whitelist entries to remove"

		Embed = discordgo.MessageEmbed{
			Title: Title,
			Color: 0xFF0000,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    username,
				IconURL: avatarURL,
				URL:     conf.Config.Discord.EmbedErrorAuthorURL,
			},
		}
	)
	return Embed
}

func WhitelistRemoveAllSure(i *discordgo.InteractionCreate) (embed discordgo.MessageEmbed, button discordgo.Button) {
	var (
		username  = i.Member.User.String()
		avatarURL = i.Member.User.AvatarURL("40")
		Title     = "Do you really want to remove all accounts from the whitelist?"

		Embed = discordgo.MessageEmbed{
			Title: Title,
			Color: 0xFF9900,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    username,
				IconURL: avatarURL,
				URL:     conf.Config.Discord.EmbedErrorAuthorURL,
			},
		}
		Button = discordgo.Button{
			Emoji: &discordgo.ComponentEmoji{
				Name: "✅",
			},
			Label:    "Yes, I want to remove all accounts",
			CustomID: "remove_yes",
			Style:    discordgo.DangerButton,
		}
	)
	return Embed, Button
}
func WhitelistRemoveAll(i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		username  = i.Member.User.String()
		avatarURL = i.Member.User.AvatarURL("40")
		Title     = "You have successful removed all accounts from the whitelist"

		Embed = discordgo.MessageEmbed{
			Title: Title,
			Color: 0x00FF00,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    username,
				IconURL: avatarURL,
				URL:     conf.Config.Discord.EmbedErrorAuthorURL,
			},
		}
	)
	return Embed
}

func WhitelistBanUserID(playerID database.UserID, reason string, i *discordgo.InteractionCreate, s *session.Session) discordgo.MessageEmbed {
	var (
		byMember, err = types.MemberFromID(playerID, s)
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		maxAccounts   = utils.GetMaxAccounts(byMember)
		Title         = fmt.Sprintf("Banning following user for the reason %v that has following whitelisted accounts", username)
		Description   = fmt.Sprintf("<@%v>", playerID)
		embedAccounts []*discordgo.MessageEmbedField
		Footer        *discordgo.MessageEmbedFooter
		Players       = whitelist.ListedAccountsOf(playerID, true)
		bannedPlayers = utils.CheckBans(playerID)
	)
	if err != nil {
		log.Printf("Error getting roles: %v", err)
	}
	var FooterText string

	if !conf.Config.Whitelist.BannedUsersToMaxAccounts {
		FooterText = fmt.Sprintf("%v had whitelisted %v accounts (max %v)", username, len(Players), maxAccounts)
	} else {
		FooterText = fmt.Sprintf("%v had whitelisted %v accounts and %v banned (max %v)", username, len(Players), len(bannedPlayers), maxAccounts)
	}
	if conf.Config.Discord.FooterIcon {
		Footer = &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		}
	} else {
		Footer = &discordgo.MessageEmbedFooter{
			Text: FooterText,
		}
	}

	for _, player := range Players {
		userURL := fmt.Sprintf("https://namemc.com/profile/%v", player)
		embedAccounts = append(embedAccounts, &discordgo.MessageEmbedField{
			Name:   string(player),
			Value:  userURL,
			Inline: false,
		})
	}
	var Fields []*discordgo.MessageEmbedField
	Fields = append(Fields, &discordgo.MessageEmbedField{
		Name:  "Reason",
		Value: reason,
	})
	Fields = append(Fields, embedAccounts...)

	Embed := discordgo.MessageEmbed{
		Title:       Title,
		Description: Description,
		Color:       0x00FF00,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    username,
			IconURL: avatarURL,
		},
		Fields: Fields,
		Footer: Footer,
	}
	return Embed
}

func WhitelistBanAccount(PlayerName database.Player, playerID database.UserID, reason string, i *discordgo.InteractionCreate, s *session.Session) discordgo.MessageEmbed {
	var (
		byMember, err = types.MemberFromID(playerID, s)
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		maxAccounts   = utils.GetMaxAccounts(byMember)
		Players       = whitelist.ListedAccountsOf(playerID, true)
		bannedPlayers = utils.CheckBans(playerID)
		Title         = fmt.Sprintf("%v is now banned from the whitelist", PlayerName)
		AuthorIconUrl = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", PlayerName)
		AuthorUrl     = fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		Footer        *discordgo.MessageEmbedFooter
		FooterText    string
		field         discordgo.MessageEmbedField
		reasonField   = discordgo.MessageEmbedField{
			Name:  "Reason",
			Value: reason,
		}
		userID = i.Member.User.ID
	)
	if err != nil {
		log.Printf("Error getting roles: %v", err)
	}
	if len(userID) > 0 {
		FieldName := fmt.Sprintf("%v was whitelisted by", PlayerName)
		FieldValue := fmt.Sprintf("<@%v>", playerID)
		field = discordgo.MessageEmbedField{
			Name:  FieldName,
			Value: FieldValue,
		}
		if !conf.Config.Whitelist.BannedUsersToMaxAccounts {
			FooterText = fmt.Sprintf("%v • He had whitelisted %v accounts (max %v)", username, len(Players), maxAccounts)
		} else {
			FooterText = fmt.Sprintf("%v • He had whitelisted %v accounts and %v banned (max %v)", username, len(Players), len(bannedPlayers), maxAccounts)
		}
		if conf.Config.Discord.FooterIcon {
			Footer = &discordgo.MessageEmbedFooter{
				Text:    FooterText,
				IconURL: avatarURL,
			}
		} else {
			Footer = &discordgo.MessageEmbedFooter{
				Text: FooterText,
			}
		}
	} else {
		FieldName := fmt.Sprintf("%v is not whitelisted", PlayerName)
		field = discordgo.MessageEmbedField{
			Name:  FieldName,
			Value: "The ban will be executed",
		}
	}
	var Embed discordgo.MessageEmbed
	if len(FooterText) > 0 {
		Embed = discordgo.MessageEmbed{
			Title: Title,
			Color: 0x00FF00,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    string(PlayerName),
				IconURL: AuthorIconUrl,
				URL:     AuthorUrl,
			},
			Fields: []*discordgo.MessageEmbedField{
				&reasonField,
				&field,
			},
			Footer: Footer,
		}
	} else {
		Embed = discordgo.MessageEmbed{
			Title: Title,
			Color: 0x00FF00,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    string(PlayerName),
				IconURL: AuthorIconUrl,
				URL:     AuthorUrl,
			},
			Fields: []*discordgo.MessageEmbedField{
				&reasonField,
				&field,
			},
		}
	}

	return Embed
}

func WhitelistUnBanUserID(playerID database.UserID, i *discordgo.InteractionCreate, s *session.Session) discordgo.MessageEmbed {
	var (
		byMember, err = types.MemberFromID(playerID, s)
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		maxAccounts   = utils.GetMaxAccounts(byMember)
		Title         = "Unbanning user"
		Description   = fmt.Sprintf("<@%v>", playerID)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
		Players       = whitelist.ListedAccountsOf(playerID, true)
		bannedPlayers = utils.CheckBans(playerID)
	)
	if err != nil {
		log.Printf("Error getting roles: %v", err)
	}
	if !conf.Config.Whitelist.BannedUsersToMaxAccounts {
		FooterText = fmt.Sprintf("%v • He has whitelisted %v accounts (max %v)", username, len(Players), maxAccounts)
	} else {
		FooterText = fmt.Sprintf("%v • He has whitelisted %v accounts and %v are banned (max %v)", username, len(Players), len(Players), maxAccounts)
	}
	if conf.Config.Discord.FooterIcon {
		Footer = &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		}
	} else {
		Footer = &discordgo.MessageEmbedFooter{
			Text: FooterText,
		}
	}
	var embedAccounts []*discordgo.MessageEmbedField
	for _, playerName := range Players {
		userURL := fmt.Sprintf("https://namemc.com/profile/%v", playerName)
		embedAccounts = append(embedAccounts, &discordgo.MessageEmbedField{
			Name:   string(playerName),
			Value:  userURL,
			Inline: false,
		})
	}
	for _, PlayerName := range bannedPlayers {
		userURL := fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		embedAccounts = append(embedAccounts, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%v (banned)", PlayerName),
			Value:  userURL,
			Inline: false,
		})
	}
	Embed := discordgo.MessageEmbed{
		Title:       Title,
		Description: Description,
		Color:       0x00FF00,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    username,
			IconURL: avatarURL,
		},
		Footer: Footer,
		Fields: embedAccounts,
	}
	return Embed
}

func WhitelistUnBanAccount(playerName database.Player, i *discordgo.InteractionCreate, s *session.Session) discordgo.MessageEmbed {
	var (
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		owner         = whitelist.GetOwner(playerName, s)
		Title         = fmt.Sprintf("%v is now unbanned from the whitelist", playerName)
		AuthorIconUrl = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", playerName)
		AuthorUrl     = fmt.Sprintf("https://namemc.com/profile/%v", playerName)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
	)
	if len(owner.PlayersWithBanned) > 0 {
		FooterText = fmt.Sprintf("%v • 1had whitelisted now %v accounts (max %v)", username, len(owner.PlayersWithBanned), owner.MaxAccounts)
	} else {
		FooterText = fmt.Sprintf("%v • No was not whitelisted", username)
	}
	if conf.Config.Discord.FooterIcon {
		Footer = &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		}
	} else {
		Footer = &discordgo.MessageEmbedFooter{
			Text: FooterText,
		}
	}

	Embed := discordgo.MessageEmbed{
		Title: Title,
		Color: 0x00FF00,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    string(playerName),
			IconURL: AuthorIconUrl,
			URL:     AuthorUrl,
		},
		Footer: Footer,
	}
	return Embed
}

func WhitelistBanAccountNotAllowed(mcName database.Player, i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		member        = types.MemberFromDG(i.Member)
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		maxAccounts   = utils.GetMaxAccounts(member)
		Title         = fmt.Sprintf("You have no permission to (un)ban %v", mcName)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
		Players       = whitelist.ListedAccountsOf(member.ID, true)
		bannedPlayers = utils.CheckBans(member.ID)
	)
	if !conf.Config.Whitelist.BannedUsersToMaxAccounts {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts (max %v)", username, len(Players), maxAccounts)
	} else {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts and %v are banned (max %v)", username, len(Players), len(bannedPlayers), maxAccounts)
	}
	if conf.Config.Discord.FooterIcon {
		Footer = &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		}
	} else {
		Footer = &discordgo.MessageEmbedFooter{
			Text: FooterText,
		}
	}

	Embed := discordgo.MessageEmbed{
		Title: Title,
		Color: 0xFF0000,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    username,
			IconURL: avatarURL,
			URL:     conf.Config.Discord.EmbedErrorAuthorURL,
		},
		Footer: Footer,
	}
	return Embed
}

func WhitelistBanUserIDNotAllowed(playerID database.UserID, i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		member        = types.MemberFromDG(i.Member)
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		maxAccounts   = utils.GetMaxAccounts(member)
		Title         = "You have no permission to (un)ban"
		Description   = fmt.Sprintf("<@%v>", playerID)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
		Players       = whitelist.ListedAccountsOf(member.ID, true)
		bannedPlayers = utils.CheckBans(member.ID)
	)
	if !conf.Config.Whitelist.BannedUsersToMaxAccounts {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts (max %v)", username, len(Players), maxAccounts)
	} else {
		FooterText = fmt.Sprintf("%v • You have whitelisted %v accounts and %v are banned (max %v)", username, len(Players), len(bannedPlayers), maxAccounts)
	}
	if conf.Config.Discord.FooterIcon {
		Footer = &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		}
	} else {
		Footer = &discordgo.MessageEmbedFooter{
			Text: FooterText,
		}
	}

	Embed := discordgo.MessageEmbed{
		Title:       Title,
		Description: Description,
		Color:       0xFF0000,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    username,
			IconURL: avatarURL,
			URL:     conf.Config.Discord.EmbedErrorAuthorURL,
		},
		Footer: Footer,
	}
	return Embed
}

func WhitelistBanned(PlayerName string, IDBan bool, reason string, i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		Title         string
		AuthorName    string
		AuthorURL     string
		AuthorIconURL string
		Description   = fmt.Sprintf("Reason: %v", reason)
	)
	if IDBan {
		Title = " You have no permission to whitelist accounts because you are banned from the whitelist"
		AuthorName = username
		AuthorURL = conf.Config.Discord.EmbedErrorAuthorURL
		AuthorIconURL = avatarURL
	} else {
		Title = fmt.Sprintf("You have no permission to add %v to the whitelist beacause the account banned from the whitelist", PlayerName)
		AuthorName = PlayerName
		AuthorURL = fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		AuthorIconURL = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", PlayerName)
	}

	Embed := discordgo.MessageEmbed{
		Title:       Title,
		Description: Description,
		Color:       0xFF0000,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    AuthorName,
			IconURL: AuthorIconURL,
			URL:     AuthorURL,
		},
	}
	return Embed

}
func WhitelistRemoveMyAccounts(PlayerNames []database.Player, bannedPlayers []database.Player, i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		username      = i.Member.User.String()
		avatarURL     = i.Member.User.AvatarURL("40")
		maxAccounts   = utils.GetMaxAccounts(types.MemberFromDG(i.Member))
		Title         = "Removing whitelisted accounts of"
		playerID      = i.Member.User.ID
		Description   = fmt.Sprintf("<@%v>", playerID)
		embedAccounts []*discordgo.MessageEmbedField
		Footer        *discordgo.MessageEmbedFooter
	)

	for _, PlayerName := range PlayerNames {
		userURL := fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		embedAccounts = append(embedAccounts, &discordgo.MessageEmbedField{
			Name:   string(PlayerName),
			Value:  userURL,
			Inline: false,
		})
	}

	var FooterText string
	if len(bannedPlayers) > 0 {
		FooterText = fmt.Sprintf("You have %v banned accounts (max %v)", len(bannedPlayers), maxAccounts)
	}
	Footer = &discordgo.MessageEmbedFooter{
		Text:    FooterText,
		IconURL: avatarURL,
	}
	var Embed discordgo.MessageEmbed
	if len(FooterText) > 0 {

		Embed = discordgo.MessageEmbed{
			Title:       Title,
			Description: Description,
			Color:       0x00FF00,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    username,
				IconURL: avatarURL,
			},
			Fields: embedAccounts,
			Footer: Footer,
		}
	} else {
		Embed = discordgo.MessageEmbed{
			Title:       Title,
			Description: Description,
			Color:       0x00FF00,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    username,
				IconURL: avatarURL,
			},
			Fields: embedAccounts,
		}
	}

	return Embed
}

func ReportPlayer(PlayerName database.Player, reason string, i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		username    = i.Member.User.String()
		avatarURL   = i.Member.User.AvatarURL("40")
		Title       = fmt.Sprintf("Reported player %v for reason %v", PlayerName, reason)
		Description = "Thanks for reporting a player, we will check it as soon as possible!"
		AuthorName  = PlayerName
		FooterText  = fmt.Sprintf("%v • Reason: %v", username, reason)
		AuthorURL   = fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		AuthorIcon  = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", PlayerName)
	)

	Embed := discordgo.MessageEmbed{
		Title:       Title,
		Description: Description,
		Color:       0x00FF00,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    string(AuthorName),
			IconURL: AuthorIcon,
			URL:     AuthorURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: avatarURL,
		},
	}
	return Embed
}

func ReportNotALlowed(i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		username    = i.Member.User.String()
		avatarURL   = i.Member.User.AvatarURL("40")
		Title       = "You have no permission to report players"
		Description = "You have to be a member of the server to report players"
	)

	Embed := discordgo.MessageEmbed{
		Title:       Title,
		Description: Description,
		Color:       0xFF0000,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    username,
			IconURL: avatarURL,
		},
	}
	return Embed
}
func ReportDisabled(i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		username    = i.Member.User.String()
		avatarURL   = i.Member.User.AvatarURL("40")
		Title       = "Reports are disabled"
		Description = "You can't report players because the reports are disabled"
	)

	Embed := discordgo.MessageEmbed{
		Title:       Title,
		Description: Description,
		Color:       0xFF4653,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    username,
			IconURL: avatarURL,
		},
	}
	return Embed
}
func AlreadyReported(PlayerName database.Player) discordgo.MessageEmbed {
	var (
		Title      = fmt.Sprintf("The player %v is already reported", PlayerName)
		AuthorName = PlayerName
		AuthorURL  = fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		AuthorIcon = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", PlayerName)
	)

	Embed := discordgo.MessageEmbed{
		Title: Title,
		Color: 0xFF0000,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    string(AuthorName),
			IconURL: AuthorIcon,
			URL:     AuthorURL,
		},
	}
	return Embed
}

func NewReport(playerName database.Player, reason string, i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		username    = i.Member.User.String()
		Title       = fmt.Sprintf("New report: %v (report from %v) ", playerName, username)
		Description = "A new player has been reported, owner has listed accounts:"
		AuthorName  = playerName
		FooterText  = fmt.Sprintf("%v • Reason: %v", username, reason)
		AuthorURL   = fmt.Sprintf("https://namemc.com/profile/%v", playerName)
		AuthorIcon  = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", playerName)
		owner       = whitelist.GetOwner(playerName, nil)
		Fields      []*discordgo.MessageEmbedField
		FooterIcon  = i.Member.User.AvatarURL("40")
	)

	log.Printf("Owner of %v is %v", playerName, owner.ID)
	for _, Account := range owner.PlayersWithBanned {
		log.Printf("Account: %v", Account)
		var (
			FieldName string
		)
		mcBanned, _, banReason := utils.CheckBanned(Account, "")
		if mcBanned {
			FieldName = fmt.Sprintf("%v (banned, reason: %v)", Account, banReason)
		} else {
			FieldName = string(Account)
		}

		Fields = append(Fields, &discordgo.MessageEmbedField{
			Name:  FieldName,
			Value: AuthorURL,
		})

	}

	Embed := discordgo.MessageEmbed{
		Title:       Title,
		Description: Description,
		Color:       0xFF6F00,
		Fields:      Fields,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    string(AuthorName),
			IconURL: AuthorIcon,
			URL:     AuthorURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: FooterIcon,
		},
	}
	return Embed
}

func ListReports(i *discordgo.InteractionCreate) discordgo.MessageEmbed {
	var (
		username    = i.Member.User.String()
		avatarURL   = i.Member.User.AvatarURL("40")
		Title       = "List of reported players"
		Description = "Here is a list of all reported players"
		Fields      []*discordgo.MessageEmbedField
		reports     = database.DB.GetReports()
	)

	if len(reports) > 0 {

		for _, r := range reports {
			banned, _, reason := utils.CheckBanned(r.ReportedPlayer, "")
			if banned {
				value := fmt.Sprintf("%v (banned, reason: %v)", r.Reason, reason)
				Fields = append(Fields, &discordgo.MessageEmbedField{
					Name:  string(r.ReportedPlayer),
					Value: value,
				})
			} else {
				Fields = append(Fields, &discordgo.MessageEmbedField{
					Name:  string(r.ReportedPlayer),
					Value: r.Reason,
				})
			}

		}
	} else {
		Description = "There are no reported players"
	}
	Embed := discordgo.MessageEmbed{
		Title:       Title,
		Description: Description,
		Color:       0x6C50FF,
		Fields:      Fields,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    username,
			IconURL: avatarURL,
		},
	}
	return Embed
}

func ReportAction(name string, action string, notifyreporter bool) discordgo.MessageEmbed {
	var (
		avatarURL   = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", name)
		AuthorURL   = fmt.Sprintf("https://namemc.com/profile/%v", name)
		Title       = fmt.Sprintf("Report %v", action)
		Description = fmt.Sprintf("The report has been %v, notifing reporter: %v", action, notifyreporter)
	)

	Embed := discordgo.MessageEmbed{
		Title:       Title,
		Description: Description,
		Color:       0xFFCB00,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    name,
			IconURL: avatarURL,
			URL:     AuthorURL,
		},
	}
	return Embed
}

func ReportUserAction(name string, dmFailed bool, userID database.UserID, s *session.Session, action string) discordgo.MessageEmbed {
	var (
		avatarURL   = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", name)
		AuthorURL   = fmt.Sprintf("https://namemc.com/profile/%v", name)
		Title       = fmt.Sprintf("Your report against %v has been %v", name, action)
		Description string
		FooterText  string
		user, _     = s.GetUserProfile(userID)
		FooterIcon  = user.AvatarURL("40")
	)
	if action == "accepted" {
		Description = "Thank you for your report. The report has been accepted."
	} else {
		Description = "Thank you for your report, but it has been rejected. If you think this is a mistake, please contact a staff member directly."
	}

	if dmFailed {
		FooterText = fmt.Sprintf("Failed to send DM to %v. Maybe you have DMs disabled? Sending to channel instead.", user.User.String())
	} else {
		FooterText = user.User.String()
	}

	Embed := discordgo.MessageEmbed{
		Title:       Title,
		Description: Description,
		Color:       0x00FFC9,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    name,
			IconURL: avatarURL,
			URL:     AuthorURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    FooterText,
			IconURL: FooterIcon,
		},
	}
	return Embed
}

func AlreadyBanned(PlayerName string) discordgo.MessageEmbed {
	var (
		Title      = fmt.Sprintf("The player %v is already banned", PlayerName)
		AuthorName = PlayerName
		AuthorURL  = fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		AuthorIcon = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", PlayerName)
	)

	Embed := discordgo.MessageEmbed{
		Title: Title,
		Color: 0xFF8836,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    AuthorName,
			IconURL: AuthorIcon,
			URL:     AuthorURL,
		},
	}
	return Embed
}

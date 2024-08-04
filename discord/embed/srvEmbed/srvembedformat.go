package srvEmbed

import (
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/discord/session"
	"github.com/Sharktheone/ScharschBot/minecraft/advancements"
	"github.com/Sharktheone/ScharschBot/types"
	"github.com/Sharktheone/ScharschBot/whitelist/whitelist"
	"github.com/bwmarrin/discordgo"
)

var (
	config     = conf.Config
	bansToMax  = config.Whitelist.BannedUsersToMaxAccounts
	footerIcon = config.Discord.FooterIcon
)

func PlayerJoin(e *types.WebsocketEvent, serverConf *conf.Server, footerIconURL, username *string, s *session.Session) *discordgo.MessageEmbed {
	var (
		owner         = whitelist.GetOwner(e.Data.Player, s)
		Title         = fmt.Sprintf("%v joined the game", e.Data.Player)
		AuthorIconUrl = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", e.Data.Player)
		AuthorUrl     = fmt.Sprintf("https://namemc.com/profile/%v", e.Data.Player)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
	)
	if serverConf.SRV.Footer {
		if owner.Whitelisted {
			if !bansToMax {
				FooterText = fmt.Sprintf("%v • The owner has whitelisted %v accounts (max %v)", username, len(owner.Players), owner.MaxAccounts)
			} else {
				FooterText = fmt.Sprintf("%v • The owner has whitelisted %v accounts and %v are banned (max %v)", username, len(owner.Players), len(owner.BannedPlayers), owner.MaxAccounts)
			}
		} else {
			FooterText = fmt.Sprintf("%v is not whitelisted", e.Data.Player)
		}
		if footerIcon {
			Footer = &discordgo.MessageEmbedFooter{
				Text:    FooterText,
				IconURL: *footerIconURL,
			}
		} else {
			Footer = &discordgo.MessageEmbedFooter{
				Text: FooterText,
			}
		}

	}
	var (
		color = 0x00FF00
		Embed *discordgo.MessageEmbed
	)
	if serverConf.SRV.Footer {
		if serverConf.SRV.OneLine {
			Embed = &discordgo.MessageEmbed{
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    Title,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
				Footer: Footer,
			}
		} else {
			Embed = &discordgo.MessageEmbed{
				Title: Title,
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    e.Data.Player,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
				Footer: Footer,
			}
		}
	} else {
		if serverConf.SRV.OneLine {
			Embed = &discordgo.MessageEmbed{
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    Title,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
			}
		} else {
			Embed = &discordgo.MessageEmbed{
				Title: Title,
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    e.Data.Player,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
			}
		}
	}
	return Embed
}
func PlayerQuit(e *types.WebsocketEvent, serverConf *conf.Server, footerIconURL, username *string, s *session.Session) *discordgo.MessageEmbed {
	var (
		owner         = whitelist.GetOwner(e.Data.Player, s)
		Title         = fmt.Sprintf("%v left the game", e.Data.Player)
		AuthorIconUrl = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", e.Data.Player)
		AuthorUrl     = fmt.Sprintf("https://namemc.com/profile/%v", e.Data.Player)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
	)
	if serverConf.SRV.Footer {
		if owner.Whitelisted {
			if !bansToMax {
				FooterText = fmt.Sprintf("%v • The owner has whitelisted %v accounts (max %v)", username, len(owner.Players), owner.MaxAccounts)
			} else {
				FooterText = fmt.Sprintf("%v • The owner has whitelisted %v accounts and %v are banned (max %v)", username, len(owner.Players), len(owner.BannedPlayers), owner.MaxAccounts)
			}
		} else {
			FooterText = fmt.Sprintf("%v is not whitelisted", e.Data.Player)
		}
		if footerIcon {
			Footer = &discordgo.MessageEmbedFooter{
				Text:    FooterText,
				IconURL: *footerIconURL,
			}
		} else {
			Footer = &discordgo.MessageEmbedFooter{
				Text: FooterText,
			}
		}
	}
	var (
		Embed *discordgo.MessageEmbed
		color = 0xFF0000
	)
	if serverConf.SRV.Footer {
		if serverConf.SRV.OneLine {
			Embed = &discordgo.MessageEmbed{
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    Title,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
				Footer: Footer,
			}

		} else {
			Embed = &discordgo.MessageEmbed{
				Title: Title,
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    e.Data.Player,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
				Footer: Footer,
			}
		}
	} else {
		if serverConf.SRV.OneLine {
			Embed = &discordgo.MessageEmbed{
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    Title,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
			}
		} else {
			Embed = &discordgo.MessageEmbed{
				Title: Title,
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    e.Data.Player,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
			}
		}
	}
	return Embed
}

func PlayerAdvancement(e *types.WebsocketEvent, serverConf *conf.Server, footerIconURL, username *string, s *session.Session) discordgo.MessageEmbed {
	var (
		PlayerName    = e.Data.Player
		advancement   = advancements.Decode(e.Data.Advancement)
		owner         = whitelist.GetOwner(PlayerName, s)
		Title         = fmt.Sprintf("%v made the Advancement %v", PlayerName, advancement)
		AuthorIconUrl = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", PlayerName)
		AuthorUrl     = fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
	)
	if serverConf.SRV.Footer {
		if owner.Whitelisted {
			if !bansToMax {
				FooterText = fmt.Sprintf("%v • The owner has whitelisted %v accounts (max %v)", username, len(owner.Players), owner.MaxAccounts)
			} else {
				FooterText = fmt.Sprintf("%v • The owner has whitelisted %v accounts and %v are banned (max %v)", username, len(owner.Players), len(owner.BannedPlayers), owner.MaxAccounts)
			}
		} else {
			FooterText = fmt.Sprintf("%v is not whitelisted", PlayerName)
		}
		if footerIcon {
			Footer = &discordgo.MessageEmbedFooter{
				Text:    FooterText,
				IconURL: *footerIconURL,
			}
		} else {
			Footer = &discordgo.MessageEmbedFooter{
				Text: FooterText,
			}
		}
	}
	var (
		Embed discordgo.MessageEmbed
		color = 0xFFFF00
	)
	if serverConf.SRV.Footer {
		if serverConf.SRV.OneLine {
			Embed = discordgo.MessageEmbed{
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    Title,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
				Footer: Footer,
			}
		} else {
			Embed = discordgo.MessageEmbed{
				Title: Title,
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    PlayerName,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
				Footer: Footer,
			}
		}
	} else {
		if serverConf.SRV.OneLine {
			Embed = discordgo.MessageEmbed{
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    Title,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
			}
		} else {
			Embed = discordgo.MessageEmbed{
				Title: Title,
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    PlayerName,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
			}
		}
	}
	return Embed
}

func PlayerDeath(e *types.WebsocketEvent, serverConf *conf.Server, footerIconURL, username *string, s *session.Session) *discordgo.MessageEmbed {
	var (
		PlayerName    = e.Data.Player
		owner         = whitelist.GetOwner(PlayerName, s)
		Title         = fmt.Sprintf("%v %v", PlayerName, e.Data.DeathMessage)
		AuthorIconUrl = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", PlayerName)
		AuthorUrl     = fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
	)
	if serverConf.SRV.Footer {
		if owner.Whitelisted {
			if !bansToMax {
				FooterText = fmt.Sprintf("%v • The owner has whitelisted %v accounts (max %v)", username, len(owner.Players), owner.MaxAccounts)
			} else {
				FooterText = fmt.Sprintf("%v • The owner has whitelisted %v accounts and %v are banned (max %v)", username, len(owner.Players), len(owner.BannedPlayers), owner.MaxAccounts)
			}
		} else {
			FooterText = fmt.Sprintf("%v is not whitelisted", PlayerName)
		}
		if footerIcon {
			Footer = &discordgo.MessageEmbedFooter{
				Text:    FooterText,
				IconURL: *footerIconURL,
			}
		} else {
			Footer = &discordgo.MessageEmbedFooter{
				Text: FooterText,
			}
		}
	}
	var (
		Embed *discordgo.MessageEmbed
		color = 0x000000
	)
	if serverConf.SRV.Footer {
		if serverConf.SRV.OneLine {
			Embed = &discordgo.MessageEmbed{
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    Title,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
				Footer: Footer,
			}
		} else {
			Embed = &discordgo.MessageEmbed{
				Title: Title,
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    PlayerName,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
				Footer: Footer,
			}
		}
	} else {
		if serverConf.SRV.OneLine {
			Embed = &discordgo.MessageEmbed{
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    Title,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
			}
		} else {
			Embed = &discordgo.MessageEmbed{
				Title: Title,
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    PlayerName,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
			}
		}
	}
	return Embed
}

func Chat(eventJson *types.WebsocketEvent, serverConf *conf.Server, footerIconURL, username *string, s *session.Session) *discordgo.MessageEmbed {
	var (
		PlayerName    = eventJson.Data.Player
		Message       = eventJson.Data.Message
		owner         = whitelist.GetOwner(PlayerName, s)
		AuthorIconUrl = fmt.Sprintf("https://mc-heads.net/avatar/%v.png", PlayerName)
		AuthorUrl     = fmt.Sprintf("https://namemc.com/profile/%v", PlayerName)
		FooterText    string
		Footer        *discordgo.MessageEmbedFooter
	)
	if serverConf.SRV.Footer {
		if owner.Whitelisted {
			if !bansToMax {
				FooterText = fmt.Sprintf("%v • The owner has whitelisted %v accounts (max %v)", username, len(owner.Players), owner.MaxAccounts)
			} else {
				FooterText = fmt.Sprintf("%v • The owner has whitelisted %v accounts and %v are banned (max %v)", username, len(owner.Players), len(owner.BannedPlayers), owner.MaxAccounts)
			}
		} else {
			FooterText = fmt.Sprintf("%v is not whitelisted", PlayerName)
		}
		if footerIcon {
			Footer = &discordgo.MessageEmbedFooter{
				Text:    FooterText,
				IconURL: *footerIconURL,
			}
		} else {
			Footer = &discordgo.MessageEmbedFooter{
				Text: FooterText,
			}
		}
	}
	var (
		Embed *discordgo.MessageEmbed
		color = 0x00AAFF
	)
	if serverConf.SRV.Footer {
		if serverConf.SRV.OneLine {
			Embed = &discordgo.MessageEmbed{
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    Message,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
				Footer: Footer,
			}
		} else {
			Embed = &discordgo.MessageEmbed{
				Title: Message,
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    PlayerName,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
				Footer: Footer,
			}
		}
	} else {
		if serverConf.SRV.OneLine {
			Embed = &discordgo.MessageEmbed{
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    Message,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
			}
		} else {
			Embed = &discordgo.MessageEmbed{
				Title: Message,
				Color: color,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    PlayerName,
					IconURL: AuthorIconUrl,
					URL:     AuthorUrl,
				},
			}
		}
	}
	return Embed
}

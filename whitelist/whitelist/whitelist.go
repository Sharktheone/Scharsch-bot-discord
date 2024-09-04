package whitelist

import (
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/discord/embed/banEmbed"
	"github.com/Sharktheone/ScharschBot/discord/session"
	"github.com/Sharktheone/ScharschBot/pterodactyl"
	"github.com/Sharktheone/ScharschBot/types"
	"github.com/Sharktheone/ScharschBot/whitelist"
	"github.com/bwmarrin/discordgo"
	"log"
)

var (
	config             = conf.Config
	addCommand         = config.Pterodactyl.WhitelistAddCommand
	removeCommand      = config.Pterodactyl.WhitelistRemoveCommand
	pterodactylEnabled = config.Pterodactyl.Enabled
)

type Player struct {
	ID                database.UserID
	Whitelisted       bool
	Name              database.Player
	Players           []database.Player
	PlayersWithBanned []database.Player
	BannedPlayers     []database.Player
	Roles             []database.Role
	MaxAccounts       int
}

type AddResult int

const (
	AlreadyListed AddResult = iota
	NotExisting
	NoFreeAccount
	NotAllowed
	McBanned
	DcBanned
	BothBanned
	Ok
)

func Add(player database.Player, member *types.Member) (AddResult, string) {
	mcBan, dcBan, reason := CheckBanned(player, member)

	if mcBan && dcBan {
		return BothBanned, reason
	}
	if mcBan {
		return McBanned, reason
	}
	if dcBan {
		return DcBanned, reason
	}

	if !CheckRoles(member, config.Whitelist.Roles.ServerRoleID) {
		return NotAllowed, ""
	}

	if !HasFreeAccount(player, member) {
		return NoFreeAccount, ""
	}

	if !AccountExists(player) {
		return NotExisting, ""
	}

	if database.DB.IsWhitelisted(player) {
		return AlreadyListed, ""
	}

	whitelist.Provider.AddToWhitelist(player, member)

	log.Printf("%v is adding %v to whitelist", member.ID, player)

	return Ok, ""
}

func Remove(username database.Player, member *types.Member) (allowed bool, onWhitelist bool) {
	entry, found := database.DB.GetWhitelistedPlayer(username)
	if !found {
		return false, false
	}

	if entry.ID != member.ID {
		if !CheckRoles(member, config.Whitelist.Roles.ServerRoleID) {
			return false, false
		}
	}

	database.DB.RemoveAccount(username)
	if pterodactylEnabled {
		command := fmt.Sprintf(removeCommand, username)
		for _, listedServer := range config.Whitelist.Servers {
			for _, server := range config.Pterodactyl.Servers {
				if server.ServerName == listedServer {
					if err := pterodactyl.SendCommand(command, server.ServerID); err != nil {
						log.Printf("Failed to send command to server %v: %v", server.ServerID, err)
					}
				}
			}
		}
	}
	log.Printf("%v is removing %v from whitelist", member.ID, username)
	return true, true
}

func RemoveAll(member *types.Member) (allowed bool, onWhitelist bool) {
	if !CheckRoles(member, config.Discord.WhitelistRemoveRoleID) {
		return false, false
	}

	entries := database.DB.GetAllWhitelistedPlayers()

	log.Printf("%v is removing all accounts from whitelist", member.ID)
	for _, entry := range entries {
		database.DB.RemoveAccount(entry.Name)

		if pterodactylEnabled {
			command := fmt.Sprintf(removeCommand, entry.Name)
			for _, listedServer := range config.Whitelist.Servers {
				for _, server := range config.Pterodactyl.Servers {
					if server.ServerName == listedServer {
						if err := pterodactyl.SendCommand(command, server.ServerID); err != nil {
							log.Printf("Failed to send command to server %v: %v", server.ServerID, err)
						}
					}
				}
			}
		}

	}

	return true, len(entries) > 0
}
func RemoveAllAllowed(member *types.Member) (allowed bool) {
	return CheckRoles(member, config.Discord.WhitelistRemoveRoleID)

}

func Whois(username database.Player, member *types.Member) (dcUserID database.UserID, allowed bool, found bool) {
	if !CheckRoles(member, config.Discord.WhitelistWhoisRoleID) {
		return "", false, false
	}

	log.Printf("%v is looking who whitelisted %v ", member.ID, username)
	result, dataFound := database.DB.GetWhitelistedPlayer(username)

	return result.ID, true, dataFound
}

func HasListed(lookupID database.UserID, member *types.Member, isSelfLookup bool) (accounts []database.Player, allowed bool, found bool, bannedPlayers []string) {
	if !isSelfLookup && !CheckRoles(member, config.Discord.WhitelistWhoisRoleID) {
		return nil, false, false, nil
	}
	log.Printf("%v is looking on whitelisted accounts of %v ", member.ID, lookupID)

	players := database.DB.Players(lookupID)
	return players, true, len(players) > 0, CheckBans(member.ID)
}

func ListedAccountsOf(userID database.UserID, banned bool) (Accounts []database.Player) {
	var (
		lastIndex  = -1
		datalen    = 0
		resultsban []database.PlayerBanData
	)
	results := database.DB.Players(userID)
	datalen += len(results)
	if banned {
		resultsban = database.DB.BannedPlayers(userID)
		datalen += len(resultsban)
	}

	if datalen > 0 {
		listedAccounts := make([]database.Player, datalen)
		for i, result := range results {
			listedAccounts[i] = result
			lastIndex = i
		}
		if banned {
			for i, result := range resultsban {
				listedAccounts[lastIndex+i+1] = result.Player
			}
		}
		return listedAccounts
	} else {
		return
	}
}

func BanUserID(member *types.Member, banID database.UserID, banAccounts bool, reason string, s *session.Session) (allowed bool, alreadyBanned bool) {

	if !CheckRoles(member, config.Discord.WhitelistBanRoleID) {
		return false, false
	}

	listedAccounts := ListedAccountsOf(banID, false)

	_, banned, _ := CheckBanned("", banID)
	if banned {
		return true, true
	}

	log.Printf("%v is banning %v", member.ID, banID)

	database.DB.BanUser(banID, reason)

	if banAccounts {
		for _, account := range listedAccounts {
			database.DB.UnWhitelistPlayer(account)

			whitelist.Provider.BanPlayer(banID, account, reason)

		}
	}

	database.DB.BanUser(banID, reason)
	messageEmbedDM := banEmbed.DMBan(false, string(banID), reason, s)
	messageEmbedDMFailed := banEmbed.DMBan(true, string(banID), reason, s)
	if err := s.SendDM(string(banID), &discordgo.MessageSend{
		Embed: &messageEmbedDM,
	}, &discordgo.MessageSend{
		Content: fmt.Sprintf("<@%v>", banID),
		Embed:   &messageEmbedDMFailed,
	},
	); err != nil {
		log.Printf("Failed to send DM to %v: %v", banID, err)
	}
	return true, false
}

func BanAccount(member *types.Member, account database.Player, reason string, s *session.Session) (bool, *Player) {
	if !CheckRoles(member, config.Discord.WhitelistBanRoleID) {
		return false, nil
	}

	owner := GetOwner(account, s)
	if owner.Whitelisted {
		alreadyBanned := database.DB.IsPlayerBanned(account)

		if !alreadyBanned {
			log.Printf("%v is banning %v", owner.ID, account)
			database.DB.BanPlayer(account, reason)

			database.DB.UnWhitelistPlayer(account)

			messageEmbedDM := banEmbed.DMBanAccount(string(account), false, string(owner.ID), reason, s)
			messageEmbedDMFailed := banEmbed.DMBanAccount(string(account), true, string(owner.ID), reason, s)
			if err := s.SendDM(string(owner.ID), &discordgo.MessageSend{
				Embed: &messageEmbedDM,
			}, &discordgo.MessageSend{
				Content: fmt.Sprintf("<@%v>", owner.ID),
				Embed:   &messageEmbedDMFailed,
			},
			); err != nil {
				log.Printf("Failed to send DM to %v: %v", owner.ID, err)
			}
			if pterodactylEnabled {
				command := fmt.Sprintf(removeCommand, account)
				for _, listedServer := range config.Whitelist.Servers {
					for _, server := range config.Pterodactyl.Servers {
						if server.ServerName == listedServer {
							if err := pterodactyl.SendCommand(command, server.ServerID); err != nil {
								log.Printf("Failed to send command to server %v: %v", server.ServerID, err)
							}
						}
					}
				}
			}
		}
	} else {
		return false, nil
	}

	return true, owner
}
func UnBanUserID(userID string, roles []string, banID string, unbanAccounts bool, s *session.Session) (allowed bool) {
	unBanAllowed := false
	for _, role := range roles {
		for _, neededRole := range config.Discord.WhitelistBanRoleID {
			if role == neededRole {
				unBanAllowed = true
				break
			}
		}
	}
	if unBanAllowed {
		log.Printf("%v is unbanning %v", userID, banID)
		database.DB.UnBanUser(database.UserID(banID))
		if unbanAccounts {
			result := database.DB.BannedPlayers(database.UserID(banID))

			for _, entry := range result {
				database.DB.UnBanPlayerFrom(database.UserID(banID), entry.Player)

			}
			messageEmbedDM := banEmbed.DMUnBan(false, banID, s)
			messageEmbedDMFailed := banEmbed.DMUnBan(true, banID, s)
			if err := s.SendDM(banID, &discordgo.MessageSend{
				Embed: &messageEmbedDM,
			}, &discordgo.MessageSend{
				Content: fmt.Sprintf("<@%v>", banID),
				Embed:   &messageEmbedDMFailed,
			},
			); err != nil {
				log.Printf("Failed to send DM to %v: %v", banID, err)
			}
		}
	}
	return unBanAllowed
}

func UnBanAccount(userID string, roles []string, account string, s *session.Session) (allowed bool) {
	unBanAllowed := false
	for _, role := range roles {
		for _, neededRole := range config.Discord.WhitelistBanRoleID {
			if role == neededRole {
				unBanAllowed = true
				break
			}
		}
	}
	if unBanAllowed {
		log.Printf("%v is unbanning %v", userID, account)
		database.DB.UnBanPlayer(database.Player(account))
		messageEmbedDM := banEmbed.DMUnBanAccount(account, false, userID, s)
		messageEmbedDMFailed := banEmbed.DMUnBanAccount(account, true, userID, s)
		if err := s.SendDM(userID, &discordgo.MessageSend{
			Embed: &messageEmbedDM,
		}, &discordgo.MessageSend{
			Content: fmt.Sprintf("<@%v>", userID),
			Embed:   &messageEmbedDMFailed,
		},
		); err != nil {
			log.Printf("Failed to send DM to %v: %v", userID, err)
		}

	}

	return unBanAllowed
}

func RemoveMyAccounts(userID string) (hadListedAccounts bool, listedAccounts []string) {

	var (
		accounts          = ListedAccountsOf(userID, false)
		hasListedAccounts = false
	)
	if len(accounts) > 0 {
		hasListedAccounts = true
		log.Printf("%v is removing his own accounts from the whitelist", userID)
		for _, account := range accounts {
			found := database.DB.IsWhitelistedBy(database.UserID(userID), database.Player(account))

			if found {
				database.DB.RemoveAccount(database.Player(account))
				if pterodactylEnabled {
					command := fmt.Sprintf(removeCommand, account)
					for _, listedServer := range config.Whitelist.Servers {
						for _, server := range config.Pterodactyl.Servers {
							if server.ServerName == listedServer {
								if err := pterodactyl.SendCommand(command, server.ServerID); err != nil {
									log.Printf("Error while sending command to server %v: %v", server.ServerID, err)
								}
							}
						}
					}
				}
			}
		}
	}

	return hasListedAccounts, accounts
}

func GetOwner(Account database.Player, s *session.Session) *Player {
	var (
		dcUser database.UserID
	)
	result, found := database.DB.GetWhitelistedPlayer(Account)
	if found {
		dcUser = result.ID
	} else {
		result, found := database.DB.GetPlayerBan(Account)
		if found {
			dcUser = result.ID
		}
	}

	if dcUser != "" {
		var roles []database.Role
		if s != nil {
			rls, err := s.GetRoles(string(dcUser))
			if err != nil {
				log.Printf("Error while getting roles of %v: %v", dcUser, err)
			}

			for _, role := range rls {
				roles = append(roles, database.Role(role))
			}

		}

		member := types.Member{
			ID:       dcUser,
			Roles:    roles,
			Username: "<unknown>",
		}

		return &Player{
			ID:                dcUser,
			Whitelisted:       true,
			Name:              Account,
			Players:           ListedAccountsOf(dcUser, false),
			PlayersWithBanned: ListedAccountsOf(dcUser, true),
			BannedPlayers:     CheckBans(dcUser),
			Roles:             roles,
			MaxAccounts:       GetMaxAccounts(&member),
		}
	}
	return &Player{
		ID:            dcUser,
		Whitelisted:   false,
		Name:          Account,
		Players:       nil,
		BannedPlayers: nil,
		Roles:         nil,
		MaxAccounts:   0,
	}
}

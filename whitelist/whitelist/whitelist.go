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
	ID                string
	Whitelisted       bool
	Name              string
	Players           []string
	PlayersWithBanned []string
	BannedPlayers     []string
	Roles             []string
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

	if removeAllowed {
		log.Printf("%v is removing all accounts from whitelist", userID)
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

	}

	return removeAllowed, len(entries) > 0
}
func RemoveAllAllowed(member *types.Member) (allowed bool) {
	return CheckRoles(member, config.Discord.WhitelistRemoveRoleID)

}

func Whois(username string, userID string, roles []string) (dcUserID string, allowed bool, found bool) {
	var whoisAllowed = false
	for _, role := range roles {
		for _, neededRole := range config.Discord.WhitelistWhoisRoleID {
			if role == neededRole {
				whoisAllowed = true
				break
			}
		}
	}
	var (
		dcUser    string
		dataFound bool
	)
	if whoisAllowed {
		log.Printf("%v is looking who whitelisted %v ", userID, username)
		result, dataFound := database.DB.GetWhitelistedPlayer(database.Player(username))

		if dataFound {
			dcUser = fmt.Sprintf("%v", result.ID)
		}
	}
	return dcUser, whoisAllowed, dataFound
}
func HasListed(lookupID string, userID string, roles []string, isSelfLookup bool) (accounts []string, allowed bool, found bool, bannedPlayers []string) {
	var listedAllowed = false
	for _, role := range roles {
		// TODO Add new Role
		for _, neededRole := range config.Discord.WhitelistRemoveRoleID {
			if role == neededRole {
				listedAllowed = true
				break
			}
		}
	}
	if isSelfLookup && !listedAllowed {
		session.HasRoleID(roles, config.Discord.WhitelistServerRoleID)
	}
	var listedAcc []string
	if listedAllowed {
		log.Printf("%v is looking on whitelisted accounts of %v ", userID, lookupID)

		players := database.DB.Players(database.UserID(lookupID))
		listedAcc = make([]string, len(players))
		for i, player := range players {
			listedAcc[i] = string(player)
		}
	}
	return listedAcc, listedAllowed, len(listedAcc) > 0, CheckBans(userID)
}

func ListedAccountsOf(userID string, banned bool) (Accounts []string) {
	var (
		lastIndex = -1
		datalen   = 0
	)
	results := database.DB.Players(database.UserID(userID))
	resultsban := database.DB.BannedPlayers(database.UserID(userID))
	datalen += len(results)
	if banned {
		datalen += len(resultsban)
	}
	if datalen > 0 {
		listedAccounts := make([]string, datalen)
		for i, result := range results {
			listedAccounts[i] = string(result)
			lastIndex = i
		}
		if banned {
			for i, result := range resultsban {
				listedAccounts[lastIndex+i+1] = string(result.Player)
			}
		}
		return listedAccounts
	} else {
		return
	}
}

func BanUserID(userID string, roles []string, banID string, banAccounts bool, reason string, s *session.Session) (allowed bool, alreadyBanned bool) {
	banAllowed := false
	listedAccounts := ListedAccountsOf(banID, false)
	for _, role := range roles {
		for _, neededRole := range config.Discord.WhitelistBanRoleID {
			if role == neededRole {
				banAllowed = true
				break
			}
		}
	}
	if banAllowed {
		_, banned, _ := CheckBanned("", banID)
		if banned {
			return true, true
		} else {

			log.Printf("%v is banning %v", userID, banID)

			database.DB.BanUser(database.UserID(banID), reason)

			if banAccounts {
				for _, account := range listedAccounts {
					database.DB.UnWhitelistPlayer(database.Player(account))
					if pterodactylEnabled {
						command := fmt.Sprintf(removeCommand, account)
						for _, listedServer := range config.Whitelist.Servers {
							for _, server := range config.Pterodactyl.Servers {
								if server.ServerName == listedServer {
									if err := pterodactyl.SendCommand(command, server.ServerID); err != nil {
										log.Printf("Failed to send command to server %v: %v", server.ServerName, err)
									}
								}
							}
						}
					}

					database.DB.BanUser(database.UserID(banID), reason)
				}
				messageEmbedDM := banEmbed.DMBan(false, banID, reason, s)
				messageEmbedDMFailed := banEmbed.DMBan(true, banID, reason, s)
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
		return banAllowed, false
	}
	return
}

func BanAccount(userID string, roles []string, account string, reason string, s *session.Session) (bool, *Player) {
	var (
		banAllowed = false
	)
	for _, role := range roles {
		for _, neededRole := range config.Discord.WhitelistBanRoleID {
			if role == neededRole {
				banAllowed = true
				break
			}
		}
	}

	owner := GetOwner(account, s)
	if owner.Whitelisted {
		alreadyBanned := database.DB.IsPlayerBanned(database.Player(account))

		if banAllowed && !alreadyBanned {
			log.Printf("%v is banning %v", userID, account)
			database.DB.BanPlayer(database.Player(account), reason)

			database.DB.UnWhitelistPlayer(database.Player(account))

			messageEmbedDM := banEmbed.DMBanAccount(account, false, owner.ID, reason, s)
			messageEmbedDMFailed := banEmbed.DMBanAccount(account, true, owner.ID, reason, s)
			if err := s.SendDM(owner.ID, &discordgo.MessageSend{
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

	return banAllowed, owner
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

func GetOwner(Account string, s *session.Session) *Player {
	var (
		dcUser string
	)
	result, found := database.DB.GetWhitelistedPlayer(database.Player(Account))
	if found {
		dcUser = fmt.Sprintf("%v", result.ID)
	} else {
		result, found := database.DB.GetPlayerBan(database.Player(Account))
		if found {
			dcUser = fmt.Sprintf("%v", result.ID)
		}
	}

	if dcUser != "" {
		var roles []string
		if s != nil {
			var err error
			roles, err = s.GetRoles(dcUser)
			if err != nil {
				log.Printf("Error while getting roles of %v: %v", dcUser, err)
			}
		}
		return &Player{
			ID:                dcUser,
			Whitelisted:       true,
			Name:              Account,
			Players:           ListedAccountsOf(dcUser, false),
			PlayersWithBanned: ListedAccountsOf(dcUser, true),
			BannedPlayers:     CheckBans(dcUser),
			Roles:             roles,
			MaxAccounts:       GetMaxAccounts(roles),
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

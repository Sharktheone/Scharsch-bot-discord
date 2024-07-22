package whitelist

import (
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/database/mongodb"
	"github.com/Sharktheone/ScharschBot/discord/embed/banEmbed"
	"github.com/Sharktheone/ScharschBot/discord/session"
	"github.com/Sharktheone/ScharschBot/pterodactyl"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	config              = conf.GetConf()
	whitelistCollection = config.Whitelist.Mongodb.MongodbWhitelistCollectionName
	banCollection       = config.Whitelist.Mongodb.MongodbBanCollectionName
	addCommand          = config.Pterodactyl.WhitelistAddCommand
	removeCommand       = config.Pterodactyl.WhitelistRemoveCommand
	pterodactylEnabled  = config.Pterodactyl.Enabled
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

func Add(username string, userID string, roles []string) (alreadyListed bool, existing bool, accountFree bool, allowed bool, mcBanned bool, dcBanned bool, banReason string) {
	var addAllowed = false
	mcBan, dcBan, reason := CheckBanned(username, userID)
	if !mcBan && !dcBan {
		for _, role := range roles {
			for _, neededRole := range config.Whitelist.Roles.ServerRoleID {
				if role == neededRole {
					addAllowed = true
					break
				}
			}
		}
	}
	var hasFreeAccount = false
	result, _ := mongodb.Read(whitelistCollection, bson.M{"dcUserID": userID})
	if GetMaxAccounts(roles) <= (len(result) + len(CheckBans(userID))) {
		hasFreeAccount = false
	} else {
		hasFreeAccount = true
	}
	var found bool
	existingAcc := existingAccount(username)
	if existingAcc && hasFreeAccount && addAllowed {
		_, found = mongodb.Read(whitelistCollection, bson.M{
			"mcAccount": username,
		})
		if !found {

			database.DB.WhitelistPlayer(database.UserID(userID), database.Player(username), []database.Role{}) //TODO: Add Roles
			if pterodactylEnabled {
				// TODO: Update to use new pterodactyl package (WS)
				// TODO: Add Waitlist if server is offline
				command := fmt.Sprintf(addCommand, username)
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
			log.Printf("%v is adding %v to whitelist", userID, username)
		}

	}
	return found, existingAcc, hasFreeAccount, addAllowed, mcBan, dcBan, reason
}

func Remove(username string, userID string, roles []string) (allowed bool, onWhitelist bool) {
	var removeAllowed = false
	for _, role := range roles {
		for _, neededRole := range config.Discord.WhitelistRemoveRoleID {
			if role == neededRole {
				removeAllowed = true
				break
			}
		}
	}
	entry, found := mongodb.Read(whitelistCollection, bson.M{
		"dcUserID":  userID,
		"mcAccount": username,
	})
	if !removeAllowed {
		if entry[0]["dcUserID"] == userID && found {
			removeAllowed = true
		}
	}
	if removeAllowed && found {
		mongodb.Remove(whitelistCollection, bson.M{
			"mcAccount": username,
		})
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
		log.Printf("%v is removing %v from whitelist", userID, username)
	}
	return removeAllowed, found
}

func RemoveAll(userID string, roles []string) (allowed bool, onWhitelist bool) {
	var removeAllowed = false
	for _, role := range roles {
		for _, neededRole := range config.Discord.WhitelistRemoveRoleID {
			if role == neededRole {
				removeAllowed = true
				break
			}
		}
	}
	entries, found := mongodb.Read(whitelistCollection, bson.M{
		"dcUserID":  bson.M{"$exists": true},
		"mcAccount": bson.M{"$exists": true},
	})

	if removeAllowed && found {
		log.Printf("%v is removing all accounts from whitelist", userID)
		for _, entry := range entries {
			mongodb.Remove(whitelistCollection, bson.M{
				"mcAccount": entry["mcAccount"],
			})
			if pterodactylEnabled {
				command := fmt.Sprintf(removeCommand, entry["mcAccount"])
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

	return removeAllowed, found
}
func RemoveAllAllowed(roles []string) (allowed bool) {
	var removeAllowed = false
	for _, role := range roles {
		for _, neededRole := range config.Discord.WhitelistRemoveRoleID {
			if role == neededRole {
				removeAllowed = true
				break
			}
		}
	}
	return removeAllowed
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
		var result []bson.M
		result, dataFound = mongodb.Read(whitelistCollection, bson.M{
			"mcAccount": username,
		})
		if dataFound {
			dcUser = fmt.Sprintf("%v", result[0]["dcUserID"])
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
	var (
		results   []bson.M
		dataFound bool
	)
	var listedAcc []string
	if listedAllowed {
		log.Printf("%v is looking on whitelisted accounts of %v ", userID, lookupID)

		results, dataFound = mongodb.Read(whitelistCollection, bson.M{
			"dcUserID": lookupID,
		})
		listedAccounts := make([]string, len(results))
		for i, result := range results {
			listedAccounts[i] = fmt.Sprintf("%v", result["mcAccount"])

		}
		listedAcc = listedAccounts
	}
	return listedAcc, listedAllowed, dataFound, CheckBans(userID)
}

func existingAccount(username string) (existing bool) {
	url := fmt.Sprintf("https://api.mojang.com/users/profiles/minecraft/%v", username)
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to make check account existebility: %v\n", err)
	}
	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Printf("Failed reading Body white account check: %v\n", err)
	}
	return len(string(body)) > 0

}
func ListedAccountsOf(userID string, banned bool) (Accounts []string) {
	var (
		lastIndex = -1
		datalen   = 0
	)
	results, dataFound := mongodb.Read(whitelistCollection, bson.M{
		"dcUserID": userID,
	})
	resultsban, dataFoundban := mongodb.Read(banCollection, bson.M{
		"dcUserID":  userID,
		"mcAccount": bson.M{"$exists": true},
	})
	if dataFound {
		datalen += len(results)
	}
	if dataFoundban && banned {
		datalen += len(resultsban)
	}
	if datalen > 0 {
		listedAccounts := make([]string, datalen)
		if dataFound {
			for i, result := range results {
				listedAccounts[i] = fmt.Sprintf("%v", result["mcAccount"])
				lastIndex = i
			}
		}
		if dataFoundban && banned {
			for i, result := range resultsban {
				listedAccounts[lastIndex+i+1] = fmt.Sprintf("%v", result["mcAccount"])
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
		_, alreadyBanned := mongodb.Read(banCollection, bson.M{
			"mcAccount": account,
		})

		if banAllowed && !alreadyBanned {
			log.Printf("%v is banning %v", userID, account)
			database.DB.BanPlayer(database.UserID(owner.ID), database.Player(account), reason)

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
		mongodb.Remove(banCollection, bson.M{
			"dcUserID":  banID,
			"mcAccount": bson.M{"$exists": false},
		})
		if unbanAccounts {
			var result []bson.M
			result, dataFound := mongodb.Read(banCollection, bson.M{
				"dcUserID": banID,
			})
			if dataFound {
				for _, entry := range result {
					mongodb.Remove(banCollection, bson.M{
						"mcAccount": entry["mcAccount"],
						"dcUserID":  banID,
					})
				}
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
		mongodb.Remove(banCollection, bson.M{
			"mcAccount": account,
		})
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

func CheckBanned(mcName string, userID string) (mcBanned bool, dcBanned bool, banReason string) {
	var (
		dataFound bool
		mcData    []bson.M
		dcData    []bson.M
		mc        = false
		dc        = false
		reason    string
	)
	mcData, dataFound = mongodb.Read(banCollection, bson.M{
		"mcAccount": mcName,
	})
	if dataFound {
		mc = true
	}

	dcData, dataFound = mongodb.Read(banCollection, bson.M{
		"dcUserID":  userID,
		"mcAccount": bson.M{"$exists": false},
	})
	if dataFound {
		dc = true
	}
	if mc {
		reason = fmt.Sprintf("%v", mcData[0]["reason"])
	}
	if dc {
		reason += fmt.Sprintf("%v", dcData[0]["reason"])
	}

	return mc, dc, reason
}

func CheckBans(userID string) (bannedPlayers []string) {
	var (
		dataFound bool
		results   []bson.M
	)
	results, dataFound = mongodb.Read(banCollection, bson.M{
		"dcUserID":  userID,
		"mcAccount": bson.M{"$exists": true},
	})

	if dataFound {
		var bannedAccounts = make([]string, len(results))
		for i, result := range results {
			bannedAccounts[i] = fmt.Sprintf("%v", result["mcAccount"])

		}
		return bannedAccounts
	} else {
		return
	}
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
			_, found := mongodb.Read(whitelistCollection, bson.M{
				"dcUserID":  userID,
				"mcAccount": account,
			})
			if found {
				mongodb.Remove(whitelistCollection, bson.M{
					"mcAccount": account,
				})
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
		dataFound bool
		result    []bson.M
		dcUser    string
	)
	result, dataFound = mongodb.Read(whitelistCollection, bson.M{
		"mcAccount": strings.ToLower(Account),
	})
	if dataFound {
		dcUser = fmt.Sprintf("%v", result[0]["dcUserID"])
	} else {
		result, dataFound = mongodb.Read(banCollection, bson.M{
			"dcUserID":  bson.M{"$exists": true},
			"mcAccount": Account,
		})
		if dataFound {
			dcUser = fmt.Sprintf("%v", result[0]["dcUserID"])
		}
	}
	if dataFound {
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
			Whitelisted:       dataFound,
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
		Whitelisted:   dataFound,
		Name:          Account,
		Players:       nil,
		BannedPlayers: nil,
		Roles:         nil,
		MaxAccounts:   0,
	}
}

func GetMaxAccounts(roleIDs []string) (maxAccounts int) {
	var max = 0
	for _, entry := range config.Whitelist.MaxAccounts {
		for _, role := range roleIDs {
			if entry.RoleID == role && entry.Max > max {
				max = entry.Max
			}
		}
	}
	return max
}

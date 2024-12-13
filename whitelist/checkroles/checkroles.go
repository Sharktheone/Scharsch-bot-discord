package checkroles

import (
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/discord/bot"
	"github.com/Sharktheone/ScharschBot/pterodactyl"
	"log"
)

func CheckRoles() {
	if conf.Config.Whitelist.KickUnWhitelisted {
		for _, server := range pterodactyl.Servers {
			for _, player := range server.OnlinePlayers.Players {
				found := database.DB.IsWhitelisted(database.Player(*player))
				if !found {
					command := fmt.Sprintf(conf.Config.Whitelist.KickCommand, player)
					if err := pterodactyl.SendCommand(command, server.Config.ServerID); err != nil {
						log.Printf("Failed to kick %v from server %v: %v", player, server.Config.ServerID, err)
					} else {
						server.OnlinePlayers.Mu.Lock()
						for i, p := range server.OnlinePlayers.Players {
							if player == p {
								players := server.OnlinePlayers.Players
								if i == len(players)-1 {
									server.OnlinePlayers.Players = players[:i]
								} else {
									server.OnlinePlayers.Players = append(players[:i], players[i+1:]...)
								}
							}
						}
						server.OnlinePlayers.Mu.Unlock()
					}
				}
			}
		}
	}
	if conf.Config.Whitelist.Roles.RemoveUserWithout {
		//entries, found := mongodb.Read(whitelistCollection, bson.M{
		//	"dcUserID":  bson.M{"$exists": true},
		//	"mcAccount": bson.M{"$exists": true},
		//})

		entries := database.DB.AllWhitelists()

		session := bot.Session
		var removedIDs []string
		for _, entry := range entries {
			userID := fmt.Sprintf("%v", entry.UserID)

			checkID := true
			for _, removeID := range removedIDs {
				if removeID == userID {
					checkID = false
				}
			}
			if checkID {
				user, _ := session.GuildMember(conf.Config.Discord.ServerID, userID)
				if user == nil {

					removedIDs = append(removedIDs, userID)
					database.DB.MoveToReWhitelist(entry.UserID, database.Role("discord"))

				} else {
					serverPerms := false
					for _, role := range user.Roles {
						for _, neededRole := range conf.Config.Whitelist.Roles.ServerRoleID {
							if database.Role(role) == neededRole {
								serverPerms = true
								break
							}
						}
					}
					if serverPerms == false {
						removedIDs = append(removedIDs, userID)

						database.DB.MoveToReWhitelist(entry.UserID, database.Role("server")) //TODO: get user roles here
					}
				}
			}
		}
		if len(removedIDs) > 0 {
			log.Printf("Removing accounts of the id(s) %v from whitelist because they have not the server role", removedIDs)
		}
	}
	if conf.Config.Whitelist.Roles.ReWhitelistWith {
		entries := database.DB.AllReWhitelists()

		session := bot.Session
		var addedIDs []string
		for _, entry := range entries {
			userID := fmt.Sprintf("%v", entry.UserID)

			checkID := true
			for _, addID := range addedIDs {
				if addID == userID {
					checkID = false
				}
			}
			if checkID {
				user, _ := session.GuildMember(conf.Config.Discord.ServerID, userID)
				if user != nil {
					serverPerms := false
					for _, role := range user.Roles {
						for _, neededRole := range conf.Config.Whitelist.Roles.ServerRoleID {
							if database.Role(role) == neededRole {
								serverPerms = true
								break
							}
						}
					}
					if serverPerms == true {

						addedIDs = append(addedIDs, userID)

						database.DB.ReWhitelist(entry.UserID, []database.Role{}) //TODO: get user roles here
					}
				}
			}
		}
		if len(addedIDs) > 0 {
			log.Printf("Adding accounts of the id(s) %v to whitelist because they have the server role again", addedIDs)
		}
	}

	//TODO: We need to check here if we need to execute role obtain & loose commands
}

package whitelist

import (
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/discord/bot/auth"
	"github.com/Sharktheone/ScharschBot/discord/embed/banEmbed"
	"github.com/Sharktheone/ScharschBot/discord/session"
	"github.com/Sharktheone/ScharschBot/types"
	"github.com/Sharktheone/ScharschBot/whitelist"
	"github.com/Sharktheone/ScharschBot/whitelist/whitelist/utils"
	"github.com/bwmarrin/discordgo"
	"log"
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
	mcBan, dcBan, reason := utils.CheckBanned(player, member.ID)

	if mcBan && dcBan {
		return BothBanned, reason
	}
	if mcBan {
		return McBanned, reason
	}
	if dcBan {
		return DcBanned, reason
	}

	if !utils.CheckRoles(member, conf.Config.Whitelist.Roles.ServerRoleID) {
		return NotAllowed, ""
	}

	if !utils.HasFreeAccount(member) {
		return NoFreeAccount, ""
	}

	if !utils.AccountExists(player) {
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
		if !utils.CheckRoles(member, conf.Config.Whitelist.Roles.ServerRoleID) {
			return false, false
		}
	}

	whitelist.Provider.UnWhitelistPlayer(username, member)

	log.Printf("%v is removing %v from whitelist", member.ID, username)
	return true, true
}

func RemoveAll(member *types.Member) (allowed bool, onWhitelist bool) {
	if !utils.CheckRoles(member, conf.Config.Discord.WhitelistRemoveRoleID) {
		return false, false
	}

	entries := database.DB.GetAllWhitelistedPlayers()

	log.Printf("%v is removing all accounts from whitelist", member.ID)
	for _, entry := range entries {
		database.DB.RemoveAccount(entry.Name)

		m, err := types.MemberFromID(entry.ID, auth.Session)

		if err != nil {
			log.Printf("Error while getting member: %v", err)

			whitelist.Provider.UnWhitelistPlayer(entry.Name, member)

			continue
		}

		whitelist.Provider.UnWhitelistPlayer(entry.Name, m)
	}

	return true, len(entries) > 0
}
func RemoveAllAllowed(member *types.Member) (allowed bool) {
	return utils.CheckRoles(member, conf.Config.Discord.WhitelistRemoveRoleID)

}

func Whois(username database.Player, member *types.Member) (dcUserID database.UserID, allowed bool, found bool) {
	if !utils.CheckRoles(member, conf.Config.Discord.WhitelistWhoisRoleID) {
		return "", false, false
	}

	log.Printf("%v is looking who whitelisted %v ", member.ID, username)
	result, dataFound := database.DB.GetWhitelistedPlayer(username)

	return result.ID, true, dataFound
}

func HasListed(lookupID database.UserID, member *types.Member, isSelfLookup bool) (accounts []database.Player, allowed bool, found bool, bannedPlayers []database.Player) {
	if !isSelfLookup && !utils.CheckRoles(member, conf.Config.Discord.WhitelistWhoisRoleID) {
		return nil, false, false, nil
	}
	log.Printf("%v is looking on whitelisted accounts of %v ", member.ID, lookupID)

	players := database.DB.Players(lookupID)
	return players, true, len(players) > 0, utils.CheckBans(member.ID)
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

	if !utils.CheckRoles(member, conf.Config.Discord.WhitelistBanRoleID) {
		return false, false
	}

	listedAccounts := ListedAccountsOf(banID, false)

	_, banned, _ := utils.CheckBanned("", banID)
	if banned {
		return true, true
	}

	log.Printf("%v is banning %v", member.ID, banID)

	database.DB.BanUser(banID, reason)

	if banAccounts {
		banMember, err := types.MemberFromID(banID, s)
		if err != nil {
			log.Printf("Error while getting member: %v", err)
		} else {
			for _, account := range listedAccounts {
				database.DB.UnWhitelistPlayer(account)

				whitelist.Provider.BanPlayer(account, banMember, reason)

			}
		}

	}

	database.DB.BanUser(banID, reason)
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
	return true, false
}

func BanAccount(member *types.Member, account database.Player, reason string, s *session.Session) (bool, *Player) {
	if !utils.CheckRoles(member, conf.Config.Discord.WhitelistBanRoleID) {
		return false, nil
	}

	owner := GetOwner(account, s)
	if owner.Whitelisted {
		alreadyBanned := database.DB.IsPlayerBanned(account)

		if !alreadyBanned {
			log.Printf("%v is banning %v", owner.ID, account)
			messageEmbedDM := banEmbed.DMBanAccount(string(account), false, owner.ID, reason, s)
			messageEmbedDMFailed := banEmbed.DMBanAccount(string(account), true, owner.ID, reason, s)
			if err := s.SendDM(owner.ID, &discordgo.MessageSend{
				Embed: &messageEmbedDM,
			}, &discordgo.MessageSend{
				Content: fmt.Sprintf("<@%v>", owner.ID),
				Embed:   &messageEmbedDMFailed,
			},
			); err != nil {
				log.Printf("Failed to send DM to %v: %v", owner.ID, err)
			}

			banMember, err := types.MemberFromID(owner.ID, s)
			if err != nil {
				log.Printf("Error while getting member: %v", err)
			} else {
				whitelist.Provider.BanPlayer(account, banMember, reason)
			}
		}
	} else {
		return false, nil
	}

	return true, owner
}

// UnBanUserID  unbans a user from the whitelist
// / returns true if the action was allowed
func UnBanUserID(member *types.Member, banID database.UserID, unbanAccounts bool, s *session.Session) bool {
	if !utils.CheckRoles(member, conf.Config.Discord.WhitelistBanRoleID) {
		return false
	}

	log.Printf("%v is unbanning %v", member.ID, banID)
	database.DB.UnBanUser(banID)
	if unbanAccounts {
		result := database.DB.BannedPlayers(banID)

		for _, entry := range result {
			database.DB.UnBanPlayerFrom(banID, entry.Player)

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

	return true
}

func UnBanAccount(member *types.Member, account database.Player, s *session.Session) bool {
	if !utils.CheckRoles(member, conf.Config.Discord.WhitelistBanRoleID) {
		return false
	}

	log.Printf("%v is unbanning %v", member.ID, account)
	database.DB.UnBanPlayer(account)
	messageEmbedDM := banEmbed.DMUnBanAccount(string(account), false, member.ID, s)
	messageEmbedDMFailed := banEmbed.DMUnBanAccount(string(account), true, member.ID, s)
	if err := s.SendDM(member.ID, &discordgo.MessageSend{
		Embed: &messageEmbedDM,
	}, &discordgo.MessageSend{
		Content: fmt.Sprintf("<@%v>", member.ID),
		Embed:   &messageEmbedDMFailed,
	},
	); err != nil {
		log.Printf("Failed to send DM to %v: %v", member.ID, err)
	}

	return true
}

func RemoveMyAccounts(userID database.UserID) *[]database.Player {
	log.Printf("%v is removing his own accounts from the whitelist", userID)

	member, err := types.MemberFromID(userID, auth.Session)
	if err != nil {
		log.Printf("Error while getting member: %v", err)
		return nil
	}

	entries := database.DB.Players(userID)

	whitelist.Provider.UnWhitelistAccount(member)

	return &entries

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
			rls, err := s.GetRoles(dcUser)
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
			BannedPlayers:     utils.CheckBans(dcUser),
			Roles:             roles,
			MaxAccounts:       utils.GetMaxAccounts(&member),
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

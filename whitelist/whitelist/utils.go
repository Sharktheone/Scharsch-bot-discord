package whitelist

import (
	"fmt"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/types"
	"io"
	"log"
	"net/http"
)

func HasFreeAccount(player database.Player, member *types.Member) bool {
	num := database.DB.NumberWhitelistedPlayers(member.ID)

	if GetMaxAccounts(member) > (num + len(CheckBans(member.ID))) {
		return true
	}

	return false
}

func CheckRoles(member *types.Member, required []database.Role) bool {
	for _, role := range required {
		if CheckRole(member, role) {
			return true
		}
	}

	return false
}

func CheckRole(member *types.Member, required database.Role) bool {
	for _, role := range member.Roles {
		if role == required {
			return true
		}
	}
	return false
}

func CheckBanned(player database.Player, userID database.UserID) (mcBanned bool, dcBanned bool, banReason string) {
	var (
		reason string
	)

	mcReason, mc := database.DB.GetPlayerBan(player)

	dcReason, dc := database.DB.GetBan(userID)

	if mc {
		reason = fmt.Sprintf("%v", mcReason)
	}
	if dc {
		reason += fmt.Sprintf("%v", dcReason)
	}

	return mc, dc, reason
}

func CheckBans(userID database.UserID) []database.Player {
	results := database.DB.BannedPlayers(userID)

	var bannedAccounts = make([]database.Player, len(results))
	for i, result := range results {
		bannedAccounts[i] = result.Player

	}
	return bannedAccounts
}

func AccountExists(username database.Player) bool {
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

func GetMaxAccounts(member *types.Member) int {
	m := 0

	for _, entry := range config.Whitelist.MaxAccounts {
		if CheckRole(member, entry.RoleID) {
			if entry.Max > m {
				m = entry.Max
			}
		}
	}
	return m
}

package gorm

import (
	"github.com/Sharktheone/ScharschBot/database"
	"gorm.io/gorm"
	"log"
	"slices"
)

type GormConnection struct {
	DB *gorm.DB
}

func (m *GormConnection) Connect() {
	var err error
	m.DB, err = gorm.Open(GetDialector())
	if err != nil {
		log.Panicf("Failed to connect to database: %v", err)
	}

	if err := m.DB.AutoMigrate(&WhitelistEntry{}, &BanEntry{}, &ReWhitelistEntry{}, &ReportEntry{}, &WaitlistEntry{}); err != nil {
		log.Panicf("Failed to migrate database: %v", err)
	}
}

func (m *GormConnection) Disconnect() {
	sqlDB, err := m.DB.DB()
	if err != nil {
		log.Panicf("Failed to get database: %v", err)
	}

	if err := sqlDB.Close(); err != nil {
		log.Panicf("Failed to close database: %v", err)
	}
}

func (m *GormConnection) WhitelistPlayer(user database.UserID, player database.Player) {
	if err := m.DB.Create(&WhitelistEntry{
		UserID: user,
		Player: player,
	}).Error; err != nil {
		log.Printf("Failed to whitelist player: %v", err)
	}
}

func (m *GormConnection) UnWhitelistAccount(user database.UserID) {
	if err := m.DB.Where(WhitelistEntry{UserID: user}).Delete(&WhitelistEntry{}).Error; err != nil {
		log.Printf("Failed to unwhitelist account: %v", err)
	}
}

func (m *GormConnection) UnWhitelistPlayer(player database.Player) {
	if err := m.DB.Where(WhitelistEntry{Player: player}).Delete(&WhitelistEntry{}).Error; err != nil {
		log.Printf("Failed to unwhitelist player: %v", err)
	}
}

func (m *GormConnection) MoveToReWhitelist(user database.UserID, missingRole database.Role) {
	if err := m.DB.Create(&ReWhitelistEntry{
		UserID:      user,
		Players:     m.Players(user),
		MissingRole: missingRole,
	}).Error; err != nil {
		log.Printf("Failed to move to rewhitelist: %v", err)
	}

	if err := m.DB.Where(WhitelistEntry{UserID: user}).Delete(&WhitelistEntry{}).Error; err != nil {
		log.Printf("Failed to unwhitelist account: %v", err)
	}
}

func (m *GormConnection) ReWhitelist(user database.UserID, roles []database.Role) {
	var entry ReWhitelistEntry
	if err := m.DB.Where(ReWhitelistEntry{UserID: user}).First(&entry).Error; err != nil {
		log.Printf("Failed to rewhitelist: %v", err)
	}

	if slices.Contains(roles, entry.MissingRole) {
		for _, player := range entry.Players {
			m.WhitelistPlayer(user, player)
		}
	}
}

func (m *GormConnection) RemoveAll() {
	if err := m.DB.Delete(&WhitelistEntry{}).Error; err != nil {
		log.Printf("Failed to remove all: %v", err)
	}
}

func (m *GormConnection) RemoveAllFrom(user database.UserID) {
	if err := m.DB.Where(WhitelistEntry{UserID: user}).Delete(&WhitelistEntry{}).Error; err != nil {
		log.Printf("Failed to remove all from: %v", err)
	}
}

func (m *GormConnection) Owner(player database.Player) database.UserID {
	var entry WhitelistEntry

	if err := m.DB.Where(WhitelistEntry{Player: player}).Find(&entry).Error; err != nil {
		log.Printf("Failed to get owner: %v", err)
		return "<unknown>"
	}

	return entry.UserID
}

func (m *GormConnection) Players(user database.UserID) []database.Player {
	var entries []WhitelistEntry
	if err := m.DB.Where(WhitelistEntry{UserID: user}).Find(&entries).Error; err != nil {
		log.Printf("Failed to get players: %v", err)
		return []database.Player{}
	}

	players := make([]database.Player, len(entries))
	for i, entry := range entries {
		players[i] = entry.Player
	}

	return players
}

func (m *GormConnection) BanUser(user database.UserID, reason string) {
	if err := m.DB.Create(&BanEntry{
		UserID:     user,
		UserBan:    true,
		UserReason: reason,
	}).Error; err != nil {
		log.Printf("Failed to ban user: %v", err)
	}

	if err := m.DB.Where(WhitelistEntry{UserID: user}).Delete(&WhitelistEntry{}).Error; err != nil {
		log.Printf("Failed to unwhitelist account: %v", err)
	}
}

func (m *GormConnection) BanPlayer(player database.Player, reason string) {
	//TODO: Ban methods should first query if there already is an ban entry for that user and add it to the existing one
	if err := m.DB.Create(&BanEntry{
		Players: []database.PlayerBanData{{Player: player, Reason: reason}},
	}).Error; err != nil {
		log.Printf("Failed to ban player: %v", err)
	}

	if err := m.DB.Where(WhitelistEntry{Player: player}).Delete(&WhitelistEntry{}).Error; err != nil {
		log.Printf("Failed to unwhitelist player: %v", err)
	}
}

func (m *GormConnection) UnBanUser(user database.UserID) {
	if err := m.DB.Where(BanEntry{UserID: user}).Delete(&BanEntry{}).Error; err != nil {
		log.Printf("Failed to unban user: %v", err)
	}

}

func (m *GormConnection) UnBanPlayer(player database.Player) {
	var entry BanEntry

	data := m.DB.Where(BanEntry{Players: []database.PlayerBanData{{Player: player}}})
	if data.Error != nil {
		log.Printf("Failed to unban player: %v", data.Error)
		return
	}

	var count int64
	data.Count(&count)

	if count == 0 {
		return
	}

	if err := data.First(&entry).Error; err != nil {
		log.Printf("Failed to unban player: %v", err)
	}

	if len(entry.Players) == 1 {
		if err := m.DB.Where(BanEntry{Players: []database.PlayerBanData{{Player: player}}}).Delete(&BanEntry{}).Error; err != nil {
			log.Printf("Failed to unban player: %v", err)
		}
	} else {
		for i, playerBan := range entry.Players {
			if playerBan.Player == player {
				entry.Players = append(entry.Players[:i], entry.Players[i+1:]...)
				break
			}
		}

		if err := m.DB.Save(&entry).Error; err != nil {
			log.Printf("Failed to unban player: %v", err)
		}
	}

}

func (m *GormConnection) UnBanPlayerFrom(user database.UserID, player database.Player) {
	var entry BanEntry
	if err := m.DB.Where(BanEntry{UserID: user}).First(&entry).Error; err != nil {
		log.Printf("Failed to unban player from: %v", err)
	}

	for i, playerBan := range entry.Players {
		if playerBan.Player == player {
			entry.Players = append(entry.Players[:i], entry.Players[i+1:]...)
			break
		}
	}

	if err := m.DB.Save(&entry).Error; err != nil {
		log.Printf("Failed to unban player from: %v", err)
	}
}

func (m *GormConnection) IsUserBanned(user database.UserID) bool {
	items := m.DB.Where(BanEntry{UserID: user})

	var count int64

	if err := items.Count(&count).Error; err != nil {
		log.Printf("Failed to check if user is banned: %v", err)
		return true // return true to prevent user from being able to whitelist
	}

	return count > 0
}

func (m *GormConnection) IsPlayerBanned(player database.Player) bool {
	var entry BanEntry

	data := m.DB.Where(BanEntry{Players: []database.PlayerBanData{{Player: player}}})
	if data.Error != nil {
		log.Printf("Failed to check if player is banned: %v", data.Error)
		return true
	}

	var count int64
	data.Count(&count)

	if count == 0 {
		return false
	}

	if err := data.First(&entry).Error; err != nil {
		return false
	}

	return true
}

func (m *GormConnection) BannedPlayers(user database.UserID) []database.PlayerBanData {
	var entry BanEntry

	data := m.DB.Where(BanEntry{UserID: user})
	if data.Error != nil {
		log.Printf("Failed to get banned players: %v", data.Error)
		return nil
	}

	var count int64
	data.Count(&count)

	if count == 0 {
		return []database.PlayerBanData{}
	}

	if err := data.First(&entry).Error; err != nil {
		return nil
	}

	return entry.Players
}

func (m *GormConnection) RemoveAccounts(user database.UserID) {
	if err := m.DB.Where(WhitelistEntry{UserID: user}).Delete(&WhitelistEntry{}).Error; err != nil {
		log.Printf("Failed to remove accounts: %v", err)
	}
}

func (m *GormConnection) AddWaitlist(user database.UserID, player database.Player) {
	if err := m.DB.Create(&WaitlistEntry{
		UserID: user,
		Player: player,
	}).Error; err != nil {
		log.Printf("Failed to add to waitlist: %v", err)
	}
}

func (m *GormConnection) RemoveWaitlist(user database.UserID, player database.Player) {
	if err := m.DB.Where(WaitlistEntry{UserID: user, Player: player}).Delete(&WaitlistEntry{}).Error; err != nil {
		log.Printf("Failed to remove from waitlist: %v", err)
	}
}

func (m *GormConnection) Report(reporter database.UserID, reported database.Player, reason string) {
	if err := m.DB.Create(&ReportEntry{
		ReporterID:     reporter,
		ReportedPlayer: reported,
		Reason:         reason,
	}).Error; err != nil {
		log.Printf("Failed to report: %v", err)
	}
}

func (m *GormConnection) AllWhitelists() []database.WhitelistEntry {
	var entries []WhitelistEntry
	if err := m.DB.Find(&entries).Error; err != nil {
		log.Printf("Failed to get all whitelists: %v", err)
	}

	whitelists := make([]database.WhitelistEntry, len(entries))
	for i, entry := range entries {
		whitelists[i] = database.WhitelistEntry{
			UserID: entry.UserID,
			Player: entry.Player,
		}
	}

	return whitelists
}

func (m *GormConnection) AllReWhitelists() []database.ReWhitelistEntry {
	var entries []ReWhitelistEntry
	if err := m.DB.Find(&entries).Error; err != nil {
		log.Printf("Failed to get all rewhitelists: %v", err)
	}

	reWhitelists := make([]database.ReWhitelistEntry, len(entries))
	for i, entry := range entries {
		reWhitelists[i] = database.ReWhitelistEntry{
			UserID:      entry.UserID,
			Players:     entry.Players,
			MissingRole: entry.MissingRole,
		}
	}

	return reWhitelists
}

func (m *GormConnection) DeleteReport(reported database.Player) {
	if err := m.DB.Where(ReportEntry{ReportedPlayer: reported}).Delete(&ReportEntry{}).Error; err != nil {
		log.Printf("Failed to delete report: %v", err)
	}
}

func (m *GormConnection) RemoveAccount(player database.Player) {
	if err := m.DB.Where(WhitelistEntry{Player: player}).Delete(&WhitelistEntry{}).Error; err != nil {
		log.Printf("Failed to remove account: %v", err)
	}
}

func (m *GormConnection) IsWhitelisted(player database.Player) bool {
	var entry WhitelistEntry

	data := m.DB.Where(WhitelistEntry{Player: player})

	if data.Error != nil {
		log.Printf("Failed to check if player is whitelisted: %v", data.Error)
		return false
	}

	var count int64
	data.Count(&count)

	if count == 0 {
		return false
	}

	if err := data.First(&entry).Error; err != nil {
		return false
	}

	return true
}

func (m *GormConnection) IsWhitelistedBy(user database.UserID, player database.Player) bool {
	var entry WhitelistEntry

	data := m.DB.Where(WhitelistEntry{UserID: user, Player: player})

	if data.Error != nil {
		log.Printf("Failed to check if player is whitelisted by user: %v", data.Error)
		return false
	}

	var count int64
	data.Count(&count)

	if count == 0 {
		return false
	}

	if err := data.First(&entry).Error; err != nil {
		return false
	}

	return true
}

func (m *GormConnection) GetReports() []database.ReportData {
	var entries []ReportEntry
	if err := m.DB.Find(&entries).Error; err != nil {
		log.Printf("Failed to get reports: %v", err)
	}

	reports := make([]database.ReportData, len(entries))
	for i, entry := range entries {
		reports[i] = database.ReportData{
			ReporterID:     entry.ReporterID,
			ReportedPlayer: entry.ReportedPlayer,
			Reason:         entry.Reason,
		}
	}

	return reports
}

func (m *GormConnection) IsAlreadyReported(reported database.Player) bool {
	var entry ReportEntry

	data := m.DB.Where(ReportEntry{ReportedPlayer: reported})

	if data.Error != nil {
		log.Printf("Failed to check if player is already reported: %v", data.Error)
		return false
	}

	var count int64
	data.Count(&count)

	if count == 0 {
		return false
	}

	if err := data.First(&entry).Error; err != nil {
		return false
	}

	return true
}

func (m *GormConnection) GetReportedPlayer(reported database.Player) (database.ReportData, bool) {
	var entry ReportEntry

	data := m.DB.Where(ReportEntry{ReportedPlayer: reported})

	if data.Error != nil {
		log.Printf("Failed to get reported player: %v", data.Error)
		return database.ReportData{}, false
	}

	if err := data.First(&entry).Error; err != nil {
		return database.ReportData{}, false
	}

	if err := data.First(&entry).Error; err != nil {
		return database.ReportData{}, false
	}

	return database.ReportData{
		ReporterID:     entry.ReporterID,
		ReportedPlayer: entry.ReportedPlayer,
		Reason:         entry.Reason,
	}, true
}

func (m *GormConnection) NumberWhitelistedPlayers(user database.UserID) int {
	return len(m.Players(user))
}

func (m *GormConnection) GetWhitelistedPlayer(player database.Player) (database.WhitelistedPlayerData, bool) {
	var entry WhitelistEntry

	if err := m.DB.Where(WhitelistEntry{Player: player}).Find(&entry).Error; err != nil {
		log.Printf("Failed to get whitelisted player: %v", err)
		return database.WhitelistedPlayerData{}, false
	}

	return database.WhitelistedPlayerData{
		ID:   entry.UserID,
		Name: entry.Player,
	}, true
}

func (m *GormConnection) GetAllWhitelistedPlayers() []database.WhitelistedPlayerData {
	var entries []WhitelistEntry
	if err := m.DB.Find(&entries).Error; err != nil {
		log.Printf("Failed to get all whitelisted players: %v", err)
	}

	whitelistedPlayers := make([]database.WhitelistedPlayerData, len(entries))
	for i, entry := range entries {
		whitelistedPlayers[i] = database.WhitelistedPlayerData{
			ID:   entry.UserID,
			Name: entry.Player,
		}
	}

	return whitelistedPlayers
}

func (m *GormConnection) GetAccountsOf(user database.UserID) []database.Player {
	return m.Players(user)
}

func (m *GormConnection) GetBan(user database.UserID) (string, bool) {
	var entry BanEntry

	data := m.DB.Where(BanEntry{UserID: user})

	if data.Error != nil {
		log.Printf("Failed to get ban: %v", data.Error)
		return "", false
	}

	var count int64
	data.Count(&count)

	if count == 0 {
		return "", false
	}

	if err := data.First(&entry).Error; err != nil {
		return "", false
	}

	return entry.UserReason, true
}

func (m *GormConnection) GetPlayerBan(player database.Player) (database.PlayerBan, bool) {
	var entry BanEntry

	data := m.DB.Where(BanEntry{Players: []database.PlayerBanData{{Player: player}}})

	if data.Error != nil {
		log.Printf("Failed to get player ban: %v", data.Error)
		return database.PlayerBan{}, false
	}

	var count int64
	data.Count(&count)

	if count == 0 {
		return database.PlayerBan{}, false
	}

	if err := data.First(&entry).Error; err != nil {
		return database.PlayerBan{}, false
	}

	return database.PlayerBan{
		ID:     entry.UserID,
		Reason: entry.UserReason,
	}, true
}

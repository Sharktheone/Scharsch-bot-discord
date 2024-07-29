package gorm

import (
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database"
)

var (
	config = conf.Config.Whitelist.Database
)

type WhitelistEntry struct {
	ID     uint            `json:"-" bson:"-" gorm:"primaryKey"`
	UserID database.UserID `json:"userID" bson:"userID" gorm:"userID"`
	Player database.Player `json:"player" bson:"player" gorm:"player"`
}

func (m *WhitelistEntry) TableName() string {
	return config.WhitelistTable
}

type BanEntry struct {
	ID         uint            `json:"-" bson:"-" gorm:"primaryKey"`
	UserID     database.UserID `json:"userID" bson:"userID" gorm:"userID"`
	UserBan    bool            `json:"userBan" bson:"userBan" gorm:"userBan"`
	Players    []PlayerBanData `json:"players" bson:"players" gorm:"players"`
	UserReason string          `json:"reason" bson:"reason" gorm:"userReason"`
}

func (m *BanEntry) TableName() string {
	return config.BanTable
}

type PlayerBanData struct {
	Player database.Player `json:"player" bson:"player" gorm:"player"`
	Reason string          `json:"reason" bson:"reason" gorm:"reason"`
}

type ReWhitelistEntry struct {
	ID          uint              `json:"-" bson:"-" gorm:"primaryKey"`
	UserID      database.UserID   `json:"userID" bson:"userID" gorm:"userID"`
	Players     []database.Player `json:"players" bson:"players" gorm:"players"`
	MissingRole database.Role     `json:"missingRole" bson:"missingRole" gorm:"missingRole"`
}

func (m *ReWhitelistEntry) TableName() string {
	return config.ReWhitelistTable
}

type ReportEntry struct {
	ID             uint            `json:"-" bson:"-" gorm:"primaryKey"`
	ReporterID     database.UserID `json:"reporterID" bson:"reporterID" gorm:"reporterID"`
	ReportedPlayer database.Player `json:"reportedPlayer" bson:"reportedPlayer" gorm:"reportedPlayer"`
	Reason         string          `json:"reason" bson:"reason"`
}

func (m *ReportEntry) TableName() string {
	return config.ReportTable
}

type WaitlistEntry struct {
	ID     uint            `json:"-" bson:"-" gorm:"primaryKey"`
	UserID database.UserID `json:"userID" bson:"userID" gorm:"userID"`
	Player database.Player `json:"player" bson:"player" gorm:"player"`
}

func (m *WaitlistEntry) TableName() string {
	return config.WaitListTable
}

package gorm

import (
	"bytes"
	"database/sql/driver"
	"fmt"
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
	ID         uint                     `json:"-" bson:"-" gorm:"primaryKey"`
	UserID     database.UserID          `json:"userID" bson:"userID" gorm:"userID"`
	UserBan    bool                     `json:"userBan" bson:"userBan" gorm:"userBan"`
	Players    []database.PlayerBanData `json:"players" bson:"players" gorm:"foreignKey:ID"`
	UserReason string                   `json:"reason" bson:"reason" gorm:"userReason"`
}

func (m *BanEntry) TableName() string {
	return config.BanTable
}

type ReWhitelistEntry struct {
	ID          uint            `json:"-" bson:"-" gorm:"primaryKey"`
	UserID      database.UserID `json:"userID" bson:"userID" gorm:"userID"`
	Players     PlayerList      `json:"players" bson:"players" gorm:"type:text[];players"`
	MissingRole database.Role   `json:"missingRole" bson:"missingRole" gorm:"missingRole"`
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

type PlayerList []database.Player

func (pl PlayerList) Value() (driver.Value, error) {
	if len(pl) == 0 {
		return nil, nil
	}

	var strbuf bytes.Buffer

	for i, p := range pl {
		if i != 0 {
			strbuf.WriteString(",")
		}
		strbuf.WriteString(string(p))
	}

	return strbuf.String(), nil
}

func (pl *PlayerList) Scan(value interface{}) error {
	str, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("PlayerList.Scan: expected string, got %T", value)
	}

	players := bytes.Split(str, []byte(","))
	*pl = make([]database.Player, len(players))

	for i, p := range players {
		(*pl)[i] = database.Player(p)
	}

	return nil
}

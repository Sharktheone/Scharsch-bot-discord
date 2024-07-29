package database

type WhitelistEntry struct {
	UserID UserID `json:"userID" bson:"userID" gorm:"userID"`
	Player Player `json:"player" bson:"player" gorm:"player"`
}

type BanEntry struct {
	UserID     UserID          `json:"userID" bson:"userID" gorm:"userID"`
	UserBan    bool            `json:"userBan" bson:"userBan" gorm:"userBan"`
	Players    []PlayerBanData `json:"players" bson:"players" gorm:"players"`
	UserReason string          `json:"reason" bson:"reason" gorm:"userReason"`
}

type PlayerBanData struct {
	Player Player `json:"player" bson:"player" gorm:"player"`
	Reason string `json:"reason" bson:"reason" gorm:"reason"`
}

type ReWhitelistEntry struct {
	UserID      UserID   `json:"userID" bson:"userID" gorm:"userID"`
	Players     []Player `json:"players" bson:"players" gorm:"players"`
	MissingRole Role     `json:"missingRole" bson:"missingRole" gorm:"missingRole"`
}

type ReportEntry struct {
	ReporterID     UserID `json:"reporterID" bson:"reporterID" gorm:"reporterID"`
	ReportedPlayer Player `json:"reportedPlayer" bson:"reportedPlayer" gorm:"reportedPlayer"`
	Reason         string `json:"reason" bson:"reason"`
}

type WaitlistEntry struct {
	UserID UserID `json:"userID" bson:"userID" gorm:"userID"`
	Player Player `json:"player" bson:"player" gorm:"player"`
}

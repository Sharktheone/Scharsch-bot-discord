package database

type WhitelistEntry struct {
	UserID UserID `json:"userID" bson:"userID"`
	Player Player `json:"player" bson:"player"`
}

type BanEntry struct {
	UserID  UserID   `json:"userID" bson:"userID"`
	Players []Player `json:"players" bson:"players"`
}

type ReWhitelsitEntry struct {
	UserID      UserID   `json:"userID" bson:"userID"`
	Players     []Player `json:"players" bson:"players"`
	MissingRole Role     `json:"missingRole" bson:"missingRole"`
}

type ReportEntry struct {
	ReporterID     UserID `json:"reporterID" bson:"reporterID"`
	ReportedPlayer Player `json:"reportedPlayer" bson:"reportedPlayer"`
	Reason         string `json:"reason" bson:"reason"`
}

type WaitlistEntry struct {
	UserID UserID `json:"userID" bson:"userID"`
	Player Player `json:"player" bson:"player"`
}

package auth

import (
	"github.com/Sharktheone/ScharschBot/discord/session"
	"github.com/Sharktheone/ScharschBot/flags"
)

var (
	GuildID = flags.String("guild")
	Session *session.Session
)

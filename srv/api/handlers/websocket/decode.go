package websocket

import (
	"fmt"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/discord/bot"
	"github.com/Sharktheone/ScharschBot/srv/playersrv"
	"github.com/Sharktheone/ScharschBot/types"
	"github.com/Sharktheone/ScharschBot/whitelist/whitelist"
	"strings"
)

var (
	s = bot.Session
)

func (h *Handler) DecodePlayer(e *types.WebsocketEvent) (*PSRVEvent, error) {
	var (
		errMsg error
		pSrv   = &PSRVEvent{
			h:       h,
			e:       e,
			session: s,
		}
	)
	if e.Data.Player == "" {
		return pSrv, nil
	}
	if owner := whitelist.GetOwner(e.Data.Player, s); owner.Whitelisted {
		pSrv.onWhitelist = &owner.Whitelisted
		pSrv.userID = &owner.ID
		if member, err := s.GetUserProfile(owner.ID); err != nil {
			errMsg = fmt.Errorf("failed to get user profile: %v", err)
			pSrv.footerIcon = &config.Discord.EmbedErrorIcon
			pSrv.username = &e.Data.Player
			pSrv.member = member
		} else {
			a := member.User.AvatarURL("40")
			u := member.User.String()
			pSrv.footerIcon = &a
			pSrv.username = (*database.Player)(&u)
			pSrv.member = member
		}
	} else {
		pSrv.onWhitelist = &owner.Whitelisted
		pSrv.footerIcon = &config.Discord.EmbedErrorIcon
	}
	playersrv.CheckAccount(database.Player(strings.ToLower(string(e.Data.Player))))

	return pSrv, errMsg
}

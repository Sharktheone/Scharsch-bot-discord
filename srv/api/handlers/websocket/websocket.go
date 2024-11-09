package websocket

import (
	"context"
	"fmt"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/discord/session"
	"github.com/Sharktheone/ScharschBot/pterodactyl"
	"github.com/Sharktheone/ScharschBot/types"
	"github.com/bwmarrin/discordgo"
	"github.com/fasthttp/websocket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// Bot => Server: SendPlayers
// Bot => Server: KickPlayer
// Bot => Server: ReportPlayer
// Bot => Server: BanPlayer
// Bot => Server: UnbanPlayer
// Bot => Server: SendCommand
// Bot => Server: SendChatMessage
// Bot => Server: AuthSuccess
// Bot => Server: AuthFailed
// Bot => Server: Error
// Bot => Server: WhitelistedPlayers
//
// Server => Bot: PlayerJoined
// Server => Bot: PlayerLeft
// Server => Bot: Players
// Server => Bot: ChatMessage
// Server => Bot: PlayerDeath
// Server => Bot: PlayerAdvancement
// Server => Bot: Auth
// Server => Bot: Error
const (
	SendPlayers        = "sendPlayers"
	KickPlayer         = "kickPlayer"
	ReportPlayer       = "reportPlayer"
	BanPlayer          = "banPlayer"
	UnbanPlayer        = "unbanPlayer"
	SendCommand        = "sendCommand"
	SendChatMessage    = "sendChatMessage"
	WhitelistAdd       = "whitelistAdd"
	WhitelistRemove    = "whitelistRemove"
	WhitelistedPlayers = "whitelistedPlayers"
	PlayerJoined       = "playerJoined"
	PlayerLeft         = "playerLeft"
	Players            = "players"
	ChatMessage        = "chatMessage"
	PlayerDeath        = "playerDeath"
	PlayerAdvancement  = "playerAdvancement"
	Auth               = "auth"
	AuthSuccess        = "authSuccess"
	AuthFailed         = "authFailed"
	Error              = "error"
)

type Handler struct {
	conn          *websocket.Conn
	server        *pterodactyl.Server
	uuid          uuid.UUID
	send          chan *types.WebsocketEvent
	receive       chan *types.WebsocketEvent
	ctx           context.Context
	cancel        context.CancelFunc
	authenticated bool
}

type PSRVEvent struct {
	h           *Handler
	e           *types.WebsocketEvent
	userID      *database.UserID
	onWhitelist *bool
	footerIcon  *string
	username    *database.Player
	member      *discordgo.Member
	session     *session.Session
}

var (
	upgrader = websocket.Upgrader{}
)

func ServerWS(c *gin.Context) {
	s, err := pterodactyl.GetServer(c.Param("serverID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid server id %v: %v", c.Param("serverID"), err),
		})
		return
	}
	handler, err := getWSHandler(s, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to get server: %v", err),
		})
		return
	}

	go handler.handleInbound()
	go handler.handleOutbound()
}

func getWSHandler(s *pterodactyl.Server, c *gin.Context) (*Handler, error) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade connection: %v", err)
	}
	u, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate uuid: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Handler{
		conn:    conn,
		server:  s,
		uuid:    u,
		ctx:     ctx,
		cancel:  cancel,
		send:    make(chan *types.WebsocketEvent),
		receive: make(chan *types.WebsocketEvent),
	}, nil
}

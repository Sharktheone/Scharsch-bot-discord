package pterodactyl

import (
	"context"
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/fasthttp/websocket"
	"sync"
)

//goland:noinspection GoUnusedConst
const (
	websocketAuthSuccess   = "auth success"
	WebsocketStatus        = "status"
	WebsocketConsoleOutput = "console output"
	WebsocketStats         = "stats"
	websocketTokenExpiring = "token expiring"
	websocketTokenExpired  = "token expired"

	PowerSignalStart   = "start"
	PowerSignalStop    = "stop"
	PowerSignalKill    = "kill"
	PowerSignalRestart = "restart"

	PowerStatusRunning  = "running"
	PowerStatusOffline  = "offline"
	PowerStatusStarting = "starting"
	PowerStatusStopping = "stopping"
)

var (
	_config   = conf.Config
	_panelURL = _config.Pterodactyl.PanelURL
	_apiKey   = fmt.Sprintf("Bearer %s", _config.Pterodactyl.APIKey)
)
var (
	Servers []*Server
	mu      sync.RWMutex
)

type ServerStatus struct {
	State   string  `json:"state"`
	Ram     int     `json:"memory_bytes"`
	RamMax  int     `json:"memory_limit_bytes"`
	Cpu     float64 `json:"cpu_absolute"`
	Network struct {
		Rx int `json:"rx_bytes"`
		Tx int `json:"tx_bytes"`
	} `json:"network"`
	Disk   int `json:"disk_bytes"`
	Uptime int `json:"uptime"`
}

type ChanData struct {
	Event string
	Data  *ServerStatus
}

type listenerCtx struct {
	id     string
	cancel context.CancelFunc
	ctx    *context.Context
}

type Server struct {
	OnlinePlayers struct {
		Players []*database.Player
		Mu      sync.Mutex
	}
	Config    *conf.Server
	Data      chan *ChanData
	Console   chan string
	Status    *ServerStatus
	socket    *websocket.Conn
	connected bool
	lCtx      struct {
		ctx []*listenerCtx
		mu  sync.Mutex
	}
	ctx *context.Context
}

func New(ctx *context.Context, config *conf.Server) *Server {
	s := &Server{
		ctx:       ctx,
		Config:    config,
		Data:      make(chan *ChanData),
		Console:   make(chan string),
		Status:    &ServerStatus{},
		connected: false,
	}
	mu.Lock()
	Servers = append(Servers, s)
	mu.Unlock()

	return s
}

func (s *Server) SendCommand(command string) error {
	var (
		commandAction = []byte(fmt.Sprintf(`{"event":"set command", "args": "%s"}`, command))
	)
	return s.socket.WriteMessage(websocket.TextMessage, commandAction)
}

func (s *Server) AddListener(listener func(ctx *context.Context, server *conf.Server, data chan *ChanData), name string) {
	ctx, cancel := context.WithCancel(*s.ctx)
	s.lCtx.ctx = append(s.lCtx.ctx, &listenerCtx{
		id:     name,
		cancel: cancel,
		ctx:    &ctx,
	})
	go listener(&ctx, s.Config, s.Data)
}

func (s *Server) RemoveListener(name string) {
	s.lCtx.mu.Lock()
	defer s.lCtx.mu.Unlock()
	for i, l := range s.lCtx.ctx {
		if l.id == name || name == "*" {
			l.cancel()
			s.lCtx.ctx = append(s.lCtx.ctx[:i], s.lCtx.ctx[i+1:]...)
			return
		}
	}
}

func (s *Server) AddConsoleListener(listener func(server *conf.Server, console chan string)) {
	go listener(s.Config, s.Console)
}

func (s *Server) Start() error {
	return s.Power(PowerSignalStart)
}

func (s *Server) Stop() error {
	return s.Power(PowerSignalStop)
}

func (s *Server) Kill() error {
	return s.Power(PowerSignalKill)
}

func (s *Server) Restart() error {
	return s.Power(PowerSignalRestart)
}

func (s *Server) Power(signal string) error {
	var (
		powerAction = []byte(fmt.Sprintf(`{"event":"set state", "args": "%s"}`, signal))
	)
	return s.socket.WriteMessage(websocket.TextMessage, powerAction)
}

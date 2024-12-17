package pterodactyl

import (
	"context"
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database"
	"github.com/Sharktheone/ScharschBot/pterodactyl/types"
	"github.com/fasthttp/websocket"
	"sync"
)

var (
	Servers []*Server
	mu      sync.RWMutex
)

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
	Data      chan *types.ChanData
	Console   chan string
	Status    *types.ServerStatus
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
		Data:      make(chan *types.ChanData),
		Console:   make(chan string),
		Status:    &types.ServerStatus{},
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

func (s *Server) AddListener(listener func(ctx *context.Context, server *conf.Server, data chan *types.ChanData), name string) {
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
	return s.Power(types.PowerSignalStart)
}

func (s *Server) Stop() error {
	return s.Power(types.PowerSignalStop)
}

func (s *Server) Kill() error {
	return s.Power(types.PowerSignalKill)
}

func (s *Server) Restart() error {
	return s.Power(types.PowerSignalRestart)
}

func (s *Server) Power(signal string) error {
	var (
		powerAction = []byte(fmt.Sprintf(`{"event":"set state", "args": "%s"}`, signal))
	)
	return s.socket.WriteMessage(websocket.TextMessage, powerAction)
}

package pterodactyl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/fasthttp/websocket"
	"log"
)

type eventType struct {
	Event string   `json:"event"`
	Args  []string `json:"args"`
}

func (s *Server) connectWS() error {
	res, err := request(fmt.Sprintf("/api/client/servers/%s/websocket", s.Config.ServerID), "GET", nil)
	if err != nil {
		return fmt.Errorf("could not connect to websocket: %w", err)
	}
	if res != nil {
		if res.StatusCode == 200 {
			var socketInfo struct {
				Data struct {
					Token  string `json:"token"`
					Socket string `json:"socket"`
				} `json:"data"`
			}
			if err := json.NewDecoder(res.Body).Decode(&socketInfo); err != nil {
				return fmt.Errorf("failed to decode pterodactyl websocket information for server %v: %s",
					s.Config.ServerName, err)
			}
			var auth = []byte(fmt.Sprintf(`{"event":"auth","args":["%v"]}`, socketInfo.Data.Token))

			if !s.connected {
				s.socket, _, err = websocket.DefaultDialer.Dial(socketInfo.Data.Socket, nil)
				if err != nil {
					return fmt.Errorf("failed to connect to pterodactyl websocket for server %v: %s", s.Config.ServerName, err)
				}
			}
			if err := s.socket.WriteMessage(websocket.TextMessage, auth); err != nil {
				return fmt.Errorf("failed to send auth to pterodactyl websocket for server %v: %s", s.Config.ServerName, err)
			}
			var (
				event eventType
			)
			if err := s.socket.ReadJSON(&event); err != nil {
				log.Printf("failed to read websocket message: %s", err)
				return err
			}
			if event.Event == websocketAuthSuccess {
				return nil
			} else {
				return fmt.Errorf("failed to authenticate to pterodactyl websocket for server %v: %s", s.Config.ServerName, err)
			}

		} else {
			return fmt.Errorf("could not connect to websocket: %v", res.Status)
		}
	} else {
		return fmt.Errorf("cannot reach pterodactyl instance with panel url %v", conf.Config.Pterodactyl.PanelURL)
	}
}

func (s *Server) Listen() error {
	if !s.connected {
		if err := s.connectWS(); err != nil {
			return err
		} else {
			s.connected = true
		}
	}
	for {
		var (
			event eventType
		)
		if err := s.socket.ReadJSON(&event); err != nil {
			log.Printf("failed to read websocket message: %s", err)
			continue
		}
		if event.Event == websocketTokenExpired || event.Event == websocketTokenExpiring {
			if event.Event == websocketTokenExpired {
				s.connected = false
			}
			if err := s.connectWS(); err != nil {
				var tries int
				for tries < 5 && err != nil {
					if err := s.connectWS(); err != nil {
						tries++
					}
				}
				if err != nil {
					return err
				}
			}
			continue
		}
		s.setStats(&event)
	}
}

func (s *Server) setStats(data *eventType) {
	switch data.Event {
	case WebsocketConsoleOutput:
		s.Console <- data.Args[0]
	case WebsocketStatus:
		s.Status.State = data.Args[0]
		s.Data <- &ChanData{
			Event: WebsocketStatus,
			Data:  s.Status,
		}
	case WebsocketStats:
		var stats ServerStatus
		if err := json.NewDecoder(bytes.NewBufferString(data.Args[0])).Decode(&stats); err != nil {
			log.Printf("failed to decode stats: %s", err)
			return
		}
		s.Status = &stats
		s.Data <- &ChanData{
			Event: WebsocketStats,
			Data:  s.Status,
		}
	}
}

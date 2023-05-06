package types

type EventJson struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Type   string `json:"type"`
	Server string `json:"server"`
}

type WebsocketEvent struct {
	Event string             `json:"event"`
	Data  WebsocketEventData `json:"data"`
}

type WebsocketEventData struct {
	Players            []string `json:"players,omitempty"`
	Player             string   `json:"player,omitempty"`
	UUID               string   `json:"uuid,omitempty"`
	Reason             string   `json:"reason,omitempty"`
	Command            string   `json:"command,omitempty"`
	Message            string   `json:"message,omitempty"`
	DeathMessage       string   `json:"death_message,omitempty"`
	MessageIsComponent bool     `json:"message_is_component,omitempty"`
	Advancement        string   `json:"advancement,omitempty"`
	Password           string   `json:"password,omitempty"`
	User               string   `json:"user,omitempty"`
	Error              string   `json:"error,omitempty"`
	Server             string   `json:"server,omitempty"`
}

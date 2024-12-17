package types

type ChanData struct {
	Event string
	Data  *ServerStatus
}

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

//goland:noinspection GoUnusedConst
const (
	WebsocketAuthSuccess   = "auth success"
	WebsocketStatus        = "status"
	WebsocketConsoleOutput = "console output"
	WebsocketStats         = "stats"
	WebsocketTokenExpiring = "token expiring"
	WebsocketTokenExpired  = "token expired"

	PowerSignalStart   = "start"
	PowerSignalStop    = "stop"
	PowerSignalKill    = "kill"
	PowerSignalRestart = "restart"

	PowerStatusRunning  = "running"
	PowerStatusOffline  = "offline"
	PowerStatusStarting = "starting"
	PowerStatusStopping = "stopping"
)

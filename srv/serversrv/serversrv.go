package serversrv

import (
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/discord/bot/auth"
	"github.com/Sharktheone/ScharschBot/pterodactyl/types"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"math"
	"strings"
	"time"
)

type replacer struct {
	string string
}

func ChannelStats(status *types.ServerStatus, server *conf.Server) {
	f := replacer{
		string: func(state string) string {
			if state != types.PowerStatusOffline {
				return server.ChannelInfo.Format
			} else {
				return server.ChannelInfo.OfflineFormat
			}
		}(status.State),
	}
	title := cases.Title(language.Und)
	if server.ChannelInfo.Enabled {
		var (
			replace = map[string]string{
				"cpu":        fmt.Sprintf("%.2f", status.Cpu),
				"state":      title.String(status.State),
				"ram":        convertSize(status.Ram),
				"disk":       convertSize(status.Disk),
				"networkIn":  convertSize(status.Network.Rx),
				"networkOut": convertSize(status.Network.Tx),
				"uptime":     convertTime(status.Uptime),
			}
		)
		for variable, value := range replace {
			f.inject(variable, value)
		}
	}
	if server.ChannelInfo.InfoState != f.string {
		server.ChannelInfo.InfoState = f.string
		for _, channelID := range server.ChannelInfo.ChannelID {
			if _, err := auth.Session.ChannelEditComplex(channelID, &discordgo.ChannelEdit{
				Topic: f.string,
			}); err != nil {
				log.Printf("Failed to edit channel topic: %v (channelID %v)", err, channelID)
			}
		}
	}

}

func (r *replacer) inject(variable string, value string) {
	var (
		formats = []string{
			"{{%v}}",
			"${%v}",
		}
	)
	for _, format := range formats {
		r.string = strings.ReplaceAll(r.string, fmt.Sprintf(format, variable), value)
	}
}

func serverStarting(server *conf.Server) {
	if server.StateMessages.StartEnabled {
		for _, channelID := range server.StateMessages.ChannelID {
			_, err := auth.Session.ChannelMessageSend(channelID, server.StateMessages.Start)
			if err != nil {
				log.Printf("Failed to send server start message to discord: %v, (channelID %v)", err, channelID)
			}
		}
	}
}

func serverStopping(server *conf.Server) {
	if server.StateMessages.StopEnabled {
		for _, channelID := range server.StateMessages.ChannelID {
			_, err := auth.Session.ChannelMessageSend(channelID, server.StateMessages.Stop)
			if err != nil {
				log.Printf("Failed to send server stop message to discord: %v (channelID: %v)", err, channelID)
			}
		}
	}
}

func serverOnline(server *conf.Server) {
	if server.StateMessages.OnlineEnabled {
		for _, channelID := range server.StateMessages.ChannelID {
			_, err := auth.Session.ChannelMessageSend(channelID, server.StateMessages.Online)
			if err != nil {
				log.Printf("Failed to send server online message to discord: %v (channelID: %v)", err, channelID)
			}
		}
	}
}

func serverOffline(server *conf.Server) {
	if server.StateMessages.OfflineEnabled {
		for _, channelID := range server.StateMessages.ChannelID {
			_, err := auth.Session.ChannelMessageSend(channelID, server.StateMessages.Offline)
			if err != nil {
				log.Printf("Failed to send server offline message to discord: %v (channelID: %v)", err, channelID)
			}
		}
	}
}

func HandlePower(status string, server *conf.Server) {
	if status == "" {
		return
	}
	switch status {
	case types.PowerStatusStarting:
		serverStarting(server)
	case types.PowerStatusStopping:
		serverStopping(server)
	case types.PowerStatusRunning:
		serverOnline(server)
	case types.PowerStatusOffline:
		serverOffline(server)
	default:
		log.Printf("Unknown power state: %v", status)
	}
}
func convertSize(bytes int) string {
	floatBytes := float64(bytes)
	if bytes < 1024 {
		return fmt.Sprintf("%.2f B", floatBytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.2f KB", floatBytes/1024)
	} else if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", floatBytes/(1024*1024))
	} else if bytes < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2f GB", floatBytes/(1024*1024*1024))
	} else {
		return fmt.Sprintf("%.2f TB", floatBytes/(1024*1024*1024*1024))
	}
}

func convertTime(milliseconds int) string {
	duration := time.Duration(milliseconds * int(time.Millisecond))
	if duration.Hours() >= 24 {
		days := math.Floor(duration.Hours() / 24)
		return fmt.Sprintf("%vd %vh %vm", days, math.Floor(duration.Hours()-(days*24)), math.Floor(duration.Minutes()-(math.Floor(duration.Hours())*60)))
	} else {
		return fmt.Sprintf("%v", duration)
	}
}

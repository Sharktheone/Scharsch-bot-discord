package main

import (
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database/mongodb"
	"github.com/Sharktheone/ScharschBot/diagnostics/pprof"
	"github.com/Sharktheone/ScharschBot/discord/bot"
	"github.com/Sharktheone/ScharschBot/discord/embed/wEmbed"
	"github.com/Sharktheone/ScharschBot/srv"
	"github.com/Sharktheone/ScharschBot/whitelist/checkroles"
	"github.com/robfig/cron"
	"log"
	"os"
	"os/signal"
)

var config = conf.Config

//TODO: Waitlist for whitelist, when server is offline

func main() {
	pprof.Start()
	mongodb.Connect()
	dcBot := bot.Session
	bot.Registration()
	if config.Whitelist.Enabled {
		checkroles.CheckRoles()
		rolesCron := cron.New()
		err := rolesCron.AddFunc("0 */10 * * * *", checkroles.CheckRoles)
		if err != nil {
			log.Fatalf("Error adding RolesCron job: %v", err)
		}
		rolesCron.Start()
	}
	wEmbed.BotAvatarURL = dcBot.State.User.AvatarURL("40")
	srv.Start()

	defer mongodb.Disconnect()
	defer mongodb.Cancel()
	defer dcBot.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
	defer bot.RemoveCommands()
}

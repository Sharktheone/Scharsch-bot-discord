package main

import (
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/database/dbprovider"
	"github.com/Sharktheone/ScharschBot/discord/bot"
	"github.com/Sharktheone/ScharschBot/discord/bot/auth"
	"github.com/Sharktheone/ScharschBot/discord/embed/wEmbed"
	"github.com/Sharktheone/ScharschBot/srv"
	"github.com/Sharktheone/ScharschBot/whitelist"
	"github.com/Sharktheone/ScharschBot/whitelist/checkroles"
	"github.com/robfig/cron"
	"log"
	"os"
	"os/signal"
)

//TODO: Waitlist for whitelist, when server is offline

func main() {
	log.Println("Starting ScharschBot")

	conf.LoadConf()
	bot.Connect()

	whitelist.SetupProvider()

	//pprof.Start()
	dbprovider.Connect()
	log.Println("Connected to MongoDB")
	dcBot := auth.Session
	bot.Registration()
	if conf.Config.Whitelist.Enabled {
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

	defer func() {
		bot.RemoveCommands()
		err := dcBot.Close()
		if err != nil {
			log.Fatalf("Error closing Discord session: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

}

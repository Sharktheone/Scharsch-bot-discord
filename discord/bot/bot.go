package bot

import (
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/console"
	"github.com/Sharktheone/ScharschBot/discord/bot/auth"
	"github.com/Sharktheone/ScharschBot/discord/interactions"
	"github.com/Sharktheone/ScharschBot/discord/session"
	"github.com/Sharktheone/ScharschBot/flags"
	"github.com/bwmarrin/discordgo"
	"log"
)

func Connect() {
	if *auth.GuildID == "" {
		auth.GuildID = &conf.Config.Discord.ServerID
	}
	var BotToken = flags.StringWithFallback("token", &conf.Config.Discord.Token)
	s, err := discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatal("Invalid Bot Configuration:", err)
	}
	auth.Session = &session.Session{Session: s, Guild: *auth.GuildID}

	if err := auth.Session.Open(); err != nil {
		log.Fatal("Cannot open connection to discord:", err)
	}
}

func Registration() {
	log.Println("Registering commands...")

	auth.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := interactions.Handlers[i.ApplicationCommandData().Name]; ok {
				h(&session.Session{Session: s, Guild: *auth.GuildID}, i)
			} else {
				log.Printf("No handler for %v", i.ApplicationCommandData().Name)
			}

		case discordgo.InteractionMessageComponent:
			if h, ok := interactions.Handlers[i.MessageComponentData().CustomID]; ok {
				h(&session.Session{Session: s, Guild: *auth.GuildID}, i)
			} else {
				log.Printf("No handler for %v", i.MessageComponentData().CustomID)
			}
		}
	})

	for _, rawCommand := range interactions.Commands {
		_, err := auth.Session.ApplicationCommandCreate(auth.Session.State.User.ID, *auth.GuildID, rawCommand)
		if err != nil {
			log.Fatalf("Failed to create %v: %v", rawCommand.Name, err)
		}
	}
	auth.Session.AddHandler(console.Handler)
	auth.Session.AddHandler(console.ChatHandler)
	log.Println("Commands registered")

}

func RemoveCommands() {
	commands, err := auth.Session.ApplicationCommands(auth.Session.State.User.ID, *auth.GuildID)
	if err != nil {
		log.Fatalf("Failed to get commands: %v", err)
	}

	for _, command := range commands {
		err := auth.Session.ApplicationCommandDelete(auth.Session.State.User.ID, *auth.GuildID, command.ID)
		if err != nil {
			log.Printf("Failed to delete %v: %v", command.Name, err)
		}
	}
}

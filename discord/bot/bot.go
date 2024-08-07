package bot

import (
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/console"
	"github.com/Sharktheone/ScharschBot/discord/interactions"
	"github.com/Sharktheone/ScharschBot/discord/session"
	"github.com/Sharktheone/ScharschBot/flags"
	"github.com/bwmarrin/discordgo"
	"log"
)

var (
	config  = conf.Config
	GuildID = flags.StringWithFallback("guild", &config.Discord.ServerID)
	Session *session.Session
)

func init() {
	var BotToken = flags.StringWithFallback("token", &config.Discord.Token)
	s, err := discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatal("Invalid Bot Configuration:", err)
	}
	Session = &session.Session{Session: s}

	if err := Session.Open(); err != nil {
		log.Fatal("Cannot open connection to discord:", err)
	}
}

func Registration() {
	log.Println("Registering commands...")

	Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := interactions.Handlers[i.ApplicationCommandData().Name]; ok {
				h(&session.Session{Session: s}, i)
			} else {
				log.Printf("No handler for %v", i.ApplicationCommandData().Name)
			}

		case discordgo.InteractionMessageComponent:
			if h, ok := interactions.Handlers[i.MessageComponentData().CustomID]; ok {
				h(&session.Session{Session: s}, i)
			} else {
				log.Printf("No handler for %v", i.MessageComponentData().CustomID)
			}
		}
	})

	for _, rawCommand := range interactions.Commands {
		_, err := Session.ApplicationCommandCreate(Session.State.User.ID, *GuildID, rawCommand)
		if err != nil {
			log.Fatalf("Failed to create %v: %v", rawCommand.Name, err)
		}
	}
	Session.AddHandler(console.Handler)
	Session.AddHandler(console.ChatHandler)
	log.Println("Commands registered")

}

func RemoveCommands() {
	commands, err := Session.ApplicationCommands(Session.State.User.ID, *GuildID)
	if err != nil {
		log.Fatalf("Failed to get commands: %v", err)
	}

	for _, command := range commands {
		err := Session.ApplicationCommandDelete(Session.State.User.ID, *GuildID, command.ID)
		if err != nil {
			log.Printf("Failed to delete %v: %v", command.Name, err)
		}
	}
}

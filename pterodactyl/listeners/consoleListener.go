package listeners

import (
	"context"
	"fmt"
	"github.com/Sharktheone/ScharschBot/conf"
	"github.com/Sharktheone/ScharschBot/discord/bot/auth"
	"log"
)

func ConsoleListener(ctx context.Context, server *conf.Server, console chan string, commandSent chan bool) {
	var (
		consoleOutput []string
	)
	t, err := newTimer(server.Console.MaxTime)
	if err != nil {
		log.Fatalf("Failed to create ticker: %v", err)
		return
	}
	for {
		select {
		case output := <-console:
			consoleOutput = append(consoleOutput, output)
			t.start()
			if len(consoleOutput) > server.Console.MessageLines {
				sendConsoleOutput(server, consoleOutput)
				consoleOutput = []string{}
				t.stop()
			}
		case <-commandSent:
			sendConsoleOutput(server, consoleOutput)
			consoleOutput = []string{}
			t.stop()
		case <-t.c:
			if len(consoleOutput) > 0 {
				sendConsoleOutput(server, consoleOutput)
				consoleOutput = []string{}
				t.stop()
			}
		case <-ctx.Done():
			return
		}
	}
}

func sendConsoleOutput(server *conf.Server, consoleOutput []string) {
	var message string
	for _, line := range consoleOutput {
		message += fmt.Sprintf("\n%v", line)
	}
	for _, channelID := range server.Console.ChannelID {
		_, err := auth.Session.ChannelMessageSend(channelID, fmt.Sprintf("```%v```", message))
		if err != nil {
			log.Printf("Failed to send console to discord: %v (ChannelID: %v)", err, channelID)
		}
	}
}

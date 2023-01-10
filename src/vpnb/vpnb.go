package vpnb

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/electrocrem/vpn_server/cmd/oss"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func BotStateMachine(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		msg.Text = "Приветы! Я предаставляю доступ к VPN серверу, чтобы получить профиль набери команду" +
			"`/genprof`\n Сервер сам находится в Германии, в будущем будет больше локаций"
		switch update.Message.Command() {
		case "help":
			msg.Text = "Я знаю только одну команду /genprof"
		case "genprof":
			msg.Text = "Профиль готов! Установи приложение OPENVPN:\n" +
				"Для PC: https://openvpn.net/community-downloads/\n" +
				"Для IOS: https://apps.apple.com/us/app/openvpn-connect/id590379981 \n" +
				"Для Android: https://play.google.com/store/apps/details?id=net.openvpn.openvpn&hl=en&gl=US \n"
			GeneratePofile(bot, update)
		default:
			msg.Text = "Я не знаю эту команду."
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}

}

func GeneratePofile(bot *tgbotapi.BotAPI, update tgbotapi.Update) tgbotapi.Message {
	randomName := "profile" + strconv.Itoa(rand.Intn(1000000))
	fmt.Printf("%v", randomName)
	oss.LaunchScript("/bin/bash", "./generate_profile.sh", randomName)
	filePath := "/profiles/" + randomName + ".ovpn"
	fmt.Printf("\n%v\n", filePath)
	profileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	profileFileBytes := tgbotapi.FileBytes{
		Name:  randomName + ".ovpn",
		Bytes: profileBytes,
	}
	message, err := bot.Send(tgbotapi.NewDocument(update.Message.Chat.ID, profileFileBytes))
	if err != nil {
		panic(err)
	}
	os.Remove(filePath)
	return message

}

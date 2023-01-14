package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func LaunchScript(cmdName string, pathToScript string, arg string) error {
	cmd := exec.Command(cmdName, pathToScript, arg)
	err := cmd.Run()

	return err

}
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
			"/genprof\n"
		switch update.Message.Command() {
		case "start":
			msg.Text = "Приветы! Я предаставляю доступ к VPN серверу, чтобы получить профиль набери команду" +
				"`/genprof`\n"
		case "help":
			msg.Text = "Я знаю только одну команду /genprof"
		case "genprof":
			msg.Text = "Профиль готов! Установи приложение **OPEN VPN**:\n" +
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
	LaunchScript("/bin/bash", "./generate_profile.sh", randomName)
	filePath := "profiles/" + randomName + ".ovpn"
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
func init() {

}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN_VPN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	BotStateMachine(bot)
}

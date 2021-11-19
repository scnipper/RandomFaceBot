package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"os"
	"strings"
)

const linkToFaceGenerator = "https://thispersondoesnotexist.com/image"

func main() {
	token := os.Getenv("TELEGRAM_BOT_API_TOKEN")
	println("Start bot")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	//bot.Debug = true

	mainBotLoop(bot)
}

func mainBotLoop(bot *tgbotapi.BotAPI) {
	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := bot.GetUpdatesChan(updateConfig)

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {
		// Telegram can send many types of updates depending on what your Bot
		// is up to. We only want to look at messages for now, so we can
		// discard any other updates.
		if update.Message == nil {
			continue
		}

		msgSend, isSendPhoto := messageHandler(update.Message)
		if isSendPhoto {
			fileImg := tgbotapi.NewInputMediaPhoto(tgbotapi.FileReader{"face", getImageReader()})

			_, err := bot.SendMediaGroup(tgbotapi.NewMediaGroup(update.Message.Chat.ID, []interface{}{fileImg}))
			if err != nil {
				panic(err)
			}
		} else {
			if _, err := bot.Send(msgSend); err != nil {
				panic(err)
			}
		}

	}
}

func messageHandler(msg *tgbotapi.Message) (tgbotapi.MessageConfig, bool) {
	text := msg.Text

	isCommand := strings.HasPrefix(text, "/")

	if isCommand {
		if text == "/help" {
			return tgbotapi.NewMessage(msg.Chat.ID, "Вас приветствует бот по рандомной генерации лиц людей которых никогда не существовало. Для получения изображения лица напишите /get_face"), false
		} else if text == "/get_face" {
			return tgbotapi.MessageConfig{}, true
		}
	}

	return tgbotapi.NewMessage(msg.Chat.ID, "Для получения помощи напишите /help"), false
}

func getImageReader() io.Reader {
	resp, err := http.Get(linkToFaceGenerator)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return resp.Body
}

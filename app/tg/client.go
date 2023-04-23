package tg

import (
	"fmt"
	"os"
	"strings"

	"chad4gpt/app/gpt"
	"chad4gpt/app/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

type TelegramBot interface {
	GetUpdatesChan(tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
}

type Telegram struct {
	Bot          TelegramBot
	Gpt          *gpt.GptClient
	updateConfig tgbotapi.UpdateConfig
}

func NewTelegram() *Telegram {
	apiKey := os.Getenv("TELEGRAM_API_KEY")
	if apiKey == "" {
		log.Fatal().Msg("required TELEGRAM_API_KEY env var is missing")
	}
	tg := Telegram{}
	tg.Gpt = gpt.NewGptClient()
	tg.Init(apiKey)
	return &tg
}

func (tg *Telegram) Init(apiKey string) {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		log.Print("")
	}
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "DEBUG" {
		bot.Debug = true
	}
	tg.updateConfig = tgbotapi.NewUpdate(0)
	tg.updateConfig.Timeout = 60
	tg.Bot = bot
}

func (tg *Telegram) Start() {
	log.Print("reading telegram updates")
	tg.readUpdates(tg.updateConfig)
}

func (tg *Telegram) readUpdates(updateConfig tgbotapi.UpdateConfig) {
	updates := tg.Bot.GetUpdatesChan(updateConfig)
	for upd := range updates {
		if upd.CallbackQuery != nil {
			tg.replyToCallback(*upd.CallbackQuery)
			continue
		}
		if upd.Message.IsCommand() {
			tg.replyToCommand(upd.Message)
			continue
		}
		log.Print("Received a message: ", upd.Message.Text)
		if upd.Message != nil {
			tg.replyToMessage(upd.Message)
		}
	}
}

func (tg *Telegram) replyToMessage(msg *tgbotapi.Message) {
	switch msg.Text {
	default:
		replyMsg := fmt.Sprintf(tg.Gpt.GenerateResponse(msg.Text))
		tg.sendMessage(replyMsg, msg.Chat.ID)
	}
}

func (tg *Telegram) replyToCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "start":
		tg.sendMessage("please send your name and desired language in format <name> <language> \nno braces needed :)", msg.Chat.ID)
	case "stop":
		{
			tg.sendMessage("no notifications anymore for you 😌", msg.Chat.ID)
		}
	}
}

func (tg *Telegram) replyToCallback(cb tgbotapi.CallbackQuery) {}

func (tg *Telegram) sendMessage(msg string, chatId int64) {
	message := tgbotapi.NewMessage(chatId, msg)
	_, err := tg.Bot.Send(message)
	if err != nil {
		log.Fatal().Err(err).Msgf("[ERROR] Occured during sending message to tg chat %d", chatId)
	}
}

func filterChat(chats []storage.Chat, id int64) []storage.Chat {
	for idx, c := range chats {
		if c.ChatId == id {
			return append(chats[0:idx], chats[idx+1:]...)
		}
	}
	return chats
}

func checkIfMessageIsInit(msgText string) bool {
	return len(strings.Split(msgText, " ")) == 2
}

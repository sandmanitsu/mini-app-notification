package bot

import (
	"mini-app-notifications/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramBot(token string) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &TelegramBot{
		bot: bot,
	}, nil
}

func (t *TelegramBot) SendMessage(users []domain.User, text string) {
	for _, user := range users {
		t.bot.Send(tgbotapi.NewMessage(user.ChatId, text))
	}
}

package delivery

// Работа с Telegram API

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

type TelegramNotifier struct {
	// Экземпляр Telegram Bot API
	bot *tgbotapi.BotAPI
}

func NewTelegramNotifier(token string) (*TelegramNotifier, error) {
	// Создание экземпляра Telegram Bot API по токену
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &TelegramNotifier{bot: bot}, nil
}

func (t *TelegramNotifier) SendNotification(message string, chatID int64) error {
	// Создаение сообщения
	msg := tgbotapi.NewMessage(chatID, message)

	// Отправка сообщения
	_, err := t.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}

	return err
}

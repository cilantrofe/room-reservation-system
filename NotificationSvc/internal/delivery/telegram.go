package delivery

// Работа с Telegram API

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

type TelegramNotifier struct {
	// Экземпляр Telegram Bot API
	bot *tgbotapi.BotAPI
	// Идентификатор чата, куда отправляются уведомления
	chatID int64
}

func NewTelegramNotifier(token string, chatID int64) (*TelegramNotifier, error) {
	// Создание экземпляра Telegram Bot API по токену
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &TelegramNotifier{bot: bot, chatID: chatID}, nil
}

func (t *TelegramNotifier) SendNotification(message string) error {
	// Создаение сообщения
	msg := tgbotapi.NewMessage(t.chatID, message)

	// Отправка сообщения
	_, err := t.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}

	return err
}

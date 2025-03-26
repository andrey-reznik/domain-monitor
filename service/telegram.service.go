package service

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"os"
	"os/signal"

	"github.com/nwesterhausen/domain-monitor/configuration"
)

type TelegramService struct {
	client *bot.Bot
}

func NewTelegramService(config configuration.TelegramConfiguration) *TelegramService {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(config.BotID, opts...)

	if nil != err {
		// panics for the sake of simplicity.
		// you should handle this error properly in your code.
		panic(err)
	}
	// Start the bot in a separate goroutine
	go func() {
		defer cancel() // Cancel context when goroutine exits
		b.Start(ctx)
	}()

	return &TelegramService{
		client: b,
	}
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   update.Message.Text,
		})
	}
}

func (m *TelegramService) TestMessage(to string) error {
	ctx := context.Background()

	_, err := m.client.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: to,
		Text:   "Это тестовое сообщение от Domain Monitor.",
	})

	if err != nil {
		log.Printf("❌ Не удалось отправить тестовое сообщение: %s", err)
		return err
	}

	log.Printf("✅ Тестовое сообщение было отправлено в Telegram, СhatID: %s", to)
	return nil
}

func (m *TelegramService) SendAlert(to string, fqdn string, alert configuration.Alert) error {
	ctx := context.Background()
	msg := fmt.Sprintf("Внимание! Срок регистрации домена %s истекает через %s. Пожалуйста, примите во внимание возможность приобретения данного доменного имени после его освобождения.", fqdn, alert)
	if alert == configuration.AlertDaily {
		msg = fmt.Sprintf("Внимание! %s.\nСрок регистрации домена %s скоро истекает. Пожалуйста, примите во внимание возможность приобретения данного доменного имени после его освобождения.", alert, fqdn)
	}
	body := msg

	_, err := m.client.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: to,
		Text:   body,
	})

	if err != nil {
		log.Printf("❌ Failed to send Telegram message: %s", err)
		return err
	}

	log.Printf("✅ Telegram message sent to chat ID: %s", to)
	return nil
}

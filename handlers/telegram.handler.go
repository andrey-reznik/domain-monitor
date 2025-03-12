package handlers

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/nwesterhausen/domain-monitor/service"
)

type TelegramHandler struct {
	TelegramService *service.TelegramService
	chatID          string
}

func NewTelegramHandler(ts *service.TelegramService, chatId string) *TelegramHandler {
	// confirm that the mailer service is not nil
	if ts == nil {
		log.Fatal("🚨 Telegram service not properly initialized.")
	}

	return &TelegramHandler{
		TelegramService: ts,
		chatID:          chatId,
	}
}

func (tgh TelegramHandler) HandleTestMessage(c echo.Context) error {
	err := tgh.TelegramService.TestMessage(tgh.chatID)
	if err != nil {
		log.Printf("❌ Failed to send test Telegram message to chat ID %s: %s", tgh.chatID, err)
		return err
	}
	log.Printf("✅ Test Telegram message sent successfully to chat ID %s", tgh.chatID)
	return c.JSON(200, "Telegram message sent")
}

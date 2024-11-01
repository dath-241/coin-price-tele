package bot

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Mock token for testing
const mockToken = "mock_token"

// Mock webhook URL for testing
const mockWebhookURL = "https://example.com/webhook"

var mockInitBot *tgbotapi.BotAPI

func TestInitBot(t *testing.T) {
	// Mock initialization of the bot with a valid token and webhook URL
	bot, err := InitBot(mockToken, mockWebhookURL)
	if err != nil {
		t.Errorf("Expected nil error, but got %v", err)
	}
	mockInitBot = bot
}

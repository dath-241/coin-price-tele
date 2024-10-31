package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Menu texts
const (
	firstMenu  = "<b>Menu 1</b>\n\nA beautiful menu with a shiny inline button."
	secondMenu = "<b>Menu 2</b>\n\nA better menu with even more shiny inline buttons."
)

// Button texts
const (
	nextButton     = "Next"
	backButton     = "Back"
	tutorialButton = "Tutorial"
)

var (
	screaming = false
)

var (
	// Keyboard layout for the first menu. One button, one row
	firstMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(nextButton, nextButton),
		),
	)

	// Keyboard layout for the second menu. Two buttons, one per row
	secondMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButton, backButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(tutorialButton, "https://core.telegram.org/bots/api"),
		),
	)
)

var commandList = []string{
	"/start - Authenticate and start using the bot",
	"/login - Log in to the bot",
	"/getinfo - Get user info",
	"/scream - Enable screaming mode",
	"/whisper - Disable screaming mode",
	"/menu - Show menu with buttons",
	"/help - Show available commands",
	"/kline - Fetches Kline data for a specific trading pair and interval",
	"/price_spot - Retrieve the latest spot price of a symbol",
	"/price_future - Retrieve the latest futures price of a symbol",
	"/funding_rate - Fetch the current funding rate for a symbol",
	"/funding_rate_countdown - Show the countdown to the next funding time for a symbol",
}

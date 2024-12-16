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
	"help - Show available commands",
	"start - Authenticate and start using the bot",
	"login - Log in to the bot",
	"getinfo - Get user info",
	"menu - Show menu with buttons",
	"p <symbol> - Fetches the price and funding rate of a specific trading pair",
	"marketcap <symbol> - Fetches the marketcap of a specific trading pair",
	"volume <symbol> - Fetches the volume of a specific trading pair",
	"price_spot <symbol> - Fetches the price spot of a specific trading pair",
	"price_futures <symbol> - Fetches the price futures of a specific trading pair",
	"funding_rate <symbol> - Fetches the funding rate of a specific trading pair",
	"all_triggers - Fetches all triggers",
	"delete_trigger <type> <symbol> - Deletes a trigger",
	"alert_price_with_threshold <type> <direction> <symbol> <threshold> - Alerts the price with a threshold",
	"price_difference <direction> <symbol> <threshold> - Alerts the price difference with a threshold",
	"funding_rate_change <direction> <symbol> <threshold> - Alerts the funding rate change with a threshold",
	"kline - Fetches Kline data for a specific trading pair and interval with two choice 'ondemand' or 'realtime'",
	"<symbol> - Fetches the price a specific trading pair",
	"mute <on/off> - Mute bot",
}

type UserInfo struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	VipRole  int    `json:"vipRole"`
	Username string `json:"username"`
	Coin     int    `json:"coin"`
}

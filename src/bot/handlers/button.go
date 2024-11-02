package handlers

import (
	"log"
	"telegram-bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handle inline button clicks
func HandleButton(query *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI) {
	//log.Println("HandleButton in")
	var text string

	token, err := services.GetUserToken(int(query.From.ID))
	if err != nil {
		log.Println("Error getting user token:", err)
	}

	symbol, err := services.GetUserSymbol(int(query.From.ID))
	if err != nil {
		log.Println("Error getting user symbol:", err)
	}

	//log.Println("symbol in HandleButton:", symbol)

	markup := tgbotapi.NewInlineKeyboardMarkup()
	message := query.Message

	if query.Data == nextButton {
		text = secondMenu
		markup = secondMenuMarkup
	} else if query.Data == backButton {
		text = firstMenu
		markup = firstMenuMarkup
	} else {
		HandlePriceCallback(query, bot, token, symbol)
	}

	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	_, err = bot.Request(callbackCfg)
	if err != nil {
		log.Println("Error sending callback:", err)
	}

	// Replace menu text and keyboard
	msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err = bot.Send(msg)
	if err != nil {
		log.Println("Error editing message:", err)
	}
	//log.Println("HandleButton out")
}

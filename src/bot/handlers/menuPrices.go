package handlers

import (
	//"log"
	//"telegram-bot/services"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// const (
// 	Menu = "<i>Menu </i>\n\n<b>Please select the information you want to view:</b>"
// )

const (
	spotPriceButton    = "🟢 |  Spot Price"
	futuresPriceButton = "🔴 |  Futures Price"
	fundingRateButton  = "⚖️ |  Funding Rate"
)

const (
	callbackSpotPrice    = "spot_price"
	callbackFuturesPrice = "futures_price"
	callbackFundingRate  = "funding_rate"
)

func GetPriceMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(spotPriceButton, callbackSpotPrice),
			tgbotapi.NewInlineKeyboardButtonData(futuresPriceButton, callbackFuturesPrice),
			tgbotapi.NewInlineKeyboardButtonData(fundingRateButton, callbackFundingRate),
		),
	)
}

// HandleMessage xử lý tin nhắn văn bản từ người dùng
func HelperMenuPrices(message *tgbotapi.Message, bot *tgbotapi.BotAPI, token string, symbol string) error {
	//fmt.Println("HelperMenuPrices in")
	chatID := message.Chat.ID

	//fmt.Printf("Processing request for symbol: %s\n", symbol)
	//fmt.Printf("Message text: %s\n", message.Text)

	var err error
	switch message.Text {
	case callbackSpotPrice:
		//fmt.Println("Processing spot price request")
		closestSymbol := FindSpotSymbol(symbol)
		if closestSymbol != "" {
			go GetSpotPriceStream(chatID, closestSymbol, bot, token)
		} else {
			fmt.Println("No symbol found.")
		}
	case callbackFuturesPrice:
		closestSymbol := FindFuturesSymbol(symbol)
		if closestSymbol != "" {
			go GetFuturesPriceStream(chatID, closestSymbol, bot, token)
		} else {
			fmt.Println("No symbol found.")
		}
		//go GetFuturesPriceStream(chatID, symbol, bot, token)
	case callbackFundingRate:
		closestSymbol := FindFuturesSymbol(symbol)
		if closestSymbol != "" {
			go GetFundingRateStream(chatID, closestSymbol, bot, token)
		} else {
			fmt.Println("No symbol found.")
		}
		//go GetFundingRateStream(chatID, symbol, bot, token)
	default:
		err = fmt.Errorf("unknown price type: %s", message.Text)
	}

	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error: %v", err))
		bot.Send(msg)
		return err
	}

	//fmt.Println("HelperMenuPrices out")
	return nil
}

func HandlePriceCallback(callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, token string, symbol string) error {
	//fmt.Println("HandlePriceCallback in")

	chatID := callback.Message.Chat.ID

	// Trả lời callback query để loại bỏ loading indicator
	callbackCfg := tgbotapi.NewCallback(callback.ID, "")
	if _, err := bot.Request(callbackCfg); err != nil {
		return fmt.Errorf("error answering callback query: %v", err)
	}

	message := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: chatID},
		Text: callback.Data,
	}

	HelperMenuPrices(message, bot, token, symbol)
	//fmt.Println("HandlePriceCallback out")
	return nil
}

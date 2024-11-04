package handlers

import (
	//"log"
	//"telegram-bot/services"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	Menu = "<i>Menu </i>\n\n<b>Please select the information you want to view:</b>"
)

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

// UserState lưu trữ trạng thái của người dùng
type UserState struct {
	AwaitingSymbol bool
	PriceType      string
}

// userStates lưu trữ trạng thái của tất cả người dùng
var userStates = make(map[int64]*UserState)

// HandleMessage xử lý tin nhắn văn bản từ người dùng
func HelperMenuPrices(message *tgbotapi.Message, bot *tgbotapi.BotAPI, token string, symbol string) error {
	//fmt.Println("HelperMenuPrices in")
	chatID := message.Chat.ID

	fmt.Printf("Processing request for symbol: %s\n", symbol)
	fmt.Printf("Message text: %s\n", message.Text)

	var err error
	switch message.Text {
	case callbackSpotPrice:
		fmt.Println("Processing spot price request")
		GetSpotPriceStream(chatID, symbol, bot, token)
	case callbackFuturesPrice:
		GetFuturesPriceStream(chatID, symbol, bot, token)
	case callbackFundingRate:
		GetFundingRateStream(chatID, symbol, bot, token)
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

	// Gọi trực tiếp HelperMenuPrices với callback.Data
	message := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: chatID},
		Text: callback.Data,
	}
	go HelperMenuPrices(message, bot, token, symbol)

	//fmt.Println("HandlePriceCallback out")
	return nil
}

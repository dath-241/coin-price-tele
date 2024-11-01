package handlers

import (
	//"log"
	//"telegram-bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	spotPriceButton    = "Spot price"
	futuresPriceButton = "Futures price"
	fundingRateButton  = "Funding rate"
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

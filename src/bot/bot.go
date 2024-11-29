package bot

import (
	"context"
	"log"
	"telegram-bot/bot/handlers"

	// "time"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	// "sync"
	// "bytes"

	// "sync"
	// "bytes"
	"io"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

var commands = []tgbotapi.BotCommand{
	{
		Command:     "start",
		Description: "Authenticate and start using the bot",
	},
	{
		Command:     "login",
		Description: "Log in to the bot",
	},
	{
		Command:     "register",
		Description: "Register new user",
	},
	{
		Command:     "forgotpassword",
		Description: "Send OTP to email to get new password",
	},
	{
		Command:     "getinfo",
		Description: "Get user information",
	},
	{
		Command:     "help",
		Description: "Show available commands",
	},
	{
		Command:     "kline",
		Description: "Get Kline data and candlestick chart for a symbol ondemand or realtime",
	},
	{
		Command:     "p",
		Description: "<symbol>",
	},
	//----------------------------------------------------------------------------------------
	{
		Command:     "alert_price_with_threshold",
		Description: "<spot/future> <lower/above> <symbol> <threshold>",
	},
	{
		Command:     "price_difference",
		Description: "<lower/above> <symbol> <threshold>",
	},
	{
		Command:     "funding_rate_change",
		Description: "<lower/above> <symbol> <threshold>",
	},
	{
		Command:     "all_triggers",
		Description: "Get all triggers",
	},
	{
		Command:     "delete_trigger",
		Description: "<spot/future/price-difference/funding-rate> <symbol>",
	},
}

// send from BE
type CoinPriceUpdate struct {
	Symbol      string  `json:"symbol"`
	Spotprice   float64 `json:"spot_price"`
	Futureprice float64 `json:"future_price"`
	Pricediff   float64 `json:"price_diff"`
	Fundingrate float64 `json:"fundingrate"`
	Threshold   float64 `json:"threshold"`
	Condition   string  `json:"condition"`
	ChatID      string  `json:"chatID"`
	Timestamp   string  `json:"timestamp"`
	Triggertype string  `json:"triggerType"` //spot, price-difference, funding-rate, future
}

type IndicatorUpdate struct {
	Symbol         string  `json:"symbol"`
	Condition      string  `json:"condition"`
	ChatID         string  `json:"chatID"`
	Timestamp      string  `json:"timestamp"`
	Indicator      string  `json:"indicator"`
	IndicatorValue float64 `json:"indicatorValue"`
	Threshold      float64 `json:"threshold"`
	Period         int     `json:"period"`
	Triggertype    string  `json:"triggerType"` //indicator
}

// Initialize the bot with the token
func InitBot(bottoken string, webhookURL string) (*tgbotapi.BotAPI, error) {
	var err error
	bot, err = tgbotapi.NewBotAPI(bottoken)
	if err != nil {
		return nil, err
	}
	bot.Debug = false // Set to true if you want to debug interactions
	webhook, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		return nil, err
	}
	//Sau khi tạo Webhook, bạn cần gửi nó đến Telegram để cấu hình:
	_, err = bot.Request(webhook)
	if err != nil {
		return nil, err
	}

	//Đặt danh sách các lệnh (commands) cho bot Telegram.
	_, err = bot.Request(tgbotapi.NewSetMyCommands(commands...))
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Start")
	return bot, nil
}

// Start listening for updates
func Start(ctx context.Context, bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Pass updates to handler
	go receiveUpdates(ctx, updates)
}

// Start listening update from webhook
func StartWebhook(bot *tgbotapi.BotAPI) {
	//Create the update channel using ListenForWebhook
	updates := bot.ListenForWebhook("/webhook")
	for update := range updates {
		if update.Message != nil {
			if update.Message.Text == "/kline" || update.Message.Text == "ondemand" || update.Message.Text == "realtime" ||
				update.Message.Text == "/Resume" || update.Message.Text == "Stop" || update.Message.Text == "Chart" ||
				handlers.UserSelections[update.Message.Chat.ID]["step"] == "coin_selection" ||
				handlers.UserSelections[update.Message.Chat.ID]["step"] == "interval_selection" ||
				handlers.UserSelections[update.Message.Chat.ID]["step"] == "other_input" ||
				handlers.UserSelections[update.Message.Chat.ID]["step"] == "fetching_data" {
				handlers.HandleKlineCommand(update.Message.Chat.ID, update.Message.Text, bot, update.Message.From)
				log.Printf("Kline: %s", update.Message.Text)
			} else {
				handlers.HandleMessage(update.Message, bot)
			}
		} else if update.CallbackQuery != nil {
			handlers.HandleButton(update.CallbackQuery, bot)
		}
	}
}

// Receive updates and pass them to handlers
func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			if update.Message != nil {
				if update.Message.Text == "/kline" || update.Message.Text == "ondemand" || update.Message.Text == "realtime" ||
					handlers.UserSelections[update.Message.Chat.ID]["step"] == "coin_selection" ||
					handlers.UserSelections[update.Message.Chat.ID]["step"] == "interval_selection" ||
					handlers.UserSelections[update.Message.Chat.ID]["step"] == "other_input" ||
					handlers.UserSelections[update.Message.Chat.ID]["step"] == "fetching_data" {
					handlers.HandleKlineCommand(update.Message.Chat.ID, update.Message.Text, bot, update.Message.From)
				} else {
					handlers.HandleMessage(update.Message, bot)
				}
			} else if update.CallbackQuery != nil {
				handlers.HandleButton(update.CallbackQuery, bot)
			}
		}
	}
}

func PriceUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var temp struct {
		TriggerType string `json:"triggerType"`
	}

	err = json.Unmarshal(body, &temp)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	if temp.TriggerType == "indicator" {
		var update IndicatorUpdate
		err = json.Unmarshal(body, &update)
		if err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}
		fmt.Printf("Received price update: Coin: %s, Trigger Type: %s, Timestamp: %s\n", update.Symbol, update.Triggertype, update.Timestamp)

		direction := "below"
		if update.Condition == ">=" || update.Condition == ">" {
			direction = "above"
		}
		chatID, err := strconv.ParseInt(update.ChatID, 10, 64)
		if err != nil {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}
		mess := fmt.Sprintf("🚨Price alert:\n👉Coin: %s is %s indicator: %s \n👉Current value: %.2f",
			update.Symbol, direction, update.Indicator, update.IndicatorValue)
		go handlers.SendMessageToUser(bot, chatID, mess)

	} else {
		var update CoinPriceUpdate
		err = json.Unmarshal(body, &update)
		if err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}
		fmt.Printf("Received price update: Coin: %s, Price: %.2f, Timestamp: %s\n", update.Symbol, update.Threshold, update.Timestamp)

		direction := "below"
		if update.Condition == ">=" || update.Condition == ">" {
			direction = "above"
		}
		chatID, err := strconv.ParseInt(update.ChatID, 10, 64)
		if err != nil {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}
		var mess string
		if update.Triggertype == "spot" {
			mess = fmt.Sprintf("🚨Price alert:\n👉Coin: %s is %s spot price threshold: %.2f\n👉Current spot price: %.2f",
				update.Symbol, direction, update.Threshold, update.Spotprice)

		} else if update.Triggertype == "future" {
			mess = fmt.Sprintf("🚨Price alert:\n👉Coin: %s is %s future price threshold: %.2f\n👉Current future price: %.2f",
				update.Symbol, direction, update.Threshold, update.Futureprice)

		} else if update.Triggertype == "funding-rate" {
			mess = fmt.Sprintf("🚨Funding rate alert:\n👉Coin: %s is %s funding rate threshold: %.2f\n👉Current funding rate: %.2f",
				update.Symbol, direction, update.Threshold, update.Fundingrate)

		} else if update.Triggertype == "price-difference" {
			mess = fmt.Sprintf("🚨Price alert:\n👉Coin: %s is %s Price-diff threshold: %.2f\n👉Current spot price: %.2f, Current future price: %.2f",
				update.Symbol, direction, update.Pricediff, update.Spotprice, update.Futureprice)
		}
		go handlers.SendMessageToUser(bot, chatID, mess)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Price update received"))
}

// {
//     "symbol": "BTC",
//     "condition": ">=",
//     "chatID": "",
//     "timestamp": "2023-10-01T12:00:00Z",
//     "indicator": "EMA",
// 		"indicatorValue" : "",
//     "value": 70.5000,
//     "period": "14",
//     "triggerType": "indicator"
// }

// {
//     "symbol": "BTC",
//     "spot_price": 45000.00,
//     "threshold": 44000.00,
//     "condition": ">=",
//     "chatID": "",
//     "timestamp": "2023-10-01T12:00:00Z",
//     "value": 45000.00,
//     "triggerType": "spot"
// }

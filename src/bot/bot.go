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
		Description: "Get Kline data on demand for a symbol",
	},
	{
		Command:     "kline_realtime",
		Description: "Get realtime Kline data for a symbol",
	},
	{
		Command:     "stop",
		Description: "Stop receiving Kline_realtime",
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
	Indicator   string  `json:"indicator"`
	Value       float64 `json:"value"`
	Period      string  `json:"period"`
	Triggertype string  `json:"triggerType"` //spot, price-difference, funding-rate, future
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
			handlers.HandleMessage(update.Message, bot)
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
				handlers.HandleMessage(update.Message, bot)
			} else if update.CallbackQuery != nil {
				handlers.HandleButton(update.CallbackQuery, bot)
			}
		}
	}
}

func PriceUpdateHandler(w http.ResponseWriter, r *http.Request) {
	//? nhan lenh post -> gui cho user
	//? print user
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var update CoinPriceUpdate
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Process the received data
	fmt.Printf("Received price update: Coin: %s, Price: %.2f, Timestamp: %s\n", update.Symbol, update.Threshold, update.Timestamp)
	// Sử dụng WaitGroup để quản lý các goroutine
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
	} else if update.Triggertype == "indicator" {
		mess = fmt.Sprintf("🚨Price alert:\n👉Coin: %s is %s indicator: %s \n👉Current value: %.2f",
			update.Symbol, direction, update.Indicator, update.Value)
	}
	go handlers.SendMessageToUser(bot, chatID, mess)

	// Respond to the sender
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Price update received"))
}

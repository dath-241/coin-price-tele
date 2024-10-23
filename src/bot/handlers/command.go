package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"

	"telegram-bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	binanceSpotWSURL   = "wss://stream.binance.com:9443/ws"
	binanceFutureWSURL = "wss://fstream.binance.com/ws"
)

var conn *websocket.Conn

func connectWebSocket(url string) error {
	var err error
	// Create a background context
	ctx := context.Background()
	conn, _, err = websocket.Dial(ctx, url, nil)
	return err
}

func subscribeToStream(stream string) error {
	// Create a background context for the write operation
	ctx := context.Background()
	subscribeMsg := fmt.Sprintf(`{"method": "SUBSCRIBE", "params":["%s"], "id": 1}`, stream)
	return conn.Write(ctx, websocket.MessageText, []byte(subscribeMsg))
}

func formatNumber(value interface{}) string {
	switch v := value.(type) {
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return fmt.Sprintf("%g", f)
		}
	case float64:
		return fmt.Sprintf("%g", v)
	}
	return fmt.Sprintf("%v", value)
}

var userTokens = struct {
	sync.RWMutex
	m map[int]string
}{m: make(map[int]string)}

// Handle incoming messages (commands or regular text)
func HandleMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	user := message.From
	text := message.Text

	log.Printf("%s wrote: %s", user.FirstName, text)

	if strings.HasPrefix(text, "/") {
		parts := strings.Fields(text)
		command := parts[0]
		args := parts[1:]
		handleCommand(message.Chat.ID, command, args, bot, user)
	} else if screaming {
		_, err := bot.Send(sendScreamedMessage(message))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	} else {
		_, err := bot.Send(copyMessage(message))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	}
}

// Handle commands (e.g., /scream, /whisper, /menu)
func handleCommand(chatID int64, command string, args []string, bot *tgbotapi.BotAPI, user *tgbotapi.User) {
	fmt.Println(user.ID)
	switch command {
	case "/help":
		_, err := bot.Send(tgbotapi.NewMessage(chatID, strings.Join(commandList, "\n")))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	case "/start":
		token, err := services.AuthenticateUser(user.ID)
		if err != nil {
			_, err := bot.Send(tgbotapi.NewMessage(chatID, "Access denied."))
			if err != nil {
				log.Println("Error sending message:", err)
			}
			return
		}

		// Send a message showing access was granted
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Access granted. Your token is: %s", token))
		// Store the token
		userTokens.Lock()
		userTokens.m[int(user.ID)] = token
		userTokens.Unlock()
		_, err = bot.Send(msg)
		if err != nil {
			log.Println("Error sending message:", err)
		}
	case "/scream":
		screaming = true
		_, err := bot.Send(tgbotapi.NewMessage(chatID, "Screaming mode enabled."))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	case "/whisper":
		screaming = false
		_, err := bot.Send(tgbotapi.NewMessage(chatID, "Screaming mode disabled."))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	case "/menu":
		_, err := bot.Send(sendMenu(chatID))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	case "/protected":
		token := userTokens.m[int(user.ID)]
		response, err := services.ValidateToken(token)
		if err != nil {
			log.Println("Error validating token:", err)
			return
		}
		_, err = bot.Send(tgbotapi.NewMessage(chatID, response))
		if err != nil {
			log.Println("Error sending message:", err)
		}

	case "/price_spot":
		if len(args) < 1 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /price_spot <symbol>")
			bot.Send(msg)
			return
		}
		symbol := args[0]
		go getSpotPrice(chatID, symbol, bot)
	case "/price_future":
		if len(args) < 1 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /price_future <symbol>")
			bot.Send(msg)
			return
		}
		symbol := args[0]
		go getFuturePrice(chatID, symbol, bot)
	case "/funding_rate":
		if len(args) < 1 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /funding_rate <symbol>")
			bot.Send(msg)
			return
		}
		symbol := args[0]
		go getFundingRate(chatID, symbol, bot)
	case "/funding_rate_countdown":
		if len(args) < 1 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /funding_rate_countdown <symbol>")
			bot.Send(msg)
			return
		}

		symbol := args[0]
		go getFundingRateCountdown(chatID, symbol, bot)
	}
}

func sendMenu(chatID int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, firstMenu)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = firstMenuMarkup
	return msg
}

func sendScreamedMessage(message *tgbotapi.Message) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, strings.ToUpper(message.Text))
	msg.ParseMode = tgbotapi.ModeHTML
	return msg
}

func copyMessage(message *tgbotapi.Message) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	return msg
}

func getSpotPrice(chatID int64, symbol string, bot *tgbotapi.BotAPI) {
	// Create a background context for WebSocket connection
	ctx := context.Background()
	if err := connectWebSocket(binanceSpotWSURL); err != nil {
		log.Println("Lỗi kết nối WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing")

	stream := strings.ToLower(symbol) + "@ticker"
	if err := subscribeToStream(stream); err != nil {
		log.Println("Lỗi đăng ký stream:", err)
		return
	}

	for {
		_, message, err := conn.Read(ctx)
		if err != nil {
			log.Println("Lỗi đọc tin nhắn:", err)
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("Lỗi giải mã JSON:", err)
			continue
		}

		if price, ok := data["c"]; ok {
			formattedPrice := formatNumber(price)
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Spot price | %s |  %s", symbol, formattedPrice))
			bot.Send(msg)
			return
		}
	}
}

func getFuturePrice(chatID int64, symbol string, bot *tgbotapi.BotAPI) {
	// Create a background context for WebSocket connection
	ctx := context.Background()
	if err := connectWebSocket(binanceFutureWSURL); err != nil {
		log.Println("Lỗi kết nối WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing")

	stream := strings.ToLower(symbol) + "@markPrice"
	if err := subscribeToStream(stream); err != nil {
		log.Println("Lỗi đăng ký stream:", err)
		return
	}

	for {
		_, message, err := conn.Read(ctx)
		if err != nil {
			log.Println("Lỗi đọc tin nhắn:", err)
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("Lỗi giải mã JSON:", err)
			continue
		}

		if price, ok := data["p"]; ok {
			formattedPrice := formatNumber(price)
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Future price | %s |  %s", symbol, formattedPrice))
			bot.Send(msg)
			return
		}
	}
}

func getFundingRate(chatID int64, symbol string, bot *tgbotapi.BotAPI) {
	// Create a background context for WebSocket connection
	ctx := context.Background()
	if err := connectWebSocket(binanceFutureWSURL); err != nil {
		log.Println("Lỗi kết nối WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing")

	stream := strings.ToLower(symbol) + "@markPrice"
	if err := subscribeToStream(stream); err != nil {
		log.Println("Lỗi đăng ký stream:", err)
		return
	}

	for {
		_, message, err := conn.Read(ctx)
		if err != nil {
			log.Println("Lỗi đọc tin nhắn:", err)
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("Lỗi giải mã JSON:", err)
			continue
		}

		if fundingRate, ok := data["r"]; ok {
			formattedFundingRate := formatNumber(fundingRate)
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Funding rate | %s |  %s", symbol, formattedFundingRate))
			bot.Send(msg)
			return
		}
	}
}

func getFundingRateCountdown(chatID int64, symbol string, bot *tgbotapi.BotAPI) {
	// Create a background context for WebSocket connection
	ctx := context.Background()
	if err := connectWebSocket(binanceFutureWSURL); err != nil {
		log.Println("Lỗi kết nối WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "closing")

	stream := strings.ToLower(symbol) + "@markPrice"
	if err := subscribeToStream(stream); err != nil {
		log.Println("Lỗi đăng ký stream:", err)
		return
	}

	for {
		_, message, err := conn.Read(ctx)
		if err != nil {
			log.Println("Lỗi đọc tin nhắn:", err)
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("Lỗi giải mã JSON:", err)
			continue
		}

		if nextFundingTime, ok := data["T"]; ok {
			nextFundingTimeMillis := int64(nextFundingTime.(float64))
			remainingTime := time.Until(time.UnixMilli(nextFundingTimeMillis))
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Funding rate countdown | %s |  %s", symbol, remainingTime))
			bot.Send(msg)
			return
		}
	}
}

package handlers

import (
	// "context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"telegram-bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func init() {
	services.InitDB()
}

var (
	globalSymbol string
	symbolMutex  sync.RWMutex
)

// Handle incoming messages (commands or regular text)
func HandleMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	user := message.From
	text := message.Text

	log.Printf("\n\n%s wrote: %s", user.FirstName+" "+user.LastName, text)

	if strings.HasPrefix(text, "/") {
		parts := strings.Fields(text)
		command := parts[0]
		args := parts[1:]
		handleCommand(message.Chat.ID, command, args, bot, user)
	} else {
		closestSymbol := FindClosestSymbol1(text, SpotSymbols)
		closestSymbol1 := FindClosestSymbol1(text, FuturesSymbols)

		if closestSymbol == "" {
			fmt.Printf("No symbol found.")
			//msg := tgbotapi.NewMessage(chatID, "No symbol found.")
			//bot.Send(msg)
			//
			return
		} else {
			message1 := "/price_spot"
			args := []string{closestSymbol}
			handleCommand(message.Chat.ID, message1, args, bot, user)

			message2 := "/price_futures"
			args = []string{closestSymbol1}
			handleCommand(message.Chat.ID, message2, args, bot, user)

		}
		// _, err := bot.Send(copyMessage(message))
		// if err != nil {
		// 	log.Println("Error sending message:", err)
		// }
	}
}

// Handle commands
func handleCommand(chatID int64, command string, args []string, bot *tgbotapi.BotAPI, user *tgbotapi.User) {
	fmt.Println("userID: ", user.ID)
	switch command {
	case "/help":
		_, err := bot.Send(tgbotapi.NewMessage(chatID, strings.Join(commandList, "\n")))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	case "/start":
		response, err := services.AuthenticateUser(user.ID)
		if err != nil {
			_, err := bot.Send(tgbotapi.NewMessage(chatID, "Access denied."))
			if err != nil {
				log.Println("Error sending message:", err)
			}
			return
		}
		_, err = bot.Send(tgbotapi.NewMessage(chatID, response))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	case "/login":
		// With two args: /login username password
		if len(args) < 2 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /login <username> <password>")
			bot.Send(msg)
			return
		}

		username := args[0]
		password := args[1]
		response, token, err := services.LogIn(username, password)

		if err != nil {
			_, _ = bot.Send(tgbotapi.NewMessage(chatID, "Error logging in: "+err.Error()))
		} else {
			_, _ = bot.Send(tgbotapi.NewMessage(chatID, response))
			err = services.StoreUserToken(int(user.ID), token)
			// Log the token
			log.Println("Token:", token)
			if err != nil {
				log.Println("Error storing token:", err)
			}
		}
	case "/getinfo":
		token, err := services.GetUserToken(int(user.ID))
		if err != nil {
			log.Println("Error retrieving token:", err)
			return
		}
		response, err := services.GetUserInfo(token)
		if err != nil {
			_, _ = bot.Send(tgbotapi.NewMessage(chatID, "Error getting user info: "+err.Error()))
		} else {
			_, _ = bot.Send(tgbotapi.NewMessage(chatID, response))
		}
	case "/kline":
		if len(args) < 2 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /kline <symbol> <interval> [limit] [startTime] [endTime]")
			bot.Send(msg)
			return
		}
		symbol := args[0]
		interval := args[1]
		limit := 5
		if len(args) == 3 {
			parsedLimit, err := strconv.Atoi(args[2])
			if err == nil {
				limit = parsedLimit
			}
		}
		data, err := getKlineData(symbol, interval, limit) // Pass parameters as needed
		if err != nil {
			_, _ = bot.Send(tgbotapi.NewMessage(chatID, "Error fetching Kline data: "+err.Error()))
		} else {
			log.Println(data)
			sendChartToTelegram(bot, chatID, klineBase(data))
		}
	case "/menu":
		_, err := bot.Send(sendMenu(chatID))
		if err != nil {
			log.Println("Error sending message:", err)
		}
	case "/p":
		if len(args) < 1 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /p <symbol>")
			bot.Send(msg)
			return
		}

		symbol := args[0]
		closestSymbol := FindClosestSymbol1(symbol, FuturesSymbols)
		if closestSymbol == "" {
			log.Println("No symbol found.")
			msg := tgbotapi.NewMessage(chatID, "No symbol found.")
			bot.Send(msg)
			return
		} else {
			symbolMutex.Lock()
			globalSymbol = closestSymbol
			symbolMutex.Unlock()
			Menu := fmt.Sprintf("<i>Menu</i>\n\n<b>                                                         %s       </b>\n\nPlease select the information you want to view:", globalSymbol)
			msg := tgbotapi.NewMessage(chatID, Menu)
			msg.ReplyMarkup = GetPriceMenu()
			msg.ParseMode = "HTML"
			if _, err := bot.Send(msg); err != nil {
				log.Println("Error sending message:", err)
			}
		}
	case "/price_spot":
		token, err := services.GetUserToken(int(user.ID))
		if err != nil {
			log.Println("Error retrieving token:", err)
			return
		}
		if len(args) < 1 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /price_spot <symbol>")
			bot.Send(msg)
			return
		}
		symbol := args[0]
		go GetSpotPriceStream(chatID, symbol, bot, token)
	case "/price_futures":
		token, err := services.GetUserToken(int(user.ID))
		if err != nil {
			log.Println("Error retrieving token:", err)
			return
		}

		if len(args) < 1 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /price_futures <symbol>")
			bot.Send(msg)
			return
		}
		symbol := args[0]
		go GetFuturesPriceStream(chatID, symbol, bot, token)
	case "/funding_rate":
		token, err := services.GetUserToken(int(user.ID))
		if err != nil {
			log.Println("Error retrieving token:", err)
			return
		}

		if len(args) < 1 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /funding_rate <symbol>")
			bot.Send(msg)
			return
		}
		symbol := args[0]
		go GetFundingRateStream(chatID, symbol, bot, token)
	// case "/funding_rate_countdown":
	// 	if len(args) < 1 {
	// 		msg := tgbotapi.NewMessage(chatID, "Usage: /funding_rate_countdown <symbol>")
	// 		bot.Send(msg)
	// 		return
	// 	}
	// 	symbol := args[0]
	// 	go GetFundingRateCountdown(chatID, symbol, bot)
	case "/kline_realtime":
		if len(args) != 2 {
			bot.Send(tgbotapi.NewMessage(chatID, "Usage: /kline <symbol> <interval>. Example: /kline BTCUSDT 1m"))
			return
		}

		symbol := args[0]
		interval := args[1]

		mapMutex.Lock()
		userConnections[chatID] = &UserConnection{isStreaming: true}
		mapMutex.Unlock()

		token, err := services.GetUserToken(int(user.ID))
		if err != nil {
			log.Println("Error retrieving token:", err)
			return
		}

		// Start fetching Kline data and sending real-time updates to the user
		go fetchKlineData(symbol, interval, token, chatID, bot)
		bot.Send(tgbotapi.NewMessage(chatID, "Fetching real-time Kline data..."))
	case "/stop":
		mapMutex.Lock()
		if userConn, ok := userConnections[chatID]; ok {
			userConn.isStreaming = false
			bot.Send(tgbotapi.NewMessage(chatID, "Stopped real-time Kline updates."))
		} else {
			bot.Send(tgbotapi.NewMessage(chatID, "No active real-time updates to stop."))
		}
		mapMutex.Unlock()

	//----------------------------------------------------------------------------------------
	case "/all_triggers":
		go GetAllTrigger(chatID, bot)
	case "/delete_trigger":
		if len(args) != 2 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /delete_trigger <spot/future/price-difference/funding-rate> <symbol>")
			bot.Send(msg)
			return
		}
		if args[0] != "spot" && args[0] != "future" && args[0] != "price-difference" && args[0] != "funding-rate" {
			msg := tgbotapi.NewMessage(chatID, "First argument must be either 'spot' or 'future' or 'price-difference' or 'funding-rate'")
			bot.Send(msg)
			return
		}
		price_type := args[0]
		symbol := args[1]
		go DeleteTrigger(chatID, bot, symbol, price_type)
	case "/alert_price_with_threshold":
		if len(args) != 4 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /alert_price_with_threshold <spot/future> <lower/above> <symbol> <threshold>")
			bot.Send(msg)
			return
		}

		// Validate price_type (arg[0])
		price_type := args[0]
		if price_type != "spot" && price_type != "future" {
			msg := tgbotapi.NewMessage(chatID, "First argument must be either 'spot' or 'future'")
			bot.Send(msg)
			return
		}

		// Validate comparison type (arg[1])
		if args[1] != "lower" && args[1] != "above" {
			msg := tgbotapi.NewMessage(chatID, "Second argument must be either 'lower' or 'above'")
			bot.Send(msg)
			return
		}

		is_lower := args[1] == "lower"
		symbol := args[2]
		threshold, err := strconv.ParseFloat(args[3], 64)
		if err != nil {
			log.Println("Error parsing threshold:", err)
			return
		}
		go RegisterPriceThreshold(chatID, symbol, threshold, is_lower, price_type, bot)
	case "/price_difference":
		if len(args) != 3 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /price_difference <lower/above> <symbol> <threshold>")
			bot.Send(msg)
			return
		}
		if args[0] != "lower" && args[0] != "above" {
			msg := tgbotapi.NewMessage(chatID, "First argument must be either 'lower' or 'above'")
			bot.Send(msg)
			return
		}
		is_lower := args[0] == "lower"
		symbol := args[1]
		threshold, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			log.Println("Error parsing threshold:", err)
			return
		}
		go RegisterPriceDifferenceAndFundingRate(chatID, symbol, threshold, is_lower, "price-difference", bot)
	case "/funding_rate_change":
		if len(args) != 3 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /funding_rate_change <lower/above> <symbol> <threshold>")
			bot.Send(msg)
			return
		}
		if args[0] != "lower" && args[0] != "above" {
			msg := tgbotapi.NewMessage(chatID, "First argument must be either 'lower' or 'above'")
			bot.Send(msg)
			return
		}
		is_lower := args[0] == "lower"
		symbol := args[1]
		threshold, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			log.Println("Error parsing threshold:", err)
			return
		}
		go RegisterPriceDifferenceAndFundingRate(chatID, symbol, threshold, is_lower, "funding-rate", bot)
	}
}

func sendMenu(chatID int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, firstMenu)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = firstMenuMarkup
	return msg
}

func copyMessage(message *tgbotapi.Message) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	return msg
}

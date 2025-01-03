package handlers

import (
	// "context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"telegram-bot/config"
	"telegram-bot/services"
	"time"

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

	// Get the mute status of the user
	isMuted, err := services.GetMute(int(user.ID))
	if err != nil {
		log.Println("Error getting mute status:", err)
	}
	if isMuted && !strings.Contains(text, "/mute") {
		return
	}

	// Check if the message is from a group
	if message.Chat.IsGroup() || message.Chat.IsSuperGroup() {
		// Check if the message is a command directed at the bot
		if strings.HasPrefix(text, "/") && strings.Contains(text, "@"+config.GetEnv("BOT_USERNAME")) {
			parts := strings.Fields(text)
			command := parts[0]
			command = strings.TrimSuffix(command, "@"+config.GetEnv("BOT_USERNAME"))
			args := parts[1:]
			handleCommand(message.Chat.ID, command, args, bot, user)
		}
	} else {
		if strings.HasPrefix(text, "/") {
			parts := strings.Fields(text)
			command := parts[0]
			args := parts[1:]
			handleCommand(message.Chat.ID, command, args, bot, user)
		} else {
			closestSymbol := FindSpotSymbol(text)
			closestSymbol1 := FindFuturesSymbol(text)

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

	case "/register":
		//syntax /signup <email> <name> <username> <password>
		if len(args) < 4 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /register <email> <name> <username> <password>")
			bot.Send(msg)
			return
		}
		email := args[0]
		name := args[1]
		username := args[2]
		password := args[3]
		response, err := services.Regsiter(email, name, username, password)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "Error in registering: "+err.Error()))
		} else {
			bot.Send(tgbotapi.NewMessage(chatID, response))
			bot.Send(tgbotapi.NewMessage(chatID, "use /login to log in"))
		}
	case "/mute":
		if len(args) != 1 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /mute <on/off>")
			bot.Send(msg)
			return
		}
		// Get the mute status of the user
		isMuted := args[0] == "on"
		err := services.SetMute(int(user.ID), isMuted)
		if err != nil {
			log.Println("Error setting mute status:", err)
			return
		}
		bot.Send(tgbotapi.NewMessage(chatID, "Mute status set to "+args[0]))
	case "/changeinfo":
		bot.Send(tgbotapi.NewMessage(chatID, "In Progress"))
	case "/getinfo":
		token, err := services.GetUserToken(int(user.ID))
		if err != nil {
			log.Println("Error retrieving token:", err)
			return
		}
		response, err := services.GetUserInfo(token)
		if err != nil {
			log.Println("Error getting user info:", err)
			return
		}
		handleUserInfo(chatID, bot, response)

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

		// Add logging to check symbols
		// log.Printf("Input symbol: %s", args[0])
		// log.Printf("Available SpotSymbols: %v", SpotSymbols)

		symbol := args[0]
		closestSymbol := FindSpotSymbol(symbol)
		nameSymbol := strings.ToUpper(symbol)
		if closestSymbol == "" {
			log.Println("No symbol found.")
			msg := tgbotapi.NewMessage(chatID, "No symbol found.")
			bot.Send(msg)
			return
		} else {
			symbolMutex.Lock()
			globalSymbol = nameSymbol
			symbolMutex.Unlock()
			Menu := fmt.Sprintf("<i>Menu</i>\n\n<b>                                                         %s       </b>\n\nPlease select the information you want to view:", nameSymbol)
			msg := tgbotapi.NewMessage(chatID, Menu)
			msg.ReplyMarkup = GetPriceMenu()
			msg.ParseMode = "HTML"
			if _, err := bot.Send(msg); err != nil {
				log.Println("Error sending message:", err)
			}
		}
	case "/marketcap":
		log.Print("in marketCap")
		token, err := services.GetUserToken(int(user.ID))
		if err != nil {
			log.Println("Error retrieving token:", err)
			return
		}

		if len(args) < 1 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /marketcap <symbol>")
			bot.Send(msg)
			return
		}

		symbol := args[0]
		go GetMarketCap(chatID, symbol, bot, token)
		log.Print("out marketCap")

	case "/volume":
		if len(args) < 1 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /volume <symbol>")
			bot.Send(msg)
			return
		}
		symbol := strings.ToUpper(args[0])
		go GetTradingVolume(chatID, symbol, bot)

	case "/forgotpassword":
		//syntax /forgotpassword <username>
		if len(args) < 1 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /forgotpassword <username>")
			bot.Send(msg)
			return
		}
		username := args[0]
		response, err := services.ForgotPassword(username)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
		} else {
			bot.Send(tgbotapi.NewMessage(chatID, response))
		}

	case "/alert_indicator":
		if len(args) < 4 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /alert_indicator <symbol> <indicator> <condition> <value>")
			bot.Send(msg)
			return
		}
		if args[1] != "EMA" && args[1] != "MA" {
			bot.Send(tgbotapi.NewMessage(chatID, "Indicator không hợp lệ (EMA/MA)"))
			return
		} else {
			bot.Send(tgbotapi.NewMessage(chatID, "Nhận indicator "+args[1]))
		}

	case "/changepassword":
		msg := tgbotapi.NewMessage(chatID, "In progress")
		bot.Send(msg)
		//syntax: /changepassword <old_password> <new_password> <confirm_newpassword>
		// if len(args) < 3 {
		// 	msg := tgbotapi.NewMessage(chatID, "Usage: /changepassword <old_password> <new_password> <confirm_newpassword>")
		// 	bot.Send(msg)
		// 	return
		// }
		// old_password := args[0]
		// new_password := args[1]
		// confirm_newpassword := args[2]

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
	case "/create_snooze":
		if len(args) != 5 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /create_snooze <spot/future> <symbol> <conditionType> <startTime> <endTime>")
			bot.Send(msg)
			return
		}
		price_type := args[0]
		symbol := args[1]
		conditionType := args[2]
		startTime := args[3]
		endTime := args[4]

		// Validate time format
		layout := "2006-01-02T15:04:05"
		_, err := time.Parse(layout, startTime)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Invalid startTime format. Please use format: YYYY-MM-DDThh:mm:ss")
			bot.Send(msg)
			return
		}

		_, err = time.Parse(layout, endTime)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Invalid endTime format. Please use format: YYYY-MM-DDThh:mm:ss")
			bot.Send(msg)
			return
		}

		go CreateSnoozeTrigger(chatID, bot, price_type, symbol, conditionType, startTime, endTime)
	case "/delete_snooze":
		if len(args) != 2 {
			msg := tgbotapi.NewMessage(chatID, "Usage: /delete_snooze <spot/future> <symbol>")
			bot.Send(msg)
			return
		}
		price_type := args[0]
		symbol := args[1]
		DeleteSnoozeTrigger(chatID, bot, symbol, price_type, true)
	}
}

func HandleKlineCommand(chatID int64, command string, bot *tgbotapi.BotAPI, user *tgbotapi.User) {
	// Get the mute status of the user
	isMuted, err := services.GetMute(int(user.ID))
	if err != nil {
		log.Println("Error getting mute status:", err)
	}
	if isMuted && !strings.Contains(command, "/mute") {
		return
	}
	switch command {
	case "/kline":
		updateSymbolUsage("BTCUSDT")
		updateSymbolUsage("ETHUSDT")
		updateSymbolUsage("BNBUSDT")

		msg := tgbotapi.NewMessage(chatID, "Choose fetch type:")
		fetchTypes := []string{"ondemand", "realtime"}
		var rows []tgbotapi.KeyboardButton
		for _, t := range fetchTypes {
			rows = append(rows, tgbotapi.NewKeyboardButton(t))
		}
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(rows)
		bot.Send(msg)
		mapMutex.Lock()
		UserSelections[chatID] = map[string]string{"step": "fetch_type_selection"}
		stopChanMap[chatID] = make(chan bool)
		mapMutex.Unlock()
	default:
		handleUserSteps(command, bot, chatID, user)
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

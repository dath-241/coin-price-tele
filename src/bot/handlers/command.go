package handlers

import (
	// "context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"telegram-bot/services"

	"time"

	// "github.com/chromedp/chromedp"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/snapshot-chromedp/render"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func init() {
	services.InitDB()
}

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
			msg := tgbotapi.NewMessage(chatID, "Usage: /delete_trigger <symbol> <spot/future/price-difference/funding-rate>")
			bot.Send(msg)
			return
		}
		symbol := args[1]
		price_type := args[0]
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

func sendScreamedMessage(message *tgbotapi.Message) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, strings.ToUpper(message.Text))
	msg.ParseMode = tgbotapi.ModeHTML
	return msg
}

func copyMessage(message *tgbotapi.Message) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	return msg
}

func getKlineData(symbol string, interval string, options ...int) ([]klineData, error) {
	apiURL := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=%s", symbol, interval)

	if len(options) > 0 {
		apiURL = fmt.Sprintf("%s&limit=%d", apiURL, options[0])
	}

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("request failed: %s", resp.Status)
	}

	var data [][]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &data)

	var kd []klineData
	for _, k := range data {
		openTime := int64(k[0].(float64)) / 1000
		date := time.Unix(openTime, 0).Format("2006-01-02 15:04:05")

		open, _ := parseFloat32(k[1].(string))
		high, _ := parseFloat32(k[2].(string))
		low, _ := parseFloat32(k[3].(string))
		close, _ := parseFloat32(k[4].(string))

		kd = append(kd, klineData{
			Date: date,
			Data: [4]float32{open, close, low, high},
		})
	}

	return kd, nil
}

func sendChartToTelegram(bot *tgbotapi.BotAPI, chatID int64, chart *charts.Kline) error {
	initialMsg := tgbotapi.NewMessage(chatID, "Uploading file...")
	sentMsg, err := bot.Send(initialMsg)
	if err != nil {
		log.Printf("failed to send initial message: %v", err)
		return fmt.Errorf("failed to send initial message: %w", err)
	}
	log.Println("Uploading file...")

	fileName := fmt.Sprintf("chart-%d%d.png", chatID, time.Now().UnixNano())
	log.Printf("File name generated: %s", fileName)

	// Render chart with Chromedp context
	err = render.MakeChartSnapshot(chart.RenderContent(), fileName)
	if err != nil {
		log.Printf("Failed to generate chart snapshot: %v", err)
		return fmt.Errorf("failed to generate chart snapshot: %w", err)
	}
	log.Printf("Chart snapshot generated: %s", fileName)

	imgBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("failed to read generated chart image: %v", err)
		return fmt.Errorf("failed to read generated chart image: %w", err)
	}
	log.Printf("Image read successfully, size: %d bytes", len(imgBytes))

	msg := tgbotapi.NewPhoto(chatID, tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: imgBytes,
	})

	if _, err := bot.Send(msg); err != nil {
		log.Printf("failed to send chart image: %v", err)
		return fmt.Errorf("failed to send chart image: %w", err)
	}

	if err := os.Remove(fileName); err != nil {
		log.Printf("warning: failed to delete image file %s: %v", fileName, err)
	}

	editMsg := tgbotapi.NewEditMessageText(chatID, sentMsg.MessageID, "File uploaded successfully!")
	if _, err := bot.Send(editMsg); err != nil {
		log.Printf("warning: failed to edit initial message: %v", err)
	}

	return nil
}

type AllTriggerResponse struct {
	ID                       string  `json:"id"`
	AlertID                  string  `json:"alert_id"`
	Username                 string  `json:"username"`
	Symbol                   string  `json:"symbol"`
	Condition                string  `json:"condition"`
	SpotPriceThreshold       float64 `json:"spotPriceThreshold"`
	FuturePriceThreshold     float64 `json:"futurePriceThreshold"`
	PriceDifferenceThreshold float64 `json:"priceDifferenceThreshold"`
	FundingRateThreshold     float64 `json:"fundingRateThreshold"`
}

func GetAllTrigger(ID int64, bot *tgbotapi.BotAPI) {
	url := "https://hcmutssps.id.vn/api/vip2/get/alerts"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cookie", fmt.Sprintf("token=%s", token))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

	var response []AllTriggerResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}

	// Format the response for sending
	var responseText string
	count := 1
	for _, trigger := range response {
		// if trigger.SpotPriceThreshold != 0 {
		// 	responseText += fmt.Sprintf("%d.\n\tSymbol: %s\n\tCondition: %s\n\tspotPriceThreshold: %f\n",
		// 		count, trigger.Symbol, trigger.Condition, trigger.SpotPriceThreshold)
		// } else if trigger.FuturePriceThreshold != 0 {
		// 	responseText += fmt.Sprintf("%d.\n\tSymbol: %s\n\tCondition: %s\n\tfuturePriceThreshold: %f\n",
		// 		count, trigger.Symbol, trigger.Condition, trigger.FuturePriceThreshold)
		// } else if trigger.PriceDifferenceThreshold != 0 {
		// 	responseText += fmt.Sprintf("%d.\n\tSymbol: %s\n\tCondition: %s\n\tpriceDifferenceThreshold: %f\n",
		// 		count, trigger.Symbol, trigger.Condition, trigger.PriceDifferenceThreshold)
		// } else if trigger.FundingRateThreshold != 0 {
		// 	responseText += fmt.Sprintf("%d.\n\tSymbol: %s\n\tCondition: %s\n\tfundingRateThreshold: %f\n",
		// 		count, trigger.Symbol, trigger.Condition, trigger.FundingRateThreshold)
		// }
		responseText += fmt.Sprintf("%d.\n\tSymbol: %s\n\tCondition: %s\n\tspotPriceThreshold: %f\n\tfuturePriceThreshold: %f\n\tpriceDifferenceThreshold: %f\n\tfundingRateThreshold: %f\n",
			count, trigger.Symbol, trigger.Condition, trigger.SpotPriceThreshold, trigger.FuturePriceThreshold, trigger.PriceDifferenceThreshold, trigger.FundingRateThreshold)
		count++
	}

	bot.Send(tgbotapi.NewMessage(ID, fmt.Sprintf("All triggers:\n%v", responseText)))
}

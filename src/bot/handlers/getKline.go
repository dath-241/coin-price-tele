package handlers

import (
	"bufio"
	"io"
	"sort"

	//"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//"strconv"
	"os"
	"strings"
	"sync"
	"time"

	"telegram-bot/services"

	// "github.com/chromedp/chromedp"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/snapshot-chromedp/render"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type KlineData struct {
	Symbol             string `json:"symbol"`
	EventTime          string `json:"eventTime"`
	KlineStartTime     string `json:"klineStartTime"`
	KlineCloseTime     string `json:"klineCloseTime"`
	OpenPrice          string `json:"openPrice"`
	ClosePrice         string `json:"closePrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	NumberOfTrades     int    `json:"numberOfTrades"`
	BaseAssetVolume    string `json:"baseAssetVolume"`
	TakerBuyVolume     string `json:"takerBuyVolume"`
	TakerBuyBaseVolume string `json:"takerBuyBaseVolume"`
	Volume             string `json:"volume"`
}

const baseUrl = "https://a2-price.thuanle.me//api/vip1/get-kline"

// UserConnection stores request state for each user
type UserConnection struct {
	isStreaming bool
}

// Tạo map để lưu trạng thái người dùng
var userConnections = make(map[int64]*UserConnection)
var mapMutex = sync.Mutex{}
var UserSelections = make(map[int64]map[string]string) // To track user choices
var stopChanMap = make(map[int64]chan bool)            // To manage stop signals for each user
var symbolUsage = make(map[string]int)                 // Track symbol usage for popularity
// Update symbol usage
func updateSymbolUsage(symbol string) {
	// mapMutex.Lock()
	// defer mapMutex.Unlock()
	symbolUsage[symbol]++
}

// Function to format JSON as Telegram code block
func formatJSONResponse(data interface{}) string {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Printf("Error formatting JSON: %v", err)
		return "Error formatting data"
	}

	// Wrap JSON in Telegram MarkdownV2 code block
	return fmt.Sprintf("```json\n%s\n```", jsonData)
}

// Fetch Kline Data
func fetchKlineDataRealtime(symbol, interval string, cookie string, chatID int64, bot *tgbotapi.BotAPI) {
	reqUrl := fmt.Sprintf("%s?symbols=%s&interval=%s", baseUrl, symbol, interval)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Error creating request: %v", err)))
		return
	}
	services.SetHeadersWithPrice(req, cookie)
	// req.Header.Set("Accept", "*/*")
	// req.Header.Set("Cookie", "token=eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJNSyIsInN1YiI6InRyYW5odXkiLCJwYXNzd29yZCI6ImFpIGNobyBjb2kgbeG6rXQga2jhuql1IiwiZXhwIjoxNzMyODUzNjE4fQ.D5MqbwKknk4ZkrGb6hvrceRRbkFdy7bTfCCNVMeg8jo")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Error fetching Kline data: %v", err)))
		return
	}
	defer resp.Body.Close()

	stopChan := stopChanMap[chatID]
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-stopChan:
			// bot.Send(tgbotapi.NewMessage(chatID, "Stopped fetching Kline data."))
			return
		default:
			if UserSelections[chatID]["isPaused"] == "true" {
				continue
			}
			line := strings.TrimPrefix(scanner.Text(), "data:")
			var klineData KlineData
			err := json.Unmarshal([]byte(line), &klineData)
			if err != nil {
				// log.Printf("Error decoding response: %v", err)
				continue
			}
			selectedData := map[string]interface{}{
				"symbol":     klineData.Symbol,
				"openPrice":  formatNumber(klineData.OpenPrice, false),
				"closePrice": formatNumber(klineData.ClosePrice, false),
				"highPrice":  formatNumber(klineData.HighPrice, false),
				"lowPrice":   formatNumber(klineData.LowPrice, false),
				"volume":     formatNumber(klineData.Volume, false),
				"eventTime":  klineData.EventTime,
				"tradeCount": klineData.NumberOfTrades,
			}

			// Format selected fields as JSON and send to Telegram
			jsonMessage := formatJSONResponse(selectedData)
			msg := tgbotapi.NewMessage(chatID, jsonMessage)
			msg.ParseMode = "MarkdownV2" // Use MarkdownV2 to display JSON properly
			bot.Send(msg)
			time.Sleep(time.Second)
		}
	}

	if err := scanner.Err(); err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Error reading response: %v", err)))
	}
}

// Get top N symbols
func getTopSymbols(n int) []string {
	// mapMutex.Lock()
	// defer mapMutex.Unlock()
	type kv struct {
		Key   string
		Value int
	}
	var sorted []kv
	for k, v := range symbolUsage {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})
	var topSymbols []string
	for i, kv := range sorted {
		if i >= n {
			break
		}
		topSymbols = append(topSymbols, kv.Key)
	}

	return topSymbols
}

// Handle "Other" input
func handleOtherInput(bot *tgbotapi.BotAPI, chatID int64, input string) {

	if isValidSymbol(input) {

		updateSymbolUsage(input)
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Symbol '%s' added to the list!", input))
		for key, value := range symbolUsage {
			fmt.Printf("Symbol: %s, Usage: %d\n", key, value)
		}
		bot.Send(msg)
		// print(input)
	} else {
		msg := tgbotapi.NewMessage(chatID, "Invalid symbol. Please try again.")
		bot.Send(msg)
	}
}

// Validate symbol (placeholder for actual validation logic)
func isValidSymbol(symbol string) bool {
	return len(symbol) > 0 // Replace with actual validation logic or API call
}

// Handle user steps
func handleUserSteps(update string, bot *tgbotapi.BotAPI, chatID int64, user *tgbotapi.User) {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	if selection, ok := UserSelections[chatID]; ok {
		switch selection["step"] {
		case "fetch_type_selection":
			fetchType := update
			if fetchType == "ondemand" || fetchType == "realtime" {
				UserSelections[chatID]["fetchType"] = fetchType
				UserSelections[chatID]["step"] = "coin_selection"
				topSymbols := getTopSymbols(3)
				topSymbols = append(topSymbols, "Other")
				msg := tgbotapi.NewMessage(chatID, "Select the coin:")
				var rows []tgbotapi.KeyboardButton
				for _, symbol := range topSymbols {
					rows = append(rows, tgbotapi.NewKeyboardButton(symbol))
				}
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(rows)
				bot.Send(msg)
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "Invalid fetch type. Please choose 'ondemand' or 'realtime'."))
			}

		case "coin_selection":
			symbol := update
			if symbol == "Other" {
				bot.Send(tgbotapi.NewMessage(chatID, "Please enter the symbol:"))
				UserSelections[chatID]["step"] = "other_input"
			} else {
				UserSelections[chatID]["coin"] = symbol
				UserSelections[chatID]["step"] = "interval_selection"

				msg := tgbotapi.NewMessage(chatID, "Choose the interval:")
				intervals := []string{"1m", "5m", "1h", "1d"}
				var rows []tgbotapi.KeyboardButton
				for _, interval := range intervals {
					rows = append(rows, tgbotapi.NewKeyboardButton(interval))
				}
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(rows)
				bot.Send(msg)
			}

		case "other_input":
			symbol := update
			handleOtherInput(bot, chatID, symbol)
			UserSelections[chatID]["coin"] = symbol
			UserSelections[chatID]["step"] = "interval_selection"

			msg := tgbotapi.NewMessage(chatID, "Choose the interval:")
			intervals := []string{"1m", "5m", "1h", "1d"}
			var rows []tgbotapi.KeyboardButton
			for _, interval := range intervals {
				rows = append(rows, tgbotapi.NewKeyboardButton(interval))
			}
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(rows)
			bot.Send(msg)

		case "interval_selection":
			interval := update
			symbol := UserSelections[chatID]["coin"]
			UserSelections[chatID]["interval"] = interval

			if UserSelections[chatID]["fetchType"] == "ondemand" {
				limit := 50
				data, err := getKlineData(symbol, interval, limit) // Pass parameters as needed
				if err != nil {
					_, _ = bot.Send(tgbotapi.NewMessage(chatID, "Error fetching Kline data: "+err.Error()))
				} else {
					// log.Println(data)
					sendChartToTelegram(bot, chatID, klineBase(data, symbol, interval))
				}
				updateSymbolUsage(symbol)
				return
			} else {
				UserSelections[chatID]["step"] = "fetching_data"
				UserSelections[chatID]["isFetching"] = "true"
				UserSelections[chatID]["isPaused"] = "false"

				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Start fetching data for %s (%s).", symbol, interval))
				buttons := []string{"Resume", "Stop", "Chart"}
				var rows []tgbotapi.KeyboardButton
				for _, btn := range buttons {
					rows = append(rows, tgbotapi.NewKeyboardButton(btn))
				}
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(rows)
				bot.Send(msg)
				// mapMutex.Lock()
				userConnections[chatID] = &UserConnection{isStreaming: true}
				// mapMutex.Unlock()
				token, err := services.GetUserToken(int(user.ID))
				if err != nil {
					log.Println("Error retrieving token:", err)
					return
				}
				if UserSelections[chatID]["fetchType"] == "realtime" {
					go fetchKlineDataRealtime(symbol, interval, token, chatID, bot)
				}

				updateSymbolUsage(symbol)
			}
		case "fetching_data":
			handleFetchingActions(update, bot, chatID)
			log.Print(update)

		default:
			bot.Send(tgbotapi.NewMessage(chatID, "Unknown step. Please restart the process."))
		}
	}
}

// Handle fetching actions
func handleFetchingActions(update string, bot *tgbotapi.BotAPI, chatID int64) {
	action := update
	stopChan := stopChanMap[chatID]

	switch action {
	case "Resume":
		if UserSelections[chatID]["isFetching"] == "true" {
			UserSelections[chatID]["isPaused"] = "true"
			UserSelections[chatID]["isFetching"] = "false"
		} else {
			UserSelections[chatID]["isPaused"] = "false"
			UserSelections[chatID]["isFetching"] = "true"
		}
		// go fetchKlineData(symbol, interval, chatID, bot)
	case "Stop":
		if stopChan != nil {
			close(stopChan)
			delete(stopChanMap, chatID)
			// bot.Send(tgbotapi.NewMessage(chatID, "Fetching data stopped."))
		}
		UserSelections[chatID]["step"] = ""
	case "Chart":
		symbol := UserSelections[chatID]["coin"]
		interval := UserSelections[chatID]["interval"]
		limit := 50
		data, err := getKlineData(symbol, interval, limit) // Pass parameters as needed
		if err != nil {
			_, _ = bot.Send(tgbotapi.NewMessage(chatID, "Error fetching Kline data: "+err.Error()))
		} else {
			// log.Println(data)
			sendChartToTelegram(bot, chatID, klineBase(data, symbol, interval))
		}
	default:
		if stopChan != nil {
			close(stopChan)
			delete(stopChanMap, chatID)
			// bot.Send(tgbotapi.NewMessage(chatID, "Fetching data stopped."))
		}
		UserSelections[chatID]["step"] = ""
		fakeMessage := &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: chatID,
			},
			From: &tgbotapi.User{
				FirstName: "System",
				LastName:  "Bot",
			},
			Text: update,
		}

		HandleMessage(fakeMessage, bot)
	}
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
	body, _ := io.ReadAll(resp.Body)
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

	// imgBytes, err := ioutil.ReadFile(fileName)
	// if err != nil {
	// 	log.Printf("failed to read generated chart image: %v", err)
	// 	return fmt.Errorf("failed to read generated chart image: %w", err)
	// }
	// log.Printf("Image read successfully, size: %d bytes", len(imgBytes))

	msg := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(fileName))

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

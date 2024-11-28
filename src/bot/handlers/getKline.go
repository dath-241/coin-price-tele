package handlers

import (
	"bufio"
	"io/ioutil"

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

const baseUrl = "https://dath.hcmutssps.id.vn/api/vip1/get-kline"

// UserConnection stores request state for each user
type UserConnection struct {
	isStreaming bool
}

// Tạo map để lưu trạng thái người dùng
var userConnections = make(map[int64]*UserConnection)
var mapMutex = sync.Mutex{}

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

// fetchKlineData sends a GET request to the backend API with cookie for security
func fetchKlineData(symbol, interval, cookie string, chatID int64, bot *tgbotapi.BotAPI) {
	reqUrl := fmt.Sprintf("%s?symbols=%s&interval=%s", baseUrl, symbol, interval)
	log.Printf("API URL: %s", reqUrl)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Request creation error: %v", err)))
		return
	}
	// req.Header.Set("Accept", "*/*")
	// req.Header.Set("Cookie", cookie)
	services.SetHeadersWithPrice(req, cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to fetch Kline data: %v", err)))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Received status code %d", resp.StatusCode)))
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	var line string
	for scanner.Scan() {
		mapMutex.Lock()
		userConn := userConnections[chatID]
		if userConn == nil || !userConn.isStreaming {
			mapMutex.Unlock()
			return // Thoát vòng lặp khi isStreaming = false
		}
		mapMutex.Unlock()
		// Lấy dòng dữ liệu và loại bỏ tiền tố "data:"
		line = strings.TrimPrefix(scanner.Text(), "data:")

		// Giải mã JSON
		var klineData KlineData
		err := json.Unmarshal([]byte(line), &klineData)
		if err != nil {
			// bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Error decoding response: %v", err)))
			continue
		}

		// Gửi dữ liệu đã giải mã đến người dùng Telegram
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

		jsonMessage := formatJSONResponse(selectedData)
		msg := tgbotapi.NewMessage(chatID, jsonMessage)
		msg.ParseMode = "MarkdownV2"
		bot.Send(msg)
		time.Sleep(time.Second)
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

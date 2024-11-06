package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"telegram-bot/services"
	"time"
	"os"
	"path/filepath"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	APIBaseURL_Spot_Price    = "https://hcmutssps.id.vn/api/get-spot-price"
	APIBaseURL_Futures_Price = "https://hcmutssps.id.vn/api/get-future-price"
	APIBaseURL_Funding_Rate  = "https://hcmutssps.id.vn/api/get-funding-rate"
)

const (
	SpotExchangeInfoURL    = "https://api.binance.com/api/v3/exchangeInfo"
	FuturesExchangeInfoURL = "https://fapi.binance.com/fapi/v1/exchangeInfo"
)

type ExchangeInfo struct {
	Symbols []struct {
		Symbol string `json:"symbol"`
	} `json:"symbols"`
}

type ErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SpotPriceResponse struct {
	Price     string `json:"price"`
	EventTime string `json:"eventTime"`
	Symbol    string `json:"symbol"`
}

type FuturesPriceResponse struct {
	Price     string `json:"price"`
	EventTime string `json:"eventTime"`
	Symbol    string `json:"symbol"`
}

type FundingRateResponse struct {
	Symbol                   string `json:"symbol"`
	FundingRate              string `json:"fundingRate"`
	FundingRateCountdown     string `json:"fundingCountdown"`
	EventTime                string `json:"eventTime"`
	AdjustedFundingRateCap   string `json:"adjustedFundingRateCap"`
	AdjustedFundingRateFloor string `json:"adjustedFundingRateFloor"`
	FundingIntervalHours     int    `json:"fundingIntervalHours"`
}
type PriceInfoSpot struct {
	Symbol    string `json:"Symbol"`
	EventTime string `json:"Event time"`
	SpotPrice string `json:"Spot price"`
}

type PriceInfoFutures struct {
	Symbol    string `json:"Symbol"`
	EventTime string `json:"Event time"`
	SpotPrice string `json:"Futures price"`
}

type PriceInfoFundingRate struct {
	Symbol                   string `json:"Symbol"`
	EventTime                string `json:"Event time"`
	FundingRate              string `json:"Funding rate"`
	FundingRateCountdown     string `json:"Time until next funding"`
	FundingRateIntervalHours string `json:"Funding rate interval"`
}

func formatPrice(input string) string {

	parts := strings.Split(input, ".")

	intPart := parts[0]
	n := len(intPart)
	if n <= 3 {
		return input
	}

	var result strings.Builder
	offset := n % 3
	if offset > 0 {
		result.WriteString(intPart[:offset])
		if n > 3 {
			result.WriteString(",")
		}
	}
	for i := offset; i < n; i += 3 {
		result.WriteString(intPart[i : i+3])
		if i+3 < n {
			result.WriteString(",")
		}
	}

	if len(parts) > 1 {
		result.WriteString(".")
		result.WriteString(parts[1])
	}

	return result.String()
}

func intToString(n int) string {
	return strconv.Itoa(n)
}

func FormatPrice1(a string) string {

	for i := len(a) - 1; i >= 0; i-- {
		if a[i] != '0' {
			if a[i] == '.' {
				return a + "00"
			}
			return a
		}
		a = a[:i]
	}

	return a
}

// test ( api BE error)
// func GetSpotPriceStream(chatID int64, symbol string, bot *tgbotapi.BotAPI, token string) {
// 	bot.Send(tgbotapi.NewMessage(chatID, "spot price"))
// }

// func GetFuturesPriceStream(chatID int64, symbol string, bot *tgbotapi.BotAPI, token string) {
// 	bot.Send(tgbotapi.NewMessage(chatID, "futures price"))
// }

// func GetFundingRateStream(chatID int64, symbol string, bot *tgbotapi.BotAPI, token string) {
// 	bot.Send(tgbotapi.NewMessage(chatID, "funding rate"))
// }

func GetSpotPriceStream(chatID int64, symbol string, bot *tgbotapi.BotAPI, token string) {

	// Create a cancellable context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure context is canceled when done

	// Create the request URL
	reqUrl := fmt.Sprintf("%s?symbols=%s", APIBaseURL_Spot_Price, symbol)
	//log.Printf("API URL: %s", reqUrl)

	// Create an HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		log.Printf("Request creation error: %v", err)
		return
	}

	// Set necessary headers for the request
	// req.Header.Set("Accept", "*/*")
	// req.Header.Set("Cookie", fmt.Sprintf("token=%s", CookieToken))

	services.SetHeadersWithPrice(req, token)

	// Create an HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch spot price: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check the status code of the response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Received status code %d", resp.StatusCode)
		if resp.StatusCode == 500 {
			errorMsg := ErrorMessage{
				Code:    "500",
				Message: "You need to authenticate before executing this command.",
			}
			jsonMsg, err := json.MarshalIndent(errorMsg, "", "  ")
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
			}
			formattedMsg := fmt.Sprintf("```json\n%s\n```", string(jsonMsg))
			msg := tgbotapi.NewMessage(chatID, formattedMsg)
			msg.ParseMode = "MarkdownV2"
			bot.Send(msg)
		}
		if resp.StatusCode == 404 {
			errorMsg := ErrorMessage{
				Code:    "404",
				Message: "Symbol is not available.",
			}
			jsonMsg, err := json.MarshalIndent(errorMsg, "", "  ")
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
			}
			formattedMsg := fmt.Sprintf("```json\n%s\n```", string(jsonMsg))
			msg := tgbotapi.NewMessage(chatID, formattedMsg)
			msg.ParseMode = "MarkdownV2"
			bot.Send(msg)
		}
		return
	}

	// Read data from the stream
	scanner := bufio.NewScanner(resp.Body)
	var line string
	for scanner.Scan() {
		// Remove the "data:" prefix from the line
		line = strings.TrimPrefix(scanner.Text(), "data:")

		// Decode JSON
		var spotPriceResponse SpotPriceResponse
		err := json.Unmarshal([]byte(line), &spotPriceResponse)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			continue // Skip the error and continue reading the next data
		}
		pricestr := FormatPrice1(spotPriceResponse.Price)
		//log.Printf("Price: %s", pricestr)
		// price, err := strconv.ParseFloat(pricestr, 64)
		// if err != nil {
		// 	log.Printf("Error converting price to float: %v", err)
		// 	continue
		// }

		// Send decoded data to Telegram user and exit
		if strings.EqualFold(spotPriceResponse.Symbol, symbol) {
			formattedPrice := formatPrice(pricestr)

			priceInfo := PriceInfoSpot{
				Symbol:    spotPriceResponse.Symbol,
				EventTime: spotPriceResponse.EventTime,
				SpotPrice: formattedPrice,
			}
			// Convert the object to a JSON string
			jsonData, err := json.MarshalIndent(priceInfo, "", "    ")
			if err != nil {
				log.Printf("Error creating JSON: %v", err)
				//bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Error creating JSON: %v", err)))
				return
			}

			// Use HTML to display the JSON string
			msg := tgbotapi.NewMessage(chatID, "<pre>"+string(jsonData)+"</pre>")
			msg.ParseMode = "HTML"
			bot.Send(msg)

			cancel()
			return // Exit immediately after sending the first data
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading stream: %v", err)
		//bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Symbol is not available. Please provide a valid symbol.")))
	}
}

func GetFuturesPriceStream(chatID int64, symbol string, bot *tgbotapi.BotAPI, token string) {

	// Create a cancellable context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // Ensure context is canceled when done

	// Create the request URL
	reqUrl := fmt.Sprintf("%s?symbols=%s", APIBaseURL_Futures_Price, symbol)
	//log.Printf("API URL: %s", reqUrl)

	// Create an HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		log.Printf("Request creation error: %v", err)
		return
	}

	// Set necessary headers for the request
	// req.Header.Set("Accept", "*/*")
	// req.Header.Set("Cookie", fmt.Sprintf("token=%s", CookieToken))

	services.SetHeadersWithPrice(req, token)

	// Create an HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch futures price: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check the status code of the response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Received status code %d", resp.StatusCode)
		if resp.StatusCode == 500 {
			errorMsg := ErrorMessage{
				Code:    "500",
				Message: "You need to authenticate before executing this command.",
			}
			jsonMsg, err := json.MarshalIndent(errorMsg, "", "  ")
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
			}
			formattedMsg := fmt.Sprintf("```json\n%s\n```", string(jsonMsg))
			msg := tgbotapi.NewMessage(chatID, formattedMsg)
			msg.ParseMode = "MarkdownV2"
			bot.Send(msg)
		}
		if resp.StatusCode == 404 {
			errorMsg := ErrorMessage{
				Code:    "404",
				Message: "Symbol is not available.",
			}
			jsonMsg, err := json.MarshalIndent(errorMsg, "", "  ")
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
			}
			formattedMsg := fmt.Sprintf("```json\n%s\n```", string(jsonMsg))
			msg := tgbotapi.NewMessage(chatID, formattedMsg)
			msg.ParseMode = "MarkdownV2"
			bot.Send(msg)
		}
		return
	}

	// Read data from the stream
	scanner := bufio.NewScanner(resp.Body)
	var line string
	for scanner.Scan() {
		// Remove the "data:" prefix from the line
		line = strings.TrimPrefix(scanner.Text(), "data:")

		// Decode JSON
		var futuresPriceResponse FuturesPriceResponse
		err := json.Unmarshal([]byte(line), &futuresPriceResponse)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			continue // Skip the error and continue reading the next data
		}
		pricestr := FormatPrice1(futuresPriceResponse.Price)

		// price, err := strconv.ParseFloat(pricestr, 64)
		// if err != nil {
		// 	log.Printf("Error converting price to float: %v", err)
		// 	continue
		// }

		// Send decoded data to Telegram user and exit
		if strings.EqualFold(futuresPriceResponse.Symbol, symbol) {
			formattedPrice := formatPrice(pricestr)

			priceInfo := PriceInfoFutures{
				Symbol:    futuresPriceResponse.Symbol,
				EventTime: futuresPriceResponse.EventTime,
				SpotPrice: formattedPrice,
			}
			// Convert the object to a JSON string
			jsonData, err := json.MarshalIndent(priceInfo, "", "    ")
			if err != nil {
				log.Printf("Error creating JSON: %v", err)
				//bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Error creating JSON: %v", err)))
				return
			}

			// Use HTML to display the JSON string
			msg := tgbotapi.NewMessage(chatID, "<pre>"+string(jsonData)+"</pre>")
			msg.ParseMode = "HTML"
			bot.Send(msg)

			cancel()
			return // Exit immediately after sending the first data
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading stream: %v", err)
		//bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Symbol is not available. Please provide a valid symbol.")))
	}
}

func GetFundingRateStream(chatID int64, symbol string, bot *tgbotapi.BotAPI, token string) {
	// Create a cancellable context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure context is canceled when done

	// Create the request URL
	reqUrl := fmt.Sprintf("%s?symbols=%s", APIBaseURL_Funding_Rate, symbol)
	//log.Printf("API URL: %s", reqUrl)

	// Create an HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		log.Printf("Request creation error: %v", err)
		return
	}

	// Set necessary headers for the request
	// req.Header.Set("Accept", "*/*")
	// req.Header.Set("Cookie", fmt.Sprintf("token=%s", CookieToken))

	services.SetHeadersWithPrice(req, token)

	// Create an HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch funding rate: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check the status code of the response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Received status code %d", resp.StatusCode)
		if resp.StatusCode == 500 {
			errorMsg := ErrorMessage{
				Code:    "500",
				Message: "You need to authenticate before executing this command.",
			}
			jsonMsg, err := json.MarshalIndent(errorMsg, "", "  ")
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
			}
			formattedMsg := fmt.Sprintf("```json\n%s\n```", string(jsonMsg))
			msg := tgbotapi.NewMessage(chatID, formattedMsg)
			msg.ParseMode = "MarkdownV2"
			bot.Send(msg)
		}
		if resp.StatusCode == 404 {
			errorMsg := ErrorMessage{
				Code:    "404",
				Message: "Symbol is not available.",
			}
			jsonMsg, err := json.MarshalIndent(errorMsg, "", "  ")
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
			}
			formattedMsg := fmt.Sprintf("```json\n%s\n```", string(jsonMsg))
			msg := tgbotapi.NewMessage(chatID, formattedMsg)
			msg.ParseMode = "MarkdownV2"
			bot.Send(msg)
		}
		return
	}

	// Read data from the stream
	scanner := bufio.NewScanner(resp.Body)
	var line string
	for scanner.Scan() {
		// Remove the "data:" prefix from the line
		line = strings.TrimPrefix(scanner.Text(), "data:")

		// Decode JSON
		var fundingRateResponse FundingRateResponse
		err := json.Unmarshal([]byte(line), &fundingRateResponse)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			continue // Skip the error and continue reading the next data
		}
		fundingstr := FormatPrice1(fundingRateResponse.FundingRate)
		// Send decoded data to Telegram user and exit
		if strings.EqualFold(fundingRateResponse.Symbol, symbol) {
			fundingRateInterval := intToString(fundingRateResponse.FundingIntervalHours)
			priceInfo := PriceInfoFundingRate{
				Symbol:                   fundingRateResponse.Symbol,
				EventTime:                fundingRateResponse.EventTime,
				FundingRate:              fundingstr,
				FundingRateCountdown:     fundingRateResponse.FundingRateCountdown,
				FundingRateIntervalHours: fundingRateInterval + " hours",
			}

			// Convert the object to a JSON string
			jsonData, err := json.MarshalIndent(priceInfo, "", "    ")
			if err != nil {
				log.Printf("Error creating JSON: %v", err)
				//bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Error creating JSON: %v", err)))
				return
			}

			// Use HTML to display the JSON string
			msg := tgbotapi.NewMessage(chatID, "<pre>"+string(jsonData)+"</pre>")
			msg.ParseMode = "HTML"
			bot.Send(msg)

			cancel()
			return // Exit immediately after sending the first data
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading stream: %v", err)
		//bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Symbol is not available. Please provide a valid symbol.")))
	}
}

// Function to get available symbols from Binance API
// func GetAvailableSymbols(exchangeInfoURL string) ([]string, error) {
// 	resp, err := http.Get(exchangeInfoURL)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var exchangeInfo ExchangeInfo
// 	if err := json.Unmarshal(body, &exchangeInfo); err != nil {
// 		return nil, err
// 	}

// 	var symbols []string
// 	for _, symbol := range exchangeInfo.Symbols {
// 		//log.Printf("Symbol: %s", symbol.Symbol)
// 		symbols = append(symbols, symbol.Symbol)
// 	}
// 	return symbols, nil
// }

func FindSpotSymbol(input string, ) string {
    suffixes := []string{"USDT", "USDC", "BTC"}
    
    // Đọc file symbols
    currentDir, err := os.Getwd()
    if err != nil {
        log.Printf("Error getting current directory: %v", err)
        return ""
    }
    filePath := filepath.Join(currentDir, "services", "spot_symbols.txt")
    
    file, err := os.Open(filePath)
    if err != nil {
        log.Printf("Error opening spot_symbols file: %v", err)
        return ""
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    // Bỏ qua 2 dòng đầu tiên
    for i := 0; i < 2; i++ {
        scanner.Scan()
    }

    // Tìm chính xác symbol = input + suffix
    for _, suffix := range suffixes {
        targetSymbol := strings.ToUpper(input + suffix)
        
        // Reset scanner về đầu file sau mỗi suffix
        file.Seek(0, 0)
        scanner = bufio.NewScanner(file)
        // Bỏ qua 2 dòng đầu
        for i := 0; i < 2; i++ {
            scanner.Scan()
        }
        
        for scanner.Scan() {
            symbol := strings.TrimSpace(scanner.Text())
            if symbol == targetSymbol {
                return symbol
            }
        }
    }

    // Tìm symbol bắt đầu bằng input
    upperInput := strings.ToUpper(input)
    file.Seek(0, 0)
    scanner = bufio.NewScanner(file)
    // Bỏ qua 2 dòng đầu
    for i := 0; i < 2; i++ {
        scanner.Scan()
    }
    
    for scanner.Scan() {
        symbol := strings.TrimSpace(scanner.Text())
        if strings.HasPrefix(symbol, upperInput) {
            return symbol
        }
    }

    return ""
}


func FindFuturesSymbol(input string) string {
    suffixes := []string{"USDT", "USDC", "BTC"}
    
    // Đọc file symbols
    currentDir, err := os.Getwd()
    if err != nil {
        log.Printf("Error getting current directory: %v", err)
        return ""
    }
    filePath := filepath.Join(currentDir, "services", "futures_symbols.txt")
    
    file, err := os.Open(filePath)
    if err != nil {
        log.Printf("Error opening spot_symbols file: %v", err)
        return ""
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    // Bỏ qua 2 dòng đầu tiên
    for i := 0; i < 2; i++ {
        scanner.Scan()
    }

    // Tìm chính xác symbol = input + suffix
    for _, suffix := range suffixes {
        targetSymbol := strings.ToUpper(input + suffix)
        
        // Reset scanner về đầu file sau mỗi suffix
        file.Seek(0, 0)
        scanner = bufio.NewScanner(file)
        // Bỏ qua 2 dòng đầu
        for i := 0; i < 2; i++ {
            scanner.Scan()
        }
        
        for scanner.Scan() {
            symbol := strings.TrimSpace(scanner.Text())
            if symbol == targetSymbol {
                return symbol
            }
        }
    }

    // Tìm symbol bắt đầu bằng input
    upperInput := strings.ToUpper(input)
    file.Seek(0, 0)
    scanner = bufio.NewScanner(file)
    // Bỏ qua 2 dòng đầu
    for i := 0; i < 2; i++ {
        scanner.Scan()
    }
    
    for scanner.Scan() {
        symbol := strings.TrimSpace(scanner.Text())
        if strings.HasPrefix(symbol, upperInput) {
            return symbol
        }
    }

    return ""
}

// Function to find the closest symbol
func FindClosestSymbol(input string, symbols []string) string {
	suffixes := []string{"USDT", "USDC", "BTC"}

	for _, suffix := range suffixes {
		targetSymbol := strings.ToUpper(input + suffix)
		//log.Printf("Target symbol: %s", targetSymbol)
		for _, symbol := range symbols {
			if targetSymbol == symbol {
				return symbol
			}
		}
	}

	for _, symbol := range symbols {
		if strings.Contains(strings.ToUpper(symbol), strings.ToUpper(input)) {
			return symbol
		}
	}

	return ""
}

func FindClosestSymbol1(input string, symbols []string) string {
	suffixes := []string{"", "USDT", "USDC", "BTC"}

	for _, suffix := range suffixes {
		targetSymbol := strings.ToUpper(input + suffix)
		//log.Printf("Target symbol: %s", targetSymbol)
		for _, symbol := range symbols {
			if targetSymbol == symbol {
				return symbol
			}
		}
	}

	return ""
}


// func GetSpotSymbols() {
//     // Đọc symbols từ file

// 	currentDir, err := os.Getwd()
//     if err != nil {
//         fmt.Printf("Error getting current directory: %v\n", err)
//         return
//     }
// 	filePath := filepath.Join(currentDir, "services", "spot_symbols.txt")

//     file, err := os.Open(filePath)
//     if err != nil {
//         log.Printf("Error opening spot_symbols file: %v", err)
//         return
//     }
//     defer file.Close()

//     scanner := bufio.NewScanner(file)
//     // Bỏ qua 2 dòng đầu tiên
//     for i := 0; i < 2; i++ {
//         scanner.Scan()
//     }
    
//     // Đọc các symbols từ dòng thứ 3
//     for scanner.Scan() {
//         symbol := scanner.Text()
//         if symbol != "" {
//             SpotSymbols = append(SpotSymbols, symbol)
//         }
//     }

//     if err := scanner.Err(); err != nil {
//         log.Printf("Error reading spot_symbols file: %v", err)
// 	}
// }

// func GetFuturesSymbols() {
// 	currentDir, err := os.Getwd()
//     if err != nil {
//         fmt.Printf("Error getting current directory: %v\n", err)
//         return
//     }
// 	filePath := filepath.Join(currentDir, "services", "futures_symbols.txt")

//     file, err := os.Open(filePath)
//     if err != nil {
//         log.Printf("Error opening futures_symbols file: %v", err)
//         return
//     }
//     defer file.Close()

//     scanner := bufio.NewScanner(file)
//     // Bỏ qua 2 dòng đầu tiên
//     for i := 0; i < 2; i++ {
//         scanner.Scan()
//     }
    
//     // Đọc các symbols từ dòng thứ 3
//     for scanner.Scan() {
//         symbol := scanner.Text()
//         if symbol != "" {
//             FuturesSymbols = append(FuturesSymbols, symbol)
//         }
//     }

//     if err := scanner.Err(); err != nil {
//         log.Printf("Error reading futures_symbols file: %v", err)
// 	}
// }

package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

func FetchSpotSymbols() {
	url := "https://api.binance.com/api/v3/exchangeInfo"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result struct {
		Symbols []struct {
			Symbol string `json:"symbol"`
			Status string `json:"status"`
		} `json:"symbols"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	var symbols []string

	for _, sym := range result.Symbols {
		if sym.Status == "TRADING" {
			symbols = append(symbols, sym.Symbol)
		}
	}

	sort.Strings(symbols)

	// Tạo nội dung để ghi vào file
	content := fmt.Sprintf("Total Spot Symbols: %d\n\n", len(symbols))
	for _, symbol := range symbols {
		content += symbol + "\n"
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}
	filePath := filepath.Join(currentDir, "services", "spot_symbols.txt")

	// Ghi vào file
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Printf("Symbols have been saved to %s\n", filePath)
}

func FetchFuturesSymbols() {
    url := "https://fapi.binance.com/fapi/v1/exchangeInfo"
    resp, err := http.Get(url)
    if err != nil {
        fmt.Printf("Error fetching futures data: %v\n", err)
        return
    }
    defer resp.Body.Close()

    var result struct {
        Symbols []struct {
            Symbol string `json:"symbol"`
            Status string `json:"status"`
        } `json:"symbols"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        fmt.Printf("Error decoding futures response: %v\n", err)
        return
    }

    var symbols []string
    
    for _, sym := range result.Symbols {
        if sym.Status == "TRADING" {
            symbols = append(symbols, sym.Symbol)
        }
    }

    sort.Strings(symbols)

    // Tạo nội dung để ghi vào file
    content := fmt.Sprintf("Total Futures Symbols: %d\n\n", len(symbols))
    for _, symbol := range symbols {
        content += symbol + "\n"
    }

    currentDir, err := os.Getwd()
    if err != nil {
        fmt.Printf("Error getting current directory: %v\n", err)
        return
    }
    filePath := filepath.Join(currentDir, "services", "futures_symbols.txt")

    // Ghi vào file
    err = os.WriteFile(filePath, []byte(content), 0644)
    if err != nil {
        fmt.Printf("Error writing to file: %v\n", err)
        return
    }

    fmt.Printf("Futures symbols have been saved to %s\n", filePath)
}

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/joho/godotenv"
)

func validateBotResponse(client *telegram.Client, botUsername string, command string, expectedResponse string, timeout time.Duration) error {
	// Send message to bot
	msg, err := client.SendMessage(botUsername, command)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	// Wait for the bot to respond
	time.Sleep(timeout)

	// Get bot's response
	res, err := client.GetMessageByID(botUsername, msg.ID+1)
	if err != nil {
		return fmt.Errorf("failed to get response: %v", err)
	}

	// Validate response
	if !strings.Contains(res.MessageText(), expectedResponse) {
		return fmt.Errorf("unexpected response: got %q, want %q", res.MessageText(), expectedResponse)
	}

	return nil
}

func main() {
	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	appID := os.Getenv("APP_ID")
	appHash := os.Getenv("APP_HASH")
	phoneNum := os.Getenv("PHONE_NUMBER")
	appIDInt, _ := strconv.Atoi(appID)
	stringSession := os.Getenv("STRING_SESSION")
	client, _ := telegram.NewClient(telegram.ClientConfig{
		AppID:         int32(appIDInt),
		AppHash:       appHash,
		LogLevel:      telegram.LogInfo,
		StringSession: stringSession, // Uncomment this line to use string session
	})

	if err := client.Connect(); err != nil {
		panic(err)
	}

	if _, err := client.Login(phoneNum); err != nil {
		panic(err)
	}
	// Test bot responses
	tests := []struct {
		command          string
		expectedResponse string
		timeout          time.Duration
	}{
		{
			command:          "/login tranhuy coinprice123",
			expectedResponse: "Đăng nhập thành công",
			timeout:          5 * time.Second,
		},
		{
			command:          "/start",
			expectedResponse: "You are authenticated", // Replace with your expected response
			timeout:          1 * time.Second,
		},
		{
			command:          "/getinfo",
			expectedResponse: "🎉 Thông tin người dùng: 🎉",
			timeout:          2 * time.Second,
		},
		{
			command:          "/login",
			expectedResponse: "Usage: /login <username> <password>",
			timeout:          1 * time.Second,
		},
		{
			command:          "/all_triggers",
			expectedResponse: "🎯 All Triggers:",
			timeout:          1 * time.Second,
		},
		{
			command:          "/register",
			expectedResponse: "Usage: /register <email> <name> <username> <password>",
			timeout:          1 * time.Second,
		},
		{
			command:          "/register " + gofakeit.Email() + " " + gofakeit.Name() + " " + gofakeit.Username() + " " + gofakeit.Password(true, true, true, true, false, 10),
			expectedResponse: "Đăng kí tài khoản thành công",
			timeout:          1 * time.Second,
		},
		{
			command:          "/p BTC",
			expectedResponse: "Please select the information you want to view:",
			timeout:          1 * time.Second,
		},
		{
			command:          "/price_difference above BTCUSDT 100000",
			expectedResponse: "Registered price-difference of BTCUSDT above 100000 threshold successfully!",
			timeout:          1 * time.Second,
		},
		{
			command:          "/funding_rate_change",
			expectedResponse: "Usage: /funding_rate_change <lower/above> <symbol> <threshold>",
			timeout:          1 * time.Second,
		},
		{
			command:          "/alert_price_with_threshold",
			expectedResponse: "Usage: /alert_price_with_threshold <spot/future> <lower/above> <symbol> <threshold>",
			timeout:          2 * time.Second,
		},
		{
			command:          "/create_snooze",
			expectedResponse: "Usage: /create_snooze <spot/future> <symbol> <conditionType> <startTime> <endTime>",
			timeout:          2 * time.Second,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		err := validateBotResponse(client, "@clgslsm_bot", test.command, test.expectedResponse, test.timeout)
		if err != nil {
			fmt.Printf("Test failed for command %q: %v\n", test.command, err)
			os.Exit(1)
		} else {
			fmt.Printf("Test passed for command %q\n", test.command)
		}
	}
	fmt.Println("All tests passed")
}

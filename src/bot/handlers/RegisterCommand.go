package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"telegram-bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ErrorResponse struct {
	AlertID string `json:"alert_id"`
	Message string `json:"message"`
}

func RegisterPriceThreshold(ID int64, symbol string, threshold float64, is_lower bool, price_type string, bot *tgbotapi.BotAPI) error {
	url := fmt.Sprintf("https://a2-price.thuanle.me/api/vip2/create?triggerType=%s", price_type)
	fmt.Println("price_type:", price_type)
	method := "POST"

	condition := ">="
	if is_lower {
		condition = "<"
	}

	payload := strings.NewReader(fmt.Sprintf(`{
	  "symbol": "%s",
	  "price": %f,
	  "condition": "%s"
	}`, symbol, threshold, condition))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	token, err := services.GetUserToken(int(ID))
	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Cookie", fmt.Sprintf("token=%s", token))
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))

	var errorResponse ErrorResponse
	if err := json.Unmarshal(body, &errorResponse); err != nil {
		fmt.Println("Error unmarshalling response:", err)
		bot.Send(tgbotapi.NewMessage(ID, "Oops! Something went wrong. Please try again later."))
		return err
	}
	//bot.Send(tgbotapi.NewMessage(ID, errorResponse.Message))
	if errorResponse.AlertID != "" {
		if condition == "<" {
			bot.Send(tgbotapi.NewMessage(ID, fmt.Sprintf("Registered %s price of %s below %s threshold successfully!", price_type, symbol, removeTrailingZeros(threshold))))
		} else {
			bot.Send(tgbotapi.NewMessage(ID, fmt.Sprintf("Registered %s price of %s above %s threshold successfully!", price_type, symbol, removeTrailingZeros(threshold))))
		}
	}

	GetAllTrigger(ID, bot)
	return nil
}

func RegisterPriceDifferenceAndFundingRate(ID int64, symbol string, threshold float64, is_lower bool, Type string, bot *tgbotapi.BotAPI) error {
	url := fmt.Sprintf("https://a2-price.thuanle.me/api/vip2/create?triggerType=%s", Type)
	fmt.Println("Type:", Type)
	method := "POST"

	condition := ">="
	if is_lower {
		condition = "<"
	}

	var payload io.Reader
	if Type == "price-difference" {
		payload = strings.NewReader(fmt.Sprintf(`{
	"symbol": "%s",
	"condition": "%s",
	"priceDifference": %f
	}`, symbol, condition, threshold))
		fmt.Println("payload:", payload)
	} else if Type == "funding-rate" {
		payload = strings.NewReader(fmt.Sprintf(`{
	"symbol": "%s",
	"condition": "%s",
	"fundingRate": %f
	}`, symbol, condition, threshold))
		fmt.Println("payload:", payload)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")

	token, err := services.GetUserToken(int(ID))
	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Cookie", fmt.Sprintf("token=%s", token))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	var errorResponse ErrorResponse
	if err := json.Unmarshal(body, &errorResponse); err != nil {
		fmt.Println("Error unmarshalling response:", err)
		bot.Send(tgbotapi.NewMessage(ID, "Oops! Something went wrong. Please try again later."))
		return err
	}
	//bot.Send(tgbotapi.NewMessage(ID, errorResponse.Message))
	if errorResponse.AlertID != "" {
		if condition == "<" {
			bot.Send(tgbotapi.NewMessage(ID, fmt.Sprintf("Registered %s of %s below %s threshold successfully!", Type, symbol, removeTrailingZeros(threshold))))
		} else {
			bot.Send(tgbotapi.NewMessage(ID, fmt.Sprintf("Registered %s of %s above %s threshold successfully!", Type, symbol, removeTrailingZeros(threshold))))
		}
	}

	GetAllTrigger(ID, bot)
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
	TriggerType              string  `json:"triggerType"`
	TriggerTypeTmp           string  `json:"trigger_type"` // this fucking thing in wrong
}

func GetAllTrigger(ID int64, bot *tgbotapi.BotAPI) {
	url := "https://a2-price.thuanle.me/api/vip2/get/alerts"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "*/*")
	token, err := services.GetUserToken(int(ID))
	if err != nil {
		fmt.Println(err)
		return
	}
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
		bot.Send(tgbotapi.NewMessage(ID, "Oops! Something went wrong. Please try again later."))
		return
	}

	// Format the response for sending
	var responseText string
	count := 1
	for _, trigger := range response {
		// Add a divider between entries
		if count > 1 {
			responseText += "───────────────\n"
		}

		if trigger.TriggerType == "spot" {
			responseText += fmt.Sprintf("*%d.* 📊\n`Symbol:` *%s*\n`Condition:` _%s_\n`Spot Price:` *%s*\n",
				count, trigger.Symbol, trigger.Condition, removeTrailingZeros(trigger.SpotPriceThreshold))
		} else if trigger.TriggerType == "future" {
			responseText += fmt.Sprintf("*%d.* 🔮\n`Symbol:` *%s*\n`Condition:` _%s_\n`Future Price:` *%s*\n",
				count, trigger.Symbol, trigger.Condition, removeTrailingZeros(trigger.FuturePriceThreshold))
		} else if trigger.TriggerType == "price-difference" {
			responseText += fmt.Sprintf("*%d.* 📈\n`Symbol:` *%s*\n`Condition:` _%s_\n`Price Diff:` *%s*\n",
				count, trigger.Symbol, trigger.Condition, removeTrailingZeros(trigger.PriceDifferenceThreshold))
		} else if trigger.TriggerType == "funding-rate" {
			responseText += fmt.Sprintf("*%d.* 💰\n`Symbol:` *%s*\n`Condition:` _%s_\n`Funding Rate:` *%s*\n",
				count, trigger.Symbol, trigger.Condition, removeTrailingZeros(trigger.FundingRateThreshold))
		} else if trigger.TriggerTypeTmp == "future" {
			responseText += fmt.Sprintf("*%d.* 🔮\n`Symbol:` *%s*\n`Condition:` _%s_\n`Future Price:` *%s*\n",
				count, trigger.Symbol, trigger.Condition, removeTrailingZeros(trigger.FuturePriceThreshold))
		}
		count++
	}

	if responseText == "" {
		bot.Send(tgbotapi.NewMessage(ID, "❌ No triggers found"))
	} else {
		msg := tgbotapi.NewMessage(ID, fmt.Sprintf("🎯 *All Triggers:*\n\n%v", responseText))
		msg.ParseMode = "Markdown"
		bot.Send(msg)
	}
}

func DeleteTrigger(ID int64, bot *tgbotapi.BotAPI, symbol string, price_type string) {
	url := fmt.Sprintf("https://a2-price.thuanle.me/api/vip2/delete/%s?triggerType=%s", symbol, price_type)
	method := "DELETE"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		bot.Send(tgbotapi.NewMessage(ID, "Oops! Something went wrong. Please try again later."))
		return
	}
	req.Header.Add("Accept", "*/*")

	token, err := services.GetUserToken(int(ID))
	if err != nil {
		fmt.Println(err)
		return
	}
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
	bot.Send(tgbotapi.NewMessage(ID, string(body)))
	DeleteSnoozeTrigger(ID, bot, symbol, price_type, false)
	GetAllTrigger(ID, bot)
}

func removeTrailingZeros(n float64) string {
	// Convert float to string with maximum precision
	str := fmt.Sprintf("%f", n)

	// Remove trailing zeros after decimal point
	str = strings.TrimRight(strings.TrimRight(str, "0"), ".")

	return str
}

func CreateSnoozeTrigger(ID int64, bot *tgbotapi.BotAPI, price_type string, symbol string, conditionType string, startTime string, endTime string) {
	url := fmt.Sprintf("https://a2-price.thuanle.me/api/vip2/create/snooze?snoozeType=%s", price_type)
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
		"symbol": "%s",
		"conditionType": "%s",
		"startTime": "%s",
		"endTime": "%s"
	}`, symbol, conditionType, startTime, endTime))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	token, err := services.GetUserToken(int(ID))
	if err != nil {
		fmt.Println(err)
		return
	}
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
	bot.Send(tgbotapi.NewMessage(ID, string(body)))
}

func DeleteSnoozeTrigger(ID int64, bot *tgbotapi.BotAPI, symbol string, price_type string , alert_error bool) {
	url := fmt.Sprintf("https://a2-price.thuanle.me/api/vip2/delete/snooze/%s?snoozeType=%s", symbol, price_type)
	method := "DELETE"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	token, err := services.GetUserToken(int(ID))
	if err != nil {
		fmt.Println(err)
		return
	}
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
	if strings.Contains(string(body), "success") {
		bot.Send(tgbotapi.NewMessage(ID, "Snooze trigger deleted successfully!"))
	} else if (strings.Contains(string(body), "Error") && alert_error) {
		bot.Send(tgbotapi.NewMessage(ID, "Failed to delete snooze trigger."))
	}
}

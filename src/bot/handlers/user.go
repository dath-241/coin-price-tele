package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleUserInfo(chatID int64, bot *tgbotapi.BotAPI, response string) {
	// Parse chuỗi JSON vào struct UserInfo
	var userInfo UserInfo
	err := json.Unmarshal([]byte(response), &userInfo)
	if err != nil {
		log.Println("Error parsing user info:", err)
		_, _ = bot.Send(tgbotapi.NewMessage(chatID, "Đã xảy ra lỗi khi lấy thông tin user."))
		return
	}

	// Tạo thông báo đẹp hơn
	successMessage := fmt.Sprintf("🎉 *Thông tin người dùng:* 🎉\n\n"+
		"👤 *Tên*: %s\n"+
		"📧 *Email*: %s\n"+
		"👑 *VIP Role*: %d\n"+
		"💼 *Username*: %s\n"+
		"💰 *Coin*: %d\n",
		userInfo.Name, userInfo.Email, userInfo.VipRole, userInfo.Username, userInfo.Coin)

	// Gửi thông báo đã định dạng
	msg := tgbotapi.NewMessage(chatID, successMessage)
	msg.ParseMode = "Markdown"
	_, _ = bot.Send(msg)
}

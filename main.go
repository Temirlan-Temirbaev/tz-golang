package main

import (
	"github.com/Temirlan-Temirbaev/tz-golang/config"
	"github.com/Temirlan-Temirbaev/tz-golang/controllers"
	"github.com/Temirlan-Temirbaev/tz-golang/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}
	config.InitAuthConfig()
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{}, &models.Task{})

	go startTelegramBot()

	e := echo.New()
	controllers.InitTaskRoutes(e)
	controllers.InitUserRoutes(e)
	e.Logger.Fatal(e.Start(":3000"))
}

func startTelegramBot() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Printf("Failed to initialize bot: %v", err)
		log.Panic(err)
	}
	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome! Click the button below to open the Mini App.")
			inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
				[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonURL("Open Mini App", os.Getenv("TELEGRAM_WEBAPP_URL"))},
			)
			msg.ReplyMarkup = inlineKeyboard
			if _, err := bot.Send(msg); err != nil {
				log.Println("Error sending message:", err)
			}
		}
	}
}

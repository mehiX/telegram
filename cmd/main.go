package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mehix/telegram"
)

func main() {
	godotenv.Load()
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	tcli := telegram.NewClient(token, chatID)
	defer tcli.Stop()

	tcli.SendMessage(fmt.Sprintf("It is %v", time.Now()))

}

package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/mehix/telegram"
)

var (
	token   string
	chatIDs string
)

func init() {
	godotenv.Load()

	token = os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	flag.StringVar(&chatIDs, "to", chatID, "send messages to these destinations")

	flag.Parse()
}

func main() {

	done := make(chan struct{})

	tcli := telegram.NewClient(token, telegram.WithChatIDs(strings.Split(chatIDs, ",")...))

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := tcli.Stop(ctx); err != nil {
			slog.Error("stopping telegram client", "error", err.Error())
		}
		<-done
	}()

	go func() {
		defer close(done)
		for err := range tcli.Errors() {
			slog.Error(err.Error())
		}
	}()

	tcli.Send(fmt.Sprintf("It is %v", time.Now()))

}

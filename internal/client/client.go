package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

const maxMessageSize = 2048

var telegramApiURL = "https://api.telegram.org"

type message struct {
	chatID string
	txt    string
}

type TClient struct {
	Token      string
	HttpClient *http.Client

	ChatIdFn func() []string

	messages chan<- message
	errors   <-chan error

	done chan struct{}
}

func (c *TClient) Run() {

	msgs := make(chan message, 16)
	c.messages = msgs

	errs := make(chan error, 1)
	c.errors = errs

	c.done = make(chan struct{})

	go func() {
		defer close(errs)
		defer close(c.done)
		for m := range msgs {
			if err := c.send(m); err != nil {
				select {
				case errs <- err:
				default:
				}
			}
		}
	}()
}

// Stop waits for all the messages to be sent before returning
func (c *TClient) Stop(ctx context.Context) error {
	close(c.messages)

	select {
	case <-c.done:
		slog.Info("all telegram messages processed")
	case <-ctx.Done():
		slog.Warn("some telegram messages have been discarded")
		return ctx.Err()
	}

	return nil
}

func (c *TClient) Errors() <-chan error {
	return c.errors
}

// SendTo sends a message to the configured chat.
// Long messages are split in chunks that are then sent individually.
func (c *TClient) SendTo(msg string, chatIDs ...string) {
	parts := splitLongMessage(msg, maxMessageSize)
	for i := range parts {
		for _, chatID := range chatIDs {
			c.messages <- message{txt: parts[i], chatID: strings.TrimSpace(chatID)}
		}
	}
}

func (c *TClient) Send(msg string) {
	parts := splitLongMessage(msg, maxMessageSize)
	for i := range parts {
		for _, chatID := range c.ChatIdFn() {
			c.messages <- message{txt: parts[i], chatID: strings.TrimSpace(chatID)}
		}
	}
}

func splitLongMessage(msg string, limit int) []string {
	if len(msg) <= limit {
		return []string{msg}
	}

	parts := make([]string, 0)
	for start := 0; start <= len(msg); start += limit {
		parts = append(parts, msg[start:min(len(msg), start+limit)])
	}

	return parts
}

func (c *TClient) send(msg message) error {
	body, _ := json.Marshal(map[string]string{
		"chat_id": msg.chatID,
		"text":    msg.txt,
	})

	url := fmt.Sprintf("%s/bot%s/sendMessage", telegramApiURL, c.Token)
	resp, err := c.HttpClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		slog.Error("could not set message to Telegram", "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("message not sent", "resp.StatusCode", resp.StatusCode)
		return fmt.Errorf("message not sent")
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("could not read response from Telegram", "error", err)
		return err
	}

	slog.Debug("message sent to Telegram", "text", msg, "to", msg.chatID, "response", string(b))

	return nil
}

func DefaultChatIdFn() []string {
	if chatID, ok := os.LookupEnv("TELEGRAM_CHAT_ID"); ok {
		return []string{chatID}
	}

	slog.Warn("Using the default ChatIdFn and no TELEGRAM_CHAT_ID provided. No Telegram messages will be sent")

	return []string{}
}

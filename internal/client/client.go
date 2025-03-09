package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

const maxMessageSize = 2048

type TClient struct {
	Token  string
	ChatID string

	messages chan<- string
	errors   <-chan error
}

func (c *TClient) Run() {

	msgs := make(chan string, 16)
	c.messages = msgs

	errs := make(chan error, 1)
	c.errors = errs

	go func() {
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

func (c *TClient) Stop() {
	close(c.messages)
}

// SendMessage sends a message to the configured chat.
// Long messages are split in chunks that are then sent individually.
func (c *TClient) SendMessage(msg string) {
	parts := splitLongMessage(msg, maxMessageSize)
	for i := range parts {
		c.messages <- parts[i]
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

func (c *TClient) send(m string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.Token)

	body, _ := json.Marshal(map[string]string{
		"chat_id": c.ChatID,
		"text":    m,
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		slog.Error("could not set message to Telegram", "error", err)
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("could not read response from Telegram", "error", err)
		return err
	}

	slog.Debug("message sent to Telegram", "text", m, "from", c.ChatID, "response", string(b))

	return nil
}

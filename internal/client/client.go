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
	Token      string
	ChatID     string
	HttpClient *http.Client

	messages chan<- string
	errors   <-chan error

	done chan struct{}
}

func (c *TClient) Run() {

	msgs := make(chan string, 16)
	c.messages = msgs

	errs := make(chan error, 1)
	c.errors = errs

	c.done = make(chan struct{})

	go func() {
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
func (c *TClient) Stop() {
	close(c.messages)
	<-c.done
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
	return c.sendTo(m, url)
}

func (c *TClient) sendTo(msg, url string) error {
	body, _ := json.Marshal(map[string]string{
		"chat_id": c.ChatID,
		"text":    msg,
	})

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

	slog.Debug("message sent to Telegram", "text", msg, "from", c.ChatID, "response", string(b))

	return nil
}

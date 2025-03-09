package telegram

import "github.com/mehix/telegram/internal/client"

type Client interface {
	SendMessage(msg string)
	Stop()
}

// NewClient returns an initialized Telegram client.
// Call Stop() before exiting to avoid memory leaks.
func NewClient(token string, chatID string) *client.TClient {
	c := &client.TClient{
		Token:  token,
		ChatID: chatID,
	}
	c.Run()

	return c
}

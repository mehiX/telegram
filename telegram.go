package telegram

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/mehix/telegram/internal/client"
)

type Client interface {
	SendMessage(msg string)
	Stop()
}

// NewClient returns an initialized Telegram client.
// Call Stop() before exiting to avoid memory leaks.
func NewClient(token string, chatID string, opts ...Option) *client.TClient {
	c := &client.TClient{
		Token:      token,
		ChatID:     chatID,
		HttpClient: http.DefaultClient,
	}

	for _, o := range opts {
		o(c)
	}

	c.Run()

	return c
}

type Option func(*client.TClient)

func WithHttpClient(hc *http.Client) Option {
	return func(t *client.TClient) {
		if hc == nil {
			slog.Error("cannot work with a nil Http Client")
			os.Exit(1)
		}
		t.HttpClient = hc
	}
}

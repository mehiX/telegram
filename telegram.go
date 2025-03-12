package telegram

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/mehix/telegram/internal/client"
)

type Client interface {
	// SendTo sends the message to chatID. if the message is too long, it will be split and sent as multiple messages
	SendTo(msg string, chatID string)
	// Send sends the message to the chatID's returned by the ChatIdFn()
	Send(msg string)
	Errors() <-chan error
	Stop(context.Context) error
}

var _ Client = &client.TClient{}

// NewClient returns an initialized Telegram client.
// Call Stop() before exiting to avoid memory leaks.
func NewClient(token string, opts ...Option) *client.TClient {
	c := &client.TClient{
		Token:      token,
		HttpClient: http.DefaultClient,
		ChatIdFn:   client.DefaultChatIdFn,
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

func WithChatIDs(chatIDs ...string) Option {
	return func(t *client.TClient) {
		t.ChatIdFn = func() []string { return chatIDs }
	}
}

func WithChatIDsFn(f func() []string) Option {
	return func(t *client.TClient) {
		t.ChatIdFn = f
	}
}

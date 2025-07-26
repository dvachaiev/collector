package publisher

import (
	"log/slog"

	"collector/message"
)

type Publisher struct{}

func (p *Publisher) Publish(msg message.Message) {
	slog.Info("New message generated", "message", msg)
}

func (p *Publisher) Close() {
}

package processor

import (
	"fmt"
	"time"

	"github.com/dvachaiev/collector/message"
)

type Sensor interface {
	Value() int
}

type Publisher interface {
	Publish(msg message.Message) error
}

type Processor chan struct{}

func New(name string, sensor Sensor, rate int, publisher Publisher) Processor {
	p := make(Processor)

	go p.periodicPoll(name, sensor, rate, publisher)

	return p
}

func (p Processor) Close() {
	close(p)
}

func (p Processor) periodicPoll(name string, sensor Sensor, rate int, publisher Publisher) {
	ticker := time.NewTicker(time.Second / time.Duration(rate))
	defer ticker.Stop()

	for {
		msg := message.New(name, sensor.Value())
		if err := publisher.Publish(msg); err != nil {
			panic(fmt.Errorf("publish message: %w", err))
		}

		select {
		case <-p:
			return
		case <-ticker.C:
		}
	}
}

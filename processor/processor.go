package processor

import (
	"time"

	"collector/message"
)

type Sensor interface {
	Value() int
}

type Publisher interface {
	Publish(msg message.Message)
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
		publisher.Publish(msg)

		select {
		case <-p:
			return
		case <-ticker.C:
		}
	}
}

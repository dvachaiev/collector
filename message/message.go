package message

import "time"

type Message struct {
	Name  string
	Value int
	Time  time.Time
}

func New(name string, value int) Message {
	return Message{
		Name:  name,
		Value: value,
		Time:  time.Now(),
	}
}

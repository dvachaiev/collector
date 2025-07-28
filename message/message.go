package message

import (
	"time"
)

type Message struct {
	Time  int64  `json:"time"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func New(name string, value int) Message {
	return Message{
		Time:  time.Now().Unix(),
		Name:  name,
		Value: value,
	}
}

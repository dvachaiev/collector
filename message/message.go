package message

import (
	"encoding/json"
	"time"
)

type Message struct {
	Name  string    `json:"name"`
	Value int       `json:"value"`
	Time  time.Time `json:"time"`
}

func New(name string, value int) Message {
	return Message{
		Name:  name,
		Value: value,
		Time:  time.Now(),
	}
}

func (m Message) MarshalJSON() ([]byte, error) {
	type alias Message

	as := struct {
		Time int64 `json:"time"`
		alias
	}{
		Time:  m.Time.Unix(),
		alias: alias(m),
	}

	return json.Marshal(as)
}

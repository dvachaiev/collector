package converter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"collector/message"
)

func Convert(data []byte) ([]byte, error) {
	var (
		out []byte
		msg message.Message
	)

	now := time.Now().Unix()

	for line := range bytes.Lines(data) {
		if err := json.Unmarshal(line, &msg); err != nil {
			return nil, fmt.Errorf("parsing json: %w", err)
		}

		if err := Validate(msg, now); err != nil {
			return nil, fmt.Errorf("validating message: %w", err)
		}

		out = fmt.Appendf(out, "%v\t%v\t%v\n", msg.Time, msg.Name, msg.Value)
	}

	return out, nil
}

func Validate(msg message.Message, now int64) error {
	if msg.Time < 0 || msg.Time > now {
		return fmt.Errorf("invalid timestamp: %v", msg.Time)
	}

	if msg.Value < 0 {
		return fmt.Errorf("invalid value: %v", msg.Value)
	}

	if msg.Name == "" {
		return fmt.Errorf("empty name")
	}

	return nil
}

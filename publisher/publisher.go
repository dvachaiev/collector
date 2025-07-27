package publisher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"collector/message"
)

type Publisher struct {
	url       string
	client    http.Client
	hasData   chan struct{}
	closing   chan struct{}
	closed    chan struct{}
	mu        sync.Mutex
	buf       *buffer
	discarded int
	reqBuf    []byte
}

func New(url string, bufferSize int) *Publisher {
	p := Publisher{
		url: url,
		client: http.Client{
			Timeout: 10 * time.Second,
		},
		buf:     newBuffer(bufferSize),
		hasData: make(chan struct{}, 1),
		closing: make(chan struct{}),
		closed:  make(chan struct{}),
	}

	go p.runPublishing()

	return &p
}

func (p *Publisher) Publish(msg message.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("encoding message to json: %w", err)
	}

	data = append(data, '\n')

	p.mu.Lock()
	p.discarded += p.buf.Append(data)
	p.mu.Unlock()

	select {
	case p.hasData <- struct{}{}:
	default:
	}

	return nil
}

func (p *Publisher) Close() {
	close(p.closing)
	<-p.closed
}

func (p *Publisher) runPublishing() {
	for {
		select {
		case <-p.closing:
			close(p.closed)
			return
		case <-p.hasData:
			if ok := p.publish(); ok {
				continue
			}

			// Request wasn't successful, retrying
			select {
			case p.hasData <- struct{}{}:
			default:
			}
		}
	}
}

func (p *Publisher) publish() (ok bool) {
	published, data := p.getData(0)

	code, err := p.sendData(data)
	if err != nil {
		slog.Warn("Failed to send data", "error", err)

		return false
	}

	if code != http.StatusOK && code != http.StatusRequestEntityTooLarge {
		slog.Warn("Not expected response code", "code", err)

		return false
	}

	for code == http.StatusRequestEntityTooLarge {
		if published == 1 {
			break
		}

		published /= 2

		slog.Info("Hit rate limits, trying to reduce number of messages in request", "messages", published)

		_, data = p.getData(published)

		code, err = p.sendData(data)
		if err != nil {
			slog.Warn("Failed to send data", "error", err)

			return false
		}
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.buf.Remove(published - p.discarded)

	return true
}

func (p *Publisher) getData(num int) (int, []byte) {
	var data []byte

	p.mu.Lock()
	defer p.mu.Unlock()

	p.discarded = 0

	if num > 0 {
		data = p.buf.Get(num)
	} else {
		num, data = p.buf.GetAll()
	}

	if cap(p.reqBuf) < len(data) {
		p.reqBuf = make([]byte, len(data))
	}

	p.reqBuf = p.reqBuf[:len(data)]
	copy(p.reqBuf, data)

	return num, p.reqBuf
}

func (p *Publisher) sendData(data []byte) (int, error) {
	resp, err := p.client.Post(p.url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}

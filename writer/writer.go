package writer

import (
	"io"
	"log/slog"
	"time"
)

type Buffered struct {
	w       io.Writer
	writes  chan writeOp
	closing chan struct{}
	closed  chan error
	buf     []byte
	used    int
}

func New(w io.Writer, size int, interval time.Duration) *Buffered {
	b := Buffered{
		w:       w,
		writes:  make(chan writeOp),
		closing: make(chan struct{}),
		closed:  make(chan error),
		buf:     make([]byte, size),
		used:    0,
	}

	go b.processWrites(interval)

	return &b
}

func (b *Buffered) Write(p []byte) (n int, err error) {
	resCh := make(chan writeRes, 1)
	b.writes <- writeOp{data: p, resCh: resCh}
	res := <-resCh

	return res.n, res.err
}

func (b *Buffered) Close() error {
	close(b.closing)
	return <-b.closed
}

func (b *Buffered) processWrites(flushInterval time.Duration) {
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-b.closing:
			b.closed <- b.flush()
			close(b.writes)
			close(b.closed)

			return
		case <-ticker.C:
			if err := b.flush(); err != nil {
				slog.Error("Error on periodic flush", "error", err)
			}
		case w := <-b.writes:
			flushed, n, err := b.write(w.data)

			if flushed {
				ticker.Reset(flushInterval)
			}

			w.resCh <- writeRes{n: n, err: err}
		}
	}
}

func (b *Buffered) write(p []byte) (flushed bool, n int, err error) {
	if b.used > 0 && len(p) > len(b.buf)-b.used {
		if err := b.flush(); err != nil {
			return false, 0, err
		}

		flushed = true
	}

	if len(p) > len(b.buf) {
		n, err = b.w.Write(p)
		return true, n, err
	}

	n = copy(b.buf[b.used:], p)
	b.used += n

	return flushed, n, nil
}

func (b *Buffered) flush() error {
	if b.used == 0 {
		return nil
	}

	n, err := b.w.Write(b.buf[:b.used])
	if n > 0 {
		n = copy(b.buf, b.buf[n:b.used])
		b.used = n
	}

	return err
}

type writeOp struct {
	data  []byte
	resCh chan writeRes
}

type writeRes struct {
	n   int
	err error
}

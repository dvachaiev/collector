package writer

import (
	"io"
	"time"
)

type Buffered struct {
	w    io.Writer
	buf  []byte
	used int
}

func New(w io.Writer, size int, interval time.Duration) *Buffered {
	return &Buffered{
		w:    w,
		buf:  make([]byte, size),
		used: 0,
	}
}

func (b *Buffered) Write(p []byte) (n int, err error) {
	if len(p) > len(b.buf)-b.used {
		if err := b.flush(); err != nil {
			return 0, err
		}
	}

	if len(p) > len(b.buf) {
		return b.w.Write(p)
	}

	n = copy(b.buf[b.used:], p)
	b.used += n

	return n, nil
}

func (b *Buffered) Close() error {
	return b.flush()
}

func (b *Buffered) flush() error {
	n, err := b.w.Write(b.buf[:b.used])
	if n > 0 {
		n = copy(b.buf, b.buf[n:b.used])
		b.used = n
	}

	return err
}

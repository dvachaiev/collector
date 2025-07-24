package handler

import (
	"net/http"
	"sync"
	"time"
)

type Limiter struct {
	next      http.Handler
	byteDur   time.Duration
	maxCount  int64
	mu        sync.Mutex
	available int64
	last      time.Time
}

func NewLimiter(next http.Handler, count int64, period time.Duration) *Limiter {
	return &Limiter{
		next:      next,
		maxCount:  count,
		byteDur:   period / time.Duration(count),
		available: count,
	}
}

func (l *Limiter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cl := r.ContentLength

	if cl < 0 {
		w.WriteHeader(http.StatusLengthRequired)

		return
	}

	if !l.allow(time.Now(), cl) {
		w.WriteHeader(http.StatusRequestEntityTooLarge)

		return
	}

	l.next.ServeHTTP(w, r)
}

func (l *Limiter) allow(now time.Time, size int64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.available = min(l.available+int64(now.Sub(l.last)/l.byteDur), l.maxCount)
	l.last = now

	if next := l.available - size; next >= 0 {
		l.available = next

		return true
	}

	return false
}

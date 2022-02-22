package tokenbucket

import (
	"sync"
	"time"
)

type TokenBucket struct {
	mu       sync.Mutex
	size     int64
	perToken time.Duration
	token    int64
	last     time.Time // last put token time
}

func New(rate int64, size int64) *TokenBucket {
	return &TokenBucket{
		size:     size,
		perToken: time.Second / time.Duration(rate),
		last:     time.Now(),
	}
}

func (t *TokenBucket) Take() bool {
	t.mu.Lock()

	now := time.Now()
	t.token += int64(now.Sub(t.last)) / int64(t.perToken)
	if t.token > t.size {
		t.token = t.size
	}
	t.last = now
	t.token -= 1
	if t.token < 0 {
		t.mu.Unlock()
		time.Sleep(time.Duration(-t.token) * t.perToken)
		return true
	}
	t.mu.Unlock()
	return true
}

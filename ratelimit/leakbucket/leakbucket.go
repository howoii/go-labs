package leakbucket

import (
	"sync"
	"time"
)

type LeakBucket struct {
	mu         sync.Mutex
	last       time.Time     // last take time
	perRequest time.Duration // ns per request
	cap        int           // max num of waiting request
	sleep      time.Duration // sleep time of request
}

func New(rate int, cap int) *LeakBucket {
	l := &LeakBucket{
		last:       time.Now(),
		perRequest: time.Second / time.Duration(rate),
		cap:        cap,
	}

	return l
}

func (l *LeakBucket) Take() bool {
	l.mu.Lock()

	sleep := l.sleep + l.perRequest - time.Now().Sub(l.last)
	if sleep > 0 {
		if sleep > l.perRequest*time.Duration(l.cap) {
			l.mu.Unlock()
			return false
		}
		l.sleep = 0
		l.last = time.Now().Add(sleep)
		l.mu.Unlock()

		time.Sleep(sleep)
		return true
	}

	l.sleep = sleep
	l.last = time.Now()
	l.mu.Unlock()
	return true
}

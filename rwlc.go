package rwlc

import (
	"errors"
	"sync"
)

var ErrClosed = errors.New("closed")

type ReadWriteLineCloser interface {
	ReadLine() (s string, err error)
	WriteLine(s string) error
	Close()
}

type rwlc struct {
	mu     sync.Mutex
	head   chan string
	tail   []string
	done   chan struct{}
	closed bool
}

func New() ReadWriteLineCloser {
	return &rwlc{
		head: make(chan string, 1),
		done: make(chan struct{}),
	}
}

func (rw *rwlc) ReadLine() (s string, err error) {
	rw.mu.Lock()
	closed := rw.closed
	rw.mu.Unlock()
	if closed {
		return s, ErrClosed
	}

	select {
	case s = <-rw.head:
	case <-rw.done:
		return s, ErrClosed
	}

	rw.mu.Lock()
	if len(rw.tail) > 0 {
		rw.tryPushLocked()
	}
	rw.mu.Unlock()

	return s, nil
}

func (rw *rwlc) WriteLine(s string) error {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	if rw.closed {
		return ErrClosed
	}

	rw.tail = append(rw.tail, s)

	rw.tryPushLocked()

	return nil
}

// tries to push the head-of-tail to the head
func (rw *rwlc) tryPushLocked() {
	select {
	case rw.head <- rw.tail[0]:
		rw.tail = rw.tail[1:]
	default:
	}
}

func (rw *rwlc) Close() {
	rw.mu.Lock()
	if rw.closed {
		rw.mu.Unlock()
		return
	}
	rw.closed = true
	rw.mu.Unlock()

	close(rw.done)
}

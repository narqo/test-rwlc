package rwlc

import (
	"errors"
)

var ErrClosed = errors.New("closed")

type ReadWriteLineCloser interface {
	ReadLine() (s string, err error)
	WriteLine(s string) error
	Close()
}

type rwlc struct {
	lines chan []string
	empty chan struct{}

	closed chan struct{}
}

func New() ReadWriteLineCloser {
	lines := make(chan []string, 1)

	empty := make(chan struct{}, 1)
	empty <- struct{}{}

	return &rwlc{
		lines: lines,
		empty: empty,

		closed: make(chan struct{}),
	}
}

func (rw *rwlc) ReadLine() (s string, err error) {
	var lines []string
	select {
	case <-rw.closed:
		return "", ErrClosed
	case lines = <-rw.lines:
	}

	s = lines[0]
	if len(lines) > 1 {
		rw.lines <- lines[1:]
	} else {
		rw.empty <- struct{}{}
	}

	return s, nil
}

func (rw *rwlc) WriteLine(s string) error {
	var lines []string

	select {
	case <-rw.closed:
		return ErrClosed
	case lines = <-rw.lines:
	case <-rw.empty:
	}

	lines = append(lines, s)
	rw.lines <- lines

	return nil
}

func (rw *rwlc) Close() {
	close(rw.closed)
}

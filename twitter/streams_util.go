package twitter

import (
	"strings"
	"time"
)

func stopped(done <-chan struct{}) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func sleepOrDone(d time.Duration, done <-chan struct{}) {
	select {
	case <-time.After(d):
		return
	case <-done:
		return
	}
}

func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := strings.Index(string(data), "\r\n"); i >= 0 {
		return i + 2, data[0:i], nil
	}
	if atEOF {
		return len(data), dropCR(data), nil
	}
	return 0, nil, nil
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\n' {
		return data[0 : len(data)-1]
	}
	return data
}

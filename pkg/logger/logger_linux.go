//go:build linux
// +build linux

package logger

import (
	"strconv"

	"golang.org/x/sys/unix"
)

func (l *Logger) getThreadId() (threadId string) {
	tid := unix.Gettid()
	threadId = strconv.Itoa(tid)
	return threadId
}

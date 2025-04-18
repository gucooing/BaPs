//go:build windows
// +build windows

package logger

import (
	"strconv"

	"golang.org/x/sys/windows"
)

func (l *Logger) getThreadId() (threadId string) {
	tid := windows.GetCurrentThreadId()
	threadId = strconv.Itoa(int(tid))
	return threadId
}

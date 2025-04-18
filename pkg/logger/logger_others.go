//go:build !windows && !linux
// +build !windows,!linux

package logger

func (l *Logger) getThreadId() string {
	return ""
}

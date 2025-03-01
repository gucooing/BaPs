//go:build !dev
// +build !dev

package logger

func (l *Logger) getLineFunc() (fileName string, line int, funcName string) {
	return "", 0, ""
}

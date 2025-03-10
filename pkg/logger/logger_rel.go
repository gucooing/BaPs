//go:build !dev
// +build !dev

package logger

func (l *Logger) getLineFunc(logInfo *LogInfo) {
	logInfo.TrackLine = false
}

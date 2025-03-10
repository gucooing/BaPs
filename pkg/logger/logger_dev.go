//go:build dev
// +build dev

package logger

import (
	"path"
	"runtime"
	"strings"
)

func (l *Logger) getLineFunc(logInfo *LogInfo) {
	var pc uintptr
	var file string
	var line int
	var ok bool
	pc, file, line, ok = runtime.Caller(3)
	if !ok {
		return
	}
	fileName := path.Base(file)
	funcName := runtime.FuncForPC(pc).Name()
	split := strings.Split(funcName, ".")
	if len(split) != 0 {
		funcName = split[len(split)-1]
	}
	logInfo.TrackLine = true
	logInfo.FileName, logInfo.Line, logInfo.FuncName = fileName, line, funcName
}

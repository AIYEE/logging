package logging

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type CallerHook struct{}

func (h *CallerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *CallerHook) Fire(entry *logrus.Entry) error {
	entry.Data["file"] = getCallerInfo()
	return nil
}

func getCallerInfo() string {
	var pc [16]uintptr
	n := runtime.Callers(2, pc[:])
	frames := runtime.CallersFrames(pc[:n])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if !strings.Contains(frame.Function, "logrus") && !strings.Contains(frame.Function, "runtime") {
			// return filepath.Base(frame.File) + ":" + string(frame.Line)
			return fmt.Sprintf("%s:%d", filepath.Base(frame.File), frame.Line)
		}
	}
	return "unknown"
}

// func getCallerInfo() string {
// 	_, file, line, ok := runtime.Caller(9)
// 	if ok {
// 		return fmt.Sprintf("%s:%d", filepath.Base(file), line)
// 	}
// 	return "unknown"
// }

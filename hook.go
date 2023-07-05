package logging

import (
	"fmt"
	"path/filepath"
	"runtime"

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
	_, file, line, ok := runtime.Caller(9)
	if ok {
		return fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
	return "unknown"
}

package logging

import (
	"os"
	"testing"
)

func Test_NewLoggerWithFile(t *testing.T) {
	logger, err := NewLogger(WithLogFile("./test.log"))
	if err != nil {
		t.Error(err)
	}

	logger.Info("test ok")
	t.Log("ok")
}

func Test_NewLoggerWithDefault(t *testing.T) {
	logger, err := NewLogger(WithWriter(os.Stdout))
	if err != nil {
		t.Error(err)
	}

	logger.Info("test ok")
	t.Log("ok")
}

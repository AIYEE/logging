package logging

import (
	"io"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

const (
	maxFiles  uint  = 10
	rotatSize int64 = 5 * 1024 * 1024
)

type Logger interface {
	Tracef(format string, args ...interface{})
	Trace(args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warningf(format string, args ...interface{})
	Warning(args ...interface{})
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	WithField(key string, value interface{}) *logrus.Entry
	WithFields(fields logrus.Fields) *logrus.Entry
	WriterLevel(logrus.Level) *io.PipeWriter
	NewEntry() *logrus.Entry
}

type logger struct {
	*logrus.Logger
}

func CreateFileWriter(fileName string) io.Writer {
	writer, _ := rotatelogs.New(
		fileName+".%Y-%m-%d-%H-%M",
		rotatelogs.WithLinkName(fileName),
		rotatelogs.WithRotationCount(maxFiles),
		rotatelogs.WithRotationSize(rotatSize),
	)
	return writer
}

func New(w io.Writer, level logrus.Level) Logger {
	l := logrus.New()
	l.SetOutput(w)
	l.SetLevel(level)
	l.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	return &logger{
		Logger: l,
	}
}

func New1(fileName string, level logrus.Level, maxFiles uint, rotatSize int64) Logger {
	writer, _ := rotatelogs.New(
		fileName+".%Y-%m-%d-%H-%M",
		rotatelogs.WithLinkName(fileName),
		rotatelogs.WithRotationCount(maxFiles),
		rotatelogs.WithRotationSize(rotatSize),
	)

	l := logrus.New()
	l.SetOutput(writer)
	l.SetLevel(level)
	l.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	return &logger{
		Logger: l,
	}
}

func (l *logger) NewEntry() *logrus.Entry {
	return logrus.NewEntry(l.Logger)
}


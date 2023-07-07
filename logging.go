package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

const (
	defaultMaxFiles  uint          = 10
	defaultRotatSize int64         = 5 * 1024 * 1024
	defaultMaxAge    time.Duration = 0
	defaultVerbosity string        = "info"
)

type LoggerParam struct {
	LogFile   string
	Writer    io.Writer
	MaxFiles  uint
	RotatSize int64
	MaxAge    time.Duration
	Verbosity string
}

type Option func(opt *LoggerParam)

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
	Close()
}

type logger struct {
	*logrus.Logger
}

func NewLogger(opts ...Option) (Logger, error) {
	params := defaultLoggerParam()
	for _, o := range opts {
		o(params)
	}
	if params.LogFile != "" {
		if err := params.createFileWriter(); err != nil {
			return nil, err
		}
	}

	return params.newLogger()
}

func New(w io.Writer, level logrus.Level) Logger {
	l := logrus.New()
	l.SetOutput(w)
	l.SetLevel(level)
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			function = fmt.Sprintf("%s()", filepath.Base(f.Function))
			file = fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line)
			return
		},
	}
	// l.AddHook(&CallerHook{})

	return &logger{
		Logger: l,
	}
}

func (l *logger) NewEntry() *logrus.Entry {
	return logrus.NewEntry(l.Logger)
}

func (l *logger) Close() {
	if l.Out != nil {
		l.Out.(*rotatelogs.RotateLogs).Close()
	}
}

func WithLogFile(file string) Option {
	return func(l *LoggerParam) {
		if file != "" {
			l.LogFile = file
		}
	}
}

func WithWriter(writer io.Writer) Option {
	return func(l *LoggerParam) {
		if writer != nil {
			l.Writer = writer
		}
	}
}

func WithMaxFiles(max uint) Option {
	return func(l *LoggerParam) {
		if max != 0 {
			l.MaxFiles = max
		} else {
			l.MaxFiles = defaultMaxFiles
		}

	}
}

func WithRotatSize(size int64) Option {
	return func(l *LoggerParam) {
		if size != 0 {
			l.RotatSize = size
		} else {
			l.RotatSize = defaultRotatSize
		}

	}
}

// in days measured
func WithMaxAge(maxAge int64) Option {
	return func(l *LoggerParam) {
		if maxAge != 0 {
			l.MaxAge = time.Duration(maxAge) * 24 * time.Hour
		} else {
			l.MaxAge = defaultMaxAge
		}

	}
}

func WithVerbosity(verbosity string) Option {
	return func(l *LoggerParam) {
		if verbosity != "" {
			l.Verbosity = strings.ToLower(verbosity)
		} else {
			l.Verbosity = defaultVerbosity
		}

	}
}

func (l *LoggerParam) newLogger() (Logger, error) {
	var logger Logger
	switch l.Verbosity {
	case "0", "silent":
		logger = New(io.Discard, 0)
	case "1", "error":
		logger = New(l.Writer, logrus.ErrorLevel)
	case "2", "warn":
		logger = New(l.Writer, logrus.WarnLevel)
	case "3", "info":
		logger = New(l.Writer, logrus.InfoLevel)
	case "4", "debug":
		logger = New(l.Writer, logrus.DebugLevel)
	case "5", "trace":
		logger = New(l.Writer, logrus.TraceLevel)
	default:
		return nil, fmt.Errorf("unknown verbosity level %q", l.Verbosity)
	}
	return logger, nil
}

func defaultLoggerParam() *LoggerParam {
	return &LoggerParam{
		Writer:    os.Stdout,
		MaxFiles:  defaultMaxFiles,
		RotatSize: defaultRotatSize,
		MaxAge:    defaultMaxAge,
		Verbosity: defaultVerbosity,
	}
}

// options MaxAge and RotationCount cannot be both set
func (l *LoggerParam) createFileWriter() error {
	var err error
	l.Writer, err = rotatelogs.New(
		l.LogFile+".%Y-%m-%d-%H-%M",
		rotatelogs.WithLinkName(l.LogFile),
		rotatelogs.WithRotationCount(l.MaxFiles),
		rotatelogs.WithRotationSize(l.RotatSize),
		rotatelogs.WithMaxAge(l.MaxAge),
	)
	return err
}

package log

import (
	"fmt"
	"io"
	"os"

	"github.com/nartvt/go-core/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/sirupsen/logrus"
)

var _ log.Logger = (*LogrusLogger)(nil)

type LogrusLogger struct {
	log *logrus.Logger
}

func NewLogrusLogger(options ...Option) log.Logger {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	logger.Out = os.Stdout
	logger.Formatter = &logrus.TextFormatter{}
	for _, option := range options {
		option(logger)
	}

	return &LogrusLogger{
		log: logger,
	}
}

func LogrusConfig(logConf *conf.Server_Log) log.Logger {
	logLevel := logrus.DebugLevel
	out := os.Stdout
	var format logrus.Formatter

	format = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "@timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	}
	if logConf != nil {
		lv, err := logrus.ParseLevel(logConf.Level)
		if err != nil {
			fmt.Println("Logger level error, fallback to default id debug ", err)
		} else {
			logLevel = lv
		}

		if logConf.Format == "text" {
			format = &logrus.TextFormatter{}
		}

		if len(logConf.File) > 0 {
			// You could set this to any `io.Writer` such as a file
			file, err := os.OpenFile(logConf.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err == nil {
				out = file
			} else {
				fmt.Println("Failed to log to file, using default stdout")
			}
		}
	}

	return NewLogrusLogger(Level(logLevel), Formatter(format), Output(out))
}

func (l *LogrusLogger) Log(level log.Level, keyvals ...interface{}) (err error) {
	var (
		logrusLevel logrus.Level
		fields      logrus.Fields = make(map[string]interface{})
		msg         string
	)

	switch level {
	case log.LevelDebug:
		logrusLevel = logrus.DebugLevel
	case log.LevelInfo:
		logrusLevel = logrus.InfoLevel
	case log.LevelWarn:
		logrusLevel = logrus.WarnLevel
	case log.LevelError:
		logrusLevel = logrus.ErrorLevel
	case log.LevelFatal:
		logrusLevel = logrus.FatalLevel
	default:
		logrusLevel = logrus.DebugLevel
	}

	if logrusLevel > l.log.Level {
		return
	}

	if len(keyvals) == 0 {
		return nil
	}
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "")
	}
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			continue
		}
		if key == logrus.FieldKeyMsg {
			msg, _ = keyvals[i+1].(string)
			continue
		}
		fields[key] = keyvals[i+1]
	}

	if len(fields) > 0 {
		l.log.WithFields(fields).Log(logrusLevel, msg)
	} else {
		l.log.Log(logrusLevel, msg)
	}

	return
}

type Option func(log *logrus.Logger)

func Level(level logrus.Level) Option {
	return func(log *logrus.Logger) {
		log.Level = level
	}
}

func Output(w io.Writer) Option {
	return func(log *logrus.Logger) {
		log.Out = w
	}
}

func Formatter(formatter logrus.Formatter) Option {
	return func(log *logrus.Logger) {
		log.Formatter = formatter
	}
}

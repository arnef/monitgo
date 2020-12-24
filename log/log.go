package log

import (
	"fmt"
	"io"
	"runtime"

	log "github.com/sirupsen/logrus"
)

type Level uint32

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

func SetOutput(out io.Writer) {
	log.SetOutput(out)
	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true})
}

func SetLevel(level Level) {
	switch level {
	case PanicLevel:
		log.SetLevel(log.PanicLevel)
	case FatalLevel:
		log.SetLevel(log.FatalLevel)
	case ErrorLevel:
		log.SetLevel(log.ErrorLevel)
	case WarnLevel:
		log.SetLevel(log.WarnLevel)
	case InfoLevel:
		log.SetLevel(log.InfoLevel)
	case DebugLevel:
		log.SetLevel(log.DebugLevel)
	case TraceLevel:
		log.SetLevel(log.TraceLevel)
	}
}

func Debug(any interface{}, args ...interface{}) {
	if any != nil {
		err := (error)(nil)
		switch any.(type) {
		case string:
			err = fmt.Errorf(any.(string), args...)
		case error:
			err = fmt.Errorf(any.(error).Error(), args...)
		default:
			err = fmt.Errorf("%v", err)
		}

		_, fn, line, _ := runtime.Caller(1)
		log.Debugf("%s:%d %v\n", fn, line, err)
	}
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}
